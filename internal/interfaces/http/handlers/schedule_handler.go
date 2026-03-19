package handlers

import (
	"fmt"
	"net/http"

	scheduleapp "apiGoShei/internal/application/schedule"
	"apiGoShei/internal/interfaces/http/dto"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	getWeeklyUC    *scheduleapp.GetWeeklyScheduleUseCase
	updateUC       *scheduleapp.UpdateWeeklyScheduleUseCase
	blockUC        *scheduleapp.BlockTimeSlotUseCase
	availabilityUC *scheduleapp.GetAvailabilityUseCase
}

func NewScheduleHandler(
	getWeeklyUC *scheduleapp.GetWeeklyScheduleUseCase,
	updateUC *scheduleapp.UpdateWeeklyScheduleUseCase,
	blockUC *scheduleapp.BlockTimeSlotUseCase,
	availabilityUC *scheduleapp.GetAvailabilityUseCase,
) *ScheduleHandler {
	return &ScheduleHandler{
		getWeeklyUC:    getWeeklyUC,
		updateUC:       updateUC,
		blockUC:        blockUC,
		availabilityUC: availabilityUC,
	}
}

// GET /schedule/weekly
func (h *ScheduleHandler) GetWeekly(c *gin.Context) {
	schedules, err := h.getWeeklyUC.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener horarios"})
		return
	}
	resp := make([]dto.WeeklyScheduleResponse, len(schedules))
	for i, s := range schedules {
		resp[i] = dto.WeeklyScheduleToResponse(s)
	}
	c.JSON(http.StatusOK, resp)
}

// PUT /schedule/weekly
func (h *ScheduleHandler) UpdateWeekly(c *gin.Context) {
	var req dto.UpdateWeeklyScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entries := make([]scheduleapp.WeeklyScheduleEntry, len(req.Schedule))
	for i, e := range req.Schedule {
		entries[i] = scheduleapp.WeeklyScheduleEntry{
			DayOfWeek:       e.DayOfWeek,
			Enabled:         e.Enabled,
			OpeningTime:     e.OpeningTime,
			ClosingTime:     e.ClosingTime,
			SlotDurationMin: e.SlotDurationMin,
		}
	}
	updated, err := h.updateUC.Execute(entries)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar horarios"})
		return
	}
	resp := make([]dto.WeeklyScheduleResponse, len(updated))
	for i, s := range updated {
		resp[i] = dto.WeeklyScheduleToResponse(s)
	}
	c.JSON(http.StatusOK, resp)
}

// POST /schedule/blocked-slots
func (h *ScheduleHandler) BlockSlot(c *gin.Context) {
	var req dto.BlockedSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	slot, err := h.blockUC.Execute(scheduleapp.BlockTimeSlotInput{
		Date:      req.Date,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Reason:    req.Reason,
		Permanent: req.Permanent,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al bloquear horario"})
		return
	}
	c.JSON(http.StatusCreated, dto.BlockedSlotToResponse(*slot))
}

// GET /schedule/availability?date=2026-03-16&duration=60
func (h *ScheduleHandler) GetAvailability(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el parametro 'date' es requerido"})
		return
	}
	var duration int
	if d := c.Query("duration"); d != "" {
		if _, err := parseDuration(d, &duration); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "duracion invalida"})
			return
		}
	}
	out, err := h.availabilityUC.Execute(scheduleapp.GetAvailabilityInput{
		Date:            date,
		DurationMinutes: duration,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al consultar disponibilidad"})
		return
	}
	c.JSON(http.StatusOK, out)
}

func parseDuration(s string, out *int) (int, error) {
	var n int
	_, err := fmt.Sscan(s, &n)
	if err == nil {
		*out = n
	}
	return n, err
}
