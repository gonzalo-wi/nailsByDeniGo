package dto

type WeeklyScheduleEntryRequest struct {
	DayOfWeek       int    `json:"day_of_week"        binding:"min=0,max=6"`
	Enabled         bool   `json:"enabled"`
	OpeningTime     string `json:"opening_time"       binding:"required"` // "09:00"
	ClosingTime     string `json:"closing_time"       binding:"required"` // "19:00"
	SlotDurationMin int    `json:"slot_duration_min"  binding:"required,min=10"`
}

type UpdateWeeklyScheduleRequest struct {
	Schedule []WeeklyScheduleEntryRequest `json:"schedule" binding:"required"`
}

type BlockedSlotRequest struct {
	Date      string `json:"date"       binding:"required"` // "2026-03-16"
	StartTime string `json:"start_time" binding:"required"` // "14:00"
	EndTime   string `json:"end_time"   binding:"required"` // "15:00"
	Reason    string `json:"reason"`
	Permanent bool   `json:"permanent"`
}
