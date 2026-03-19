package models

import "gorm.io/gorm"

// WeeklyScheduleModel almacena la configuración horaria por día de la semana.
// DayOfWeek: 0=Domingo … 6=Sábado.
type WeeklyScheduleModel struct {
	gorm.Model
	DayOfWeek       int    `gorm:"uniqueIndex;not null"`
	Enabled         bool   `gorm:"default:true"`
	OpeningTime     string `gorm:"type:varchar(5);not null"` // "09:00"
	ClosingTime     string `gorm:"type:varchar(5);not null"` // "19:00"
	SlotDurationMin int    `gorm:"default:30"`
}

func (WeeklyScheduleModel) TableName() string { return "weekly_schedules" }

// BlockedSlotModel representa un bloqueo puntual o recurrente en la agenda.
type BlockedSlotModel struct {
	gorm.Model
	Date      string `gorm:"type:varchar(10);index;not null"` // "2006-01-02"
	StartTime string `gorm:"type:varchar(5);not null"`        // "14:00"
	EndTime   string `gorm:"type:varchar(5);not null"`        // "15:00"
	Reason    string
	Permanent bool `gorm:"default:false"`
}

func (BlockedSlotModel) TableName() string { return "blocked_slots" }
