package appointment

type AppointmentFilters struct {
	Status   *AppointmentStatus
	ClientID *uint
	DateFrom *string
	DateTo   *string
}

type ClientStats struct {
	AppointmentCount int64
	TotalSpent       float64
}

type PeriodStats struct {
	Total     int64
	Completed int64
	Cancelled int64
	Pending   int64
	Confirmed int64
	Revenue   float64
	Deposits  float64
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
	GetAllClientStats() (map[uint]ClientStats, error)
}
