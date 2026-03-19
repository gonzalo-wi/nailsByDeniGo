package repositories

import (
	"errors"
	"time"

	"apiGoShei/internal/domain/schedule"
	"apiGoShei/internal/infrastructure/persistence/postgres/models"

	"gorm.io/gorm"
)

type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) FindWeeklySchedule() ([]schedule.WeeklySchedule, error) {
	var rows []models.WeeklyScheduleModel
	if err := r.db.Order("day_of_week ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]schedule.WeeklySchedule, len(rows))
	for i, m := range rows {
		result[i] = toWeeklyScheduleDomain(&m)
	}
	return result, nil
}

func (r *ScheduleRepository) UpsertWeeklySchedule(schedules []schedule.WeeklySchedule) error {
	for _, s := range schedules {
		var m models.WeeklyScheduleModel
		result := r.db.Where("day_of_week = ?", s.DayOfWeek).First(&m)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		m.DayOfWeek = s.DayOfWeek
		m.Enabled = s.Enabled
		m.OpeningTime = s.OpeningTime
		m.ClosingTime = s.ClosingTime
		m.SlotDurationMin = s.SlotDurationMin
		if err := r.db.Save(&m).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *ScheduleRepository) CreateBlockedSlot(slot *schedule.BlockedSlot) error {
	m := &models.BlockedSlotModel{
		Date:      slot.Date,
		StartTime: slot.StartTime,
		EndTime:   slot.EndTime,
		Reason:    slot.Reason,
		Permanent: slot.Permanent,
	}
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	slot.ID = m.ID
	slot.CreatedAt = m.CreatedAt
	return nil
}

func (r *ScheduleRepository) FindBlockedSlots(date string) ([]schedule.BlockedSlot, error) {
	var rows []models.BlockedSlotModel
	if err := r.db.Where("date = ? OR permanent = ?", date, true).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]schedule.BlockedSlot, len(rows))
	for i, m := range rows {
		result[i] = toBlockedSlotDomain(&m)
	}
	return result, nil
}

// IsAvailable verifica si un slot [startTime, startTime+durationMinutes) está disponible
// para la fecha dada. Comprueba: horario habilitado + no bloqueado.
func (r *ScheduleRepository) IsAvailable(date, startTime string, durationMinutes int) (bool, error) {
	// Calcular día de semana a partir de la fecha
	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false, err
	}
	dayOfWeek := int(d.Weekday())

	// Obtener configuración semanal
	var ws models.WeeklyScheduleModel
	if err := r.db.Where("day_of_week = ?", dayOfWeek).First(&ws).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // día no configurado = no disponible
		}
		return false, err
	}
	if !ws.Enabled {
		return false, nil
	}

	// Calcular endTime como string "HH:MM"
	st, err := time.Parse("15:04", startTime)
	if err != nil {
		return false, err
	}
	endTime := st.Add(time.Duration(durationMinutes) * time.Minute).Format("15:04")

	// Verificar que el slot esté dentro del horario de apertura/cierre
	if startTime < ws.OpeningTime || endTime > ws.ClosingTime {
		return false, nil
	}

	// Verificar que no haya bloqueos solapados
	var count int64
	r.db.Model(&models.BlockedSlotModel{}).
		Where("(date = ? OR permanent = ?) AND start_time < ? AND end_time > ?",
			date, true, endTime, startTime).
		Count(&count)

	return count == 0, nil
}

func toWeeklyScheduleDomain(m *models.WeeklyScheduleModel) schedule.WeeklySchedule {
	return schedule.WeeklySchedule{
		ID:              m.ID,
		DayOfWeek:       m.DayOfWeek,
		Enabled:         m.Enabled,
		OpeningTime:     m.OpeningTime,
		ClosingTime:     m.ClosingTime,
		SlotDurationMin: m.SlotDurationMin,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func toBlockedSlotDomain(m *models.BlockedSlotModel) schedule.BlockedSlot {
	return schedule.BlockedSlot{
		ID:        m.ID,
		Date:      m.Date,
		StartTime: m.StartTime,
		EndTime:   m.EndTime,
		Reason:    m.Reason,
		Permanent: m.Permanent,
		CreatedAt: m.CreatedAt,
	}
}
