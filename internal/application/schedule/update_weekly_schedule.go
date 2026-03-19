package scheduleapp

import "apiGoShei/internal/domain/schedule"

type WeeklyScheduleEntry struct {
	DayOfWeek       int
	Enabled         bool
	OpeningTime     string
	ClosingTime     string
	SlotDurationMin int
}

type UpdateWeeklyScheduleUseCase struct {
	scheduleRepo schedule.Repository
}

func NewUpdateWeeklyScheduleUseCase(repo schedule.Repository) *UpdateWeeklyScheduleUseCase {
	return &UpdateWeeklyScheduleUseCase{scheduleRepo: repo}
}

func (uc *UpdateWeeklyScheduleUseCase) Execute(entries []WeeklyScheduleEntry) ([]schedule.WeeklySchedule, error) {
	schedules := make([]schedule.WeeklySchedule, len(entries))
	for i, e := range entries {
		schedules[i] = schedule.WeeklySchedule{
			DayOfWeek:       e.DayOfWeek,
			Enabled:         e.Enabled,
			OpeningTime:     e.OpeningTime,
			ClosingTime:     e.ClosingTime,
			SlotDurationMin: e.SlotDurationMin,
		}
	}
	if err := uc.scheduleRepo.UpsertWeeklySchedule(schedules); err != nil {
		return nil, err
	}
	return uc.scheduleRepo.FindWeeklySchedule()
}

// ─── GetWeeklySchedule ─────────────────────────────────────────────────────────────────

type GetWeeklyScheduleUseCase struct {
	scheduleRepo schedule.Repository
}

func NewGetWeeklyScheduleUseCase(repo schedule.Repository) *GetWeeklyScheduleUseCase {
	return &GetWeeklyScheduleUseCase{scheduleRepo: repo}
}

func (uc *GetWeeklyScheduleUseCase) Execute() ([]schedule.WeeklySchedule, error) {
	return uc.scheduleRepo.FindWeeklySchedule()
}
