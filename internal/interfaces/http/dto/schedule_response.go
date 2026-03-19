package dto

import "apiGoShei/internal/domain/schedule"

type WeeklyScheduleResponse struct {
	ID              uint   `json:"id"`
	DayOfWeek       int    `json:"day_of_week"`
	DayName         string `json:"day_name"`
	Enabled         bool   `json:"enabled"`
	OpeningTime     string `json:"opening_time"`
	ClosingTime     string `json:"closing_time"`
	SlotDurationMin int    `json:"slot_duration_min"`
}

var dayNames = map[int]string{
	0: "Domingo", 1: "Lunes", 2: "Martes", 3: "Miércoles",
	4: "Jueves", 5: "Viernes", 6: "Sábado",
}

func WeeklyScheduleToResponse(s schedule.WeeklySchedule) WeeklyScheduleResponse {
	return WeeklyScheduleResponse{
		ID:              s.ID,
		DayOfWeek:       s.DayOfWeek,
		DayName:         dayNames[s.DayOfWeek],
		Enabled:         s.Enabled,
		OpeningTime:     s.OpeningTime,
		ClosingTime:     s.ClosingTime,
		SlotDurationMin: s.SlotDurationMin,
	}
}

type BlockedSlotResponse struct {
	ID        uint   `json:"id"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Reason    string `json:"reason,omitempty"`
	Permanent bool   `json:"permanent"`
}

func BlockedSlotToResponse(b schedule.BlockedSlot) BlockedSlotResponse {
	return BlockedSlotResponse{
		ID:        b.ID,
		Date:      b.Date,
		StartTime: b.StartTime,
		EndTime:   b.EndTime,
		Reason:    b.Reason,
		Permanent: b.Permanent,
	}
}
