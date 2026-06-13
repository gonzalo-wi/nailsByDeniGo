package repositories

import (
	"errors"

	"apiGoShei/internal/domain/appointment"
	"apiGoShei/internal/infrastructure/persistence/postgres/models"

	"gorm.io/gorm"
)

type AppointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

func (r *AppointmentRepository) Create(entity *appointment.Appointment) error {
	m := models.AppointmentModel{
		ClientID:       entity.ClientID,
		ServiceID:      entity.ServiceID,
		ProfessionalID: entity.ProfessionalID,
		Date:           entity.Date,
		StartTime:      entity.StartTime,
		EndTime:        entity.EndTime,
		BasePrice:      entity.BasePrice,
		ExtrasAmount:   entity.ExtrasAmount,
		DepositAmount:  entity.DepositAmount,
		FinalPrice:     entity.FinalPrice,
		Status:         string(entity.Status),
		Notes:          entity.Notes,
	}
	if err := r.db.Create(&m).Error; err != nil {
		return err
	}
	entity.ID = m.ID
	entity.CreatedAt = m.CreatedAt
	entity.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *AppointmentRepository) Update(entity *appointment.Appointment) error {
	return r.db.Model(&models.AppointmentModel{}).Where("id = ?", entity.ID).Updates(map[string]interface{}{
		"status":         string(entity.Status),
		"service_id":     entity.ServiceID,
		"base_price":     entity.BasePrice,
		"extras_amount":  entity.ExtrasAmount,
		"extras_note":    entity.ExtrasNote,
		"deposit_amount": entity.DepositAmount,
		"final_price":    entity.FinalPrice,
		"notes":          entity.Notes,
		"penalty_amount": entity.PenaltyAmount,
		"penalty_note":   entity.PenaltyNote,
	}).Error
}

func (r *AppointmentRepository) FindByID(id uint) (*appointment.Appointment, error) {
	var m models.AppointmentModel
	if err := r.db.Preload("Client").Preload("Service").First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toAppointmentDomain(&m), nil
}

