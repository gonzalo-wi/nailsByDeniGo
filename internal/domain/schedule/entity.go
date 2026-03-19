package schedule

import "time"

// WeeklySchedule representa la configuración horaria de un día de la semana.
// DayOfWeek: 0 = Domingo, 1 = Lunes … 6 = Sábado (convención Go/time.Weekday).
type WeeklySchedule struct {
	ID              uint
	DayOfWeek       int
	Enabled         bool
	OpeningTime     string // "09:00"
	ClosingTime     string // "19:00"
	SlotDurationMin int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// BlockedSlot representa un bloqueo puntual en la agenda.
type BlockedSlot struct {
	ID        uint
	Date      string // "2006-01-02"
	StartTime string // "14:00"
	EndTime   string // "15:00"
	Reason    string
	Permanent bool
	CreatedAt time.Time
}
