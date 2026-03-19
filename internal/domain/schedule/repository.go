package schedule

type Repository interface {
	FindWeeklySchedule() ([]WeeklySchedule, error)
	UpsertWeeklySchedule(schedules []WeeklySchedule) error
	CreateBlockedSlot(slot *BlockedSlot) error
	FindBlockedSlots(date string) ([]BlockedSlot, error)
	IsAvailable(date, startTime string, durationMinutes int) (bool, error)
}
