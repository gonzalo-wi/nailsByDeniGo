package scheduleapp

import "apiGoShei/internal/domain/schedule"

type BlockTimeSlotInput struct {
	Date      string
	StartTime string
	EndTime   string
	Reason    string
	Permanent bool
}

type BlockTimeSlotUseCase struct {
	scheduleRepo schedule.Repository
}

func NewBlockTimeSlotUseCase(repo schedule.Repository) *BlockTimeSlotUseCase {
	return &BlockTimeSlotUseCase{scheduleRepo: repo}
}

func (uc *BlockTimeSlotUseCase) Execute(input BlockTimeSlotInput) (*schedule.BlockedSlot, error) {
	slot := &schedule.BlockedSlot{
		Date:      input.Date,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		Reason:    input.Reason,
		Permanent: input.Permanent,
	}
	if err := uc.scheduleRepo.CreateBlockedSlot(slot); err != nil {
		return nil, err
	}
	return slot, nil
}