func (r *AppointmentRepository) FindAll(filters appointment.AppointmentFilters) ([]appointment.Appointment, error) {
	var rows []models.AppointmentModel
	q := r.db.Preload("Client").Preload("Service")

	if filters.Status != nil {
		q = q.Where("status = ?", string(*filters.Status))
	}
	if filters.ClientID != nil {
		q = q.Where("client_id = ?", *filters.ClientID)
	}
	if filters.DateFrom != nil {
		q = q.Where("date >= ?", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		q = q.Where("date <= ?", *filters.DateTo)
	}

	if err := q.Order("date ASC, start_time ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]appointment.Appointment, len(rows))
	for i, m := range rows {
		result[i] = *toAppointmentDomain(&m)
	}
	return result, nil
}

func (r *AppointmentRepository) GetAllClientStats() (map[uint]appointment.ClientStats, error) {
	type row struct {
		ClientID         uint
		AppointmentCount int64
		TotalSpent       float64
	}
	var rows []row
	err := r.db.Model(&models.AppointmentModel{}).
		Select(`client_id,
			COUNT(*) FILTER (WHERE status != 'CANCELLED') AS appointment_count,
			COALESCE(SUM(CASE WHEN status = 'DONE' THEN final_price ELSE 0 END), 0) AS total_spent`).
		Group("client_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make(map[uint]appointment.ClientStats, len(rows))
	for _, row := range rows {
		result[row.ClientID] = appointment.ClientStats{
			AppointmentCount: row.AppointmentCount,
			TotalSpent:       row.TotalSpent,
		}
	}
	return result, nil
}

func (r *AppointmentRepository) FindByDateRange(from, to string) ([]appointment.Appointment, error) {
	var rows []models.AppointmentModel
	if err := r.db.Preload("Client").Preload("Service").
		Where("date >= ? AND date <= ?", from, to).
		Order("date ASC, start_time ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]appointment.Appointment, len(rows))
	for i, m := range rows {
		result[i] = *toAppointmentDomain(&m)
	}
	return result, nil
}

// FindNextByClient devuelve el próximo turno activo (no cancelado/ausente) del cliente a partir de `from`.
func (r *AppointmentRepository) FindNextByClient(clientID uint, from string) (*appointment.Appointment, error) {
	var m models.AppointmentModel
	err := r.db.Preload("Client").Preload("Service").
		Where("client_id = ? AND date >= ? AND status NOT IN ?", clientID, from, []string{"CANCELLED", "ABSENT"}).
		Order("date ASC, start_time ASC").
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toAppointmentDomain(&m), nil
}

// ExistsByClientAndDate devuelve true si el cliente ya tiene un turno activo (no cancelado/ausente) en esa fecha.
func (r *AppointmentRepository) ExistsByClientAndDate(clientID uint, date string) (bool, error) {
	var count int64
	err := r.db.Model(&models.AppointmentModel{}).
		Where("client_id = ? AND date = ? AND status NOT IN ?", clientID, date, []string{"CANCELLED", "ABSENT"}).
		Count(&count).Error
	return count > 0, err
}

// ExistsOverlap verifica si hay un turno que se solapa con el rango [startTime, endTime)
// en la fecha dada. Los turnos cancelados o ausentes se ignoran.
func (r *AppointmentRepository) ExistsOverlap(date, startTime, endTime string, professionalID *uint) (bool, error) {
	q := r.db.Model(&models.AppointmentModel{}).
		Where("date = ? AND start_time < ? AND end_time > ? AND status NOT IN ?",
			date, endTime, startTime, []string{"CANCELLED", "ABSENT"})

	if professionalID != nil {
		q = q.Where("professional_id = ?", *professionalID)
	}

	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *AppointmentRepository) GetMetrics(today, weekStart, monthStart, yearStart string) (*appointment.DashboardMetrics, error) {
	var row struct {
		// Hoy
		TodayTotal     int64
		TodayCompleted int64
		TodayCancelled int64
		TodayPending   int64
		TodayConfirmed int64
		TodayRevenue   float64
		TodayDeposits  float64
		// Semana
		WeekTotal     int64
		WeekCompleted int64
		WeekCancelled int64
		WeekPending   int64
		WeekConfirmed int64
		WeekRevenue   float64
		WeekDeposits  float64
		// Mes
		MonthTotal     int64
		MonthCompleted int64
		MonthCancelled int64
		MonthPending   int64
		MonthConfirmed int64
		MonthRevenue   float64
		MonthDeposits  float64
		// Año
		YearTotal     int64
		YearCompleted int64
		YearCancelled int64
		YearPending   int64
		YearConfirmed int64
		YearRevenue   float64
		YearDeposits  float64
	}

	err := r.db.Raw(`
		SELECT
			COUNT(*) FILTER (WHERE date = ?)                                                       AS today_total,
			COUNT(*) FILTER (WHERE date = ? AND status = 'DONE')                                   AS today_completed,
			COUNT(*) FILTER (WHERE date = ? AND status = 'CANCELLED')                              AS today_cancelled,
			COUNT(*) FILTER (WHERE date = ? AND status = 'PENDING')                               AS today_pending,
			COUNT(*) FILTER (WHERE date = ? AND status = 'CONFIRMED')                             AS today_confirmed,
			COALESCE(SUM(final_price)    FILTER (WHERE date = ? AND status = 'DONE'), 0)           AS today_revenue,
			COALESCE(SUM(deposit_amount) FILTER (WHERE date = ? AND status IN ('CONFIRMED','DONE')), 0) AS today_deposits,

			COUNT(*) FILTER (WHERE date >= ? AND date <= ?)                                        AS week_total,
			COUNT(*) FILTER (WHERE date >= ? AND date <= ? AND status = 'DONE')                    AS week_completed,
			COUNT(*) FILTER (WHERE date >= ? AND date <= ? AND status = 'CANCELLED')               AS week_cancelled,
			COUNT(*) FILTER (WHERE date >= ? AND date <= ? AND status = 'PENDING')                 AS week_pending,
			COUNT(*) FILTER (WHERE date >= ? AND date <= ? AND status = 'CONFIRMED')               AS week_confirmed,
			COALESCE(SUM(final_price)    FILTER (WHERE date >= ? AND date <= ? AND status = 'DONE'), 0) AS week_revenue,
			COALESCE(SUM(deposit_amount) FILTER (WHERE date >= ? AND date <= ? AND status IN ('CONFIRMED','DONE')), 0) AS week_deposits,

			COUNT(*) FILTER (WHERE date >= ?)                                                      AS month_total,
			COUNT(*) FILTER (WHERE date >= ? AND status = 'DONE')                                  AS month_completed,
			COUNT(*) FILTER (WHERE date >= ? AND status = 'CANCELLED')                             AS month_cancelled,
			COUNT(*) FILTER (WHERE date >= ? AND status = 'PENDING')                              AS month_pending,
			COUNT(*) FILTER (WHERE date >= ? AND status = 'CONFIRMED')                            AS month_confirmed,
			COALESCE(SUM(final_price)    FILTER (WHERE date >= ? AND status = 'DONE'), 0)           AS month_revenue,
			COALESCE(SUM(deposit_amount) FILTER (WHERE date >= ? AND status IN ('CONFIRMED','DONE')), 0) AS month_deposits,

			COUNT(*) FILTER (WHERE date >= ?)                                                      AS year_total,
			COUNT(*) FILTER (WHERE date >= ? AND status = 'DONE')                                  AS year_completed,
			COUNT(*) FILTER (WHERE date >= ? AND status = 'CANCELLED')                             AS year_cancelled,
			COUNT(*) FILTER (WHERE date >= ? AND status = 'PENDING')                              AS year_pending,
			COUNT(*) FILTER (WHERE date >= ? AND status = 'CONFIRMED')                            AS year_confirmed,
			COALESCE(SUM(final_price)    FILTER (WHERE date >= ? AND status = 'DONE'), 0)           AS year_revenue,
			COALESCE(SUM(deposit_amount) FILTER (WHERE date >= ? AND status IN ('CONFIRMED','DONE')), 0) AS year_deposits
		FROM appointments
		WHERE deleted_at IS NULL`,
		// today (7)
		today, today, today, today, today, today, today,
		// week (2 params cada filtro × 7 = 14)
		weekStart, today,
		weekStart, today,
		weekStart, today,
		weekStart, today,
		weekStart, today,
		weekStart, today,
		weekStart, today,
		// month (7)
		monthStart, monthStart, monthStart, monthStart, monthStart, monthStart, monthStart,
		// year (7)
		yearStart, yearStart, yearStart, yearStart, yearStart, yearStart, yearStart,
	).Scan(&row).Error
	if err != nil {
		return nil, err
	}

	return &appointment.DashboardMetrics{
		Today: appointment.PeriodStats{
			Total: row.TodayTotal, Completed: row.TodayCompleted, Cancelled: row.TodayCancelled,
			Pending: row.TodayPending, Confirmed: row.TodayConfirmed,
			Revenue: row.TodayRevenue, Deposits: row.TodayDeposits,
		},
		Week: appointment.PeriodStats{
			Total: row.WeekTotal, Completed: row.WeekCompleted, Cancelled: row.WeekCancelled,
			Pending: row.WeekPending, Confirmed: row.WeekConfirmed,
			Revenue: row.WeekRevenue, Deposits: row.WeekDeposits,
		},
		Month: appointment.PeriodStats{
			Total: row.MonthTotal, Completed: row.MonthCompleted, Cancelled: row.MonthCancelled,
			Pending: row.MonthPending, Confirmed: row.MonthConfirmed,
			Revenue: row.MonthRevenue, Deposits: row.MonthDeposits,
		},
		Year: appointment.PeriodStats{
			Total: row.YearTotal, Completed: row.YearCompleted, Cancelled: row.YearCancelled,
			Pending: row.YearPending, Confirmed: row.YearConfirmed,
			Revenue: row.YearRevenue, Deposits: row.YearDeposits,
		},
	}, nil
}

func toAppointmentDomain(m *models.AppointmentModel) *appointment.Appointment {
	return &appointment.Appointment{
		ID:             m.ID,
		ClientID:       m.ClientID,
		ServiceID:      m.ServiceID,
		ProfessionalID: m.ProfessionalID,
		Date:           m.Date,
		StartTime:      m.StartTime,
		EndTime:        m.EndTime,
		BasePrice:      m.BasePrice,
		ExtrasAmount:   m.ExtrasAmount,
		ExtrasNote:     m.ExtrasNote,
		DepositAmount:  m.DepositAmount,
		FinalPrice:     m.FinalPrice,
		Status:         appointment.AppointmentStatus(m.Status),
		Notes:          m.Notes,
		PenaltyAmount:  m.PenaltyAmount,
		PenaltyNote:    m.PenaltyNote,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}
