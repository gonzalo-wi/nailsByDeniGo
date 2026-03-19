package appointment

type AppointmentFilters struct {
	Status   *AppointmentStatus
	ClientID *uint
	DateFrom *string
	DateTo   *string
}

// ClientStats contiene los acumulados de un cliente calculados desde los turnos.
type ClientStats struct {
	AppointmentCount int64
	TotalSpent       float64
}

// PeriodStats agrupa todas las métricas de un período de tiempo.
type PeriodStats struct {
	Total     int64
	Completed int64   // DONE
	Cancelled int64   // CANCELLED
	Pending   int64   // PENDING
	Confirmed int64   // CONFIRMED
	Revenue   float64 // SUM final_price WHERE DONE
	Deposits  float64 // SUM deposit_amount WHERE CONFIRMED|DONE
}

type DashboardMetrics struct {
	Today PeriodStats
	Week  PeriodStats
	Month PeriodStats
	Year  PeriodStats
}

type Repository interface {
	Create(appointment *Appointment) error
	Update(appointment *Appointment) error
	FindByID(id uint) (*Appointment, error)
	FindAll(filters AppointmentFilters) ([]Appointment, error)
	FindByDateRange(from, to string) ([]Appointment, error)
	ExistsOverlap(date, startTime, endTime string, professionalID *uint) (bool, error)
	ExistsByClientAndDate(clientID uint, date string) (bool, error)
	FindNextByClient(clientID uint, from string) (*Appointment, error)
	GetMetrics(today, weekStart, monthStart, yearStart string) (*DashboardMetrics, error)
	// GetAllClientStats devuelve en una sola query los acumulados por cliente.
	GetAllClientStats() (map[uint]ClientStats, error)
}
