package dto

type WeeklyScheduleEntryRequest struct {
	DayOfWeek       int    `json:"day_of_week"        binding:"min=0,max=6"`
	Enabled         bool   `json:"enabled"`
	OpeningTime     string `json:"opening_time"       binding:"required"`
	ClosingTime     string `json:"closing_time"       binding:"required"`
	SlotDurationMin int    `json:"slot_duration_min"  binding:"required,min=10"`
}

type UpdateWeeklyScheduleRequest struct {
	Schedule []WeeklyScheduleEntryRequest `json:"schedule" binding:"required"`
}

type BlockedSlotRequest struct {
	Date      string `json:"date"       binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time"   binding:"required"`
	Reason    string `json:"reason"`
	Permanent bool   `json:"permanent"`
}
