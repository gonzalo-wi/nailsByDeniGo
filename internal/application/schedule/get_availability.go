package scheduleapp

import (
	"time"

	"apiGoShei/internal/domain/appointment"
	"apiGoShei/internal/domain/schedule"
)

type GetAvailabilityInput struct {
	Date            string
	DurationMinutes int // Si es 0, usa el slot_duration de la configuración
}

type AvailableSlot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type GetAvailabilityOutput struct {
	Date           string          `json:"date"`
	OpeningTime    string          `json:"opening_time"`
	ClosingTime    string          `json:"closing_time"`
	AvailableSlots []AvailableSlot `json:"available_slots"`
}

type GetAvailabilityUseCase struct {
	scheduleRepo    schedule.Repository
	appointmentRepo appointment.Repository
}

func NewGetAvailabilityUseCase(sr schedule.Repository, ar appointment.Repository) *GetAvailabilityUseCase {
	return &GetAvailabilityUseCase{scheduleRepo: sr, appointmentRepo: ar}
}

func (uc *GetAvailabilityUseCase) Execute(input GetAvailabilityInput) (*GetAvailabilityOutput, error) {
	// Parsear fecha y obtener día de semana
	d, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return nil, err
	}
	dayOfWeek := int(d.Weekday())

	// Obtener configuración semanal
	schedules, err := uc.scheduleRepo.FindWeeklySchedule()
	if err != nil {
		return nil, err
	}

	var daySchedule *schedule.WeeklySchedule
	for i := range schedules {
		if schedules[i].DayOfWeek == dayOfWeek {
			daySchedule = &schedules[i]
			break
		}
	}

	if daySchedule == nil || !daySchedule.Enabled {
		return &GetAvailabilityOutput{Date: input.Date, AvailableSlots: []AvailableSlot{}}, nil
	}

	duration := input.DurationMinutes
	if duration <= 0 {
		duration = daySchedule.SlotDurationMin
	}

	// Obtener turnos y bloqueos existentes para la fecha
	existingAppointments, err := uc.appointmentRepo.FindByDateRange(input.Date, input.Date)
	if err != nil {
		return nil, err
	}
	blockedSlots, err := uc.scheduleRepo.FindBlockedSlots(input.Date)
	if err != nil {
		return nil, err
	}

	// Generar todos los slots posibles y filtrar los ocupados
	slots := generateAvailableSlots(
		daySchedule.OpeningTime,
		daySchedule.ClosingTime,
		daySchedule.SlotDurationMin,
		duration,
		existingAppointments,
		blockedSlots,
	)

	return &GetAvailabilityOutput{
		Date:           input.Date,
		OpeningTime:    daySchedule.OpeningTime,
		ClosingTime:    daySchedule.ClosingTime,
		AvailableSlots: slots,
	}, nil
}

// generateAvailableSlots recorre los slots en pasos de stepMin y devuelve
// solo aquellos que no se solapan con turnos existentes ni con bloqueos.
func generateAvailableSlots(
	opening, closing string,
	stepMin, durationMin int,
	appointments []appointment.Appointment,
	blocked []schedule.BlockedSlot,
) []AvailableSlot {
	const layout = "15:04"
	start, err := time.Parse(layout, opening)
	if err != nil {
		return nil
	}
	close, err := time.Parse(layout, closing)
	if err != nil {
		return nil
	}

	var result []AvailableSlot
	for t := start; ; t = t.Add(time.Duration(stepMin) * time.Minute) {
		slotEnd := t.Add(time.Duration(durationMin) * time.Minute)
		if slotEnd.After(close) {
			break
		}
		stStr := t.Format(layout)
		etStr := slotEnd.Format(layout)

		if !isOccupied(stStr, etStr, appointments, blocked) {
			result = append(result, AvailableSlot{StartTime: stStr, EndTime: etStr})
		}
	}
	return result
}

func isOccupied(startTime, endTime string, appts []appointment.Appointment, blocked []schedule.BlockedSlot) bool {
	for _, a := range appts {
		if a.Status == appointment.StatusCancelled || a.Status == appointment.StatusAbsent {
			continue
		}
		// Overlap: a.StartTime < endTime AND a.EndTime > startTime
		if a.StartTime < endTime && a.EndTime > startTime {
			return true
		}
	}
	for _, b := range blocked {
		if b.StartTime < endTime && b.EndTime > startTime {
			return true
		}
	}
	return false
}
