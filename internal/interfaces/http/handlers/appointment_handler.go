package handlers

import (
	"net/http"
	"strconv"
	"time"

	appointmentapp "apiGoShei/internal/application/appointment"
	"apiGoShei/internal/domain/appointment"
	"apiGoShei/internal/infrastructure/security"
	"apiGoShei/internal/interfaces/http/dto"
	"apiGoShei/internal/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	createUC        *appointmentapp.CreateAppointmentUseCase
	cancelUC        *appointmentapp.CancelAppointmentUseCase
	confirmUC       *appointmentapp.ConfirmAppointmentUseCase
	completeUC      *appointmentapp.CompleteAppointmentUseCase
	listUC          *appointmentapp.ListAppointmentUseCase
	getByIDUC       *appointmentapp.GetAppointmentByIDUseCase
	calendarUC      *appointmentapp.ListCalendarAppointmentUseCase
	finalPriceUC    *appointmentapp.UpdateFinalPriceUseCase
	updateDepositUC *appointmentapp.UpdateDepositUseCase
	nextUC          *appointmentapp.NextAppointmentUseCase
}

func NewAppointmentHandler(
	createUC *appointmentapp.CreateAppointmentUseCase,
	cancelUC *appointmentapp.CancelAppointmentUseCase,
	confirmUC *appointmentapp.ConfirmAppointmentUseCase,
	completeUC *appointmentapp.CompleteAppointmentUseCase,
	listUC *appointmentapp.ListAppointmentUseCase,
	getByIDUC *appointmentapp.GetAppointmentByIDUseCase,
	calendarUC *appointmentapp.ListCalendarAppointmentUseCase,
	finalPriceUC *appointmentapp.UpdateFinalPriceUseCase,
	updateDepositUC *appointmentapp.UpdateDepositUseCase,
	nextUC *appointmentapp.NextAppointmentUseCase,
) *AppointmentHandler {
	return &AppointmentHandler{
		createUC:        createUC,
		cancelUC:        cancelUC,
		confirmUC:       confirmUC,
		completeUC:      completeUC,
		listUC:          listUC,
		getByIDUC:       getByIDUC,
		calendarUC:      calendarUC,
		finalPriceUC:    finalPriceUC,
		updateDepositUC: updateDepositUC,
		nextUC:          nextUC,
	}
}

func (h *AppointmentHandler) Create(c *gin.Context) {
	var req dto.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := c.MustGet(middleware.ClaimsKey).(*security.Claims)
	clientID := req.ClientID
	if claims.Role == "client" {

		clientID = claims.UserID
	} else if clientID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id es requerido"})
		return
	}

	result, err := h.createUC.Execute(appointmentapp.CreateAppointmentInput{
		ClientID:       clientID,
		ServiceID:      req.ServiceID,
		ProfessionalID: req.ProfessionalID,
		Date:           req.Date,
		StartTime:      req.StartTime,
		Notes:          req.Notes,
		IsAdmin:        claims.Role == "admin" || claims.Role == "superadmin",
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dto.AppointmentToResponse(result))
}

func (h *AppointmentHandler) NextAppointment(c *gin.Context) {
	claims := c.MustGet(middleware.ClaimsKey).(*security.Claims)
	appt, err := h.nextUC.Execute(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener próximo turno"})
		return
	}
	if appt == nil {
		c.JSON(http.StatusOK, gin.H{"next_appointment": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"next_appointment": dto.AppointmentToResponse(appt)})
}

func (h *AppointmentHandler) List(c *gin.Context) {
	filters := appointment.AppointmentFilters{}
	if s := c.Query("status"); s != "" {
		st := appointment.AppointmentStatus(s)
		filters.Status = &st
	}
	if cid := c.Query("client_id"); cid != "" {
		if id, err := strconv.ParseUint(cid, 10, 64); err == nil {
			uid := uint(id)
			filters.ClientID = &uid
		}
	}
	if df := c.Query("date_from"); df != "" {
		filters.DateFrom = &df
	} else {

		defaultFrom := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
		filters.DateFrom = &defaultFrom
	}
	if dt := c.Query("date_to"); dt != "" {
		filters.DateTo = &dt
	}
	appointments, err := h.listUC.Execute(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al listar turnos"})
		return
	}
	resp := make([]dto.AppointmentResponse, len(appointments))
	for i, a := range appointments {
		a := a
		resp[i] = dto.AppointmentToResponse(&a)
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AppointmentHandler) ListMine(c *gin.Context) {
	claims := c.MustGet(middleware.ClaimsKey).(*security.Claims)
	clientID := claims.UserID

	filters := appointment.AppointmentFilters{ClientID: &clientID}
	if s := c.Query("status"); s != "" {
		st := appointment.AppointmentStatus(s)
		filters.Status = &st
	}
	appointments, err := h.listUC.Execute(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al listar turnos"})
		return
	}
	resp := make([]dto.AppointmentResponse, len(appointments))
	for i, a := range appointments {
		a := a
		resp[i] = dto.AppointmentToResponse(&a)
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AppointmentHandler) ListByRole(c *gin.Context) {
	claims := c.MustGet(middleware.ClaimsKey).(*security.Claims)
	if claims.Role == "client" {
		h.ListMine(c)
		return
	}
	h.List(c)
}

func (h *AppointmentHandler) Calendar(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "los parametros 'from' y 'to' son requeridos"})
		return
	}
	appointments, err := h.calendarUC.Execute(appointmentapp.ListCalendarAppointmentInput{From: from, To: to})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener calendario"})
		return
	}
	resp := make([]dto.AppointmentResponse, len(appointments))
	for i, a := range appointments {
		a := a
		resp[i] = dto.AppointmentToResponse(&a)
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AppointmentHandler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}
	appt, err := h.getByIDUC.Execute(id)
	if err != nil {
		if err == appointment.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener turno"})
		return
	}
	c.JSON(http.StatusOK, dto.AppointmentToResponse(appt))
}

func (h *AppointmentHandler) Confirm(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}
	appt, err := h.confirmUC.Execute(id)
	if err != nil {
		code := http.StatusBadRequest
		if err == appointment.ErrNotFound {
			code = http.StatusNotFound
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.AppointmentToResponse(appt))
}

func (h *AppointmentHandler) Cancel(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}
	appt, err := h.cancelUC.Execute(id)
	if err != nil {
		code := http.StatusBadRequest
		if err == appointment.ErrNotFound {
			code = http.StatusNotFound
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.AppointmentToResponse(appt))
}

func (h *AppointmentHandler) Complete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}
	appt, err := h.completeUC.Execute(id)
	if err != nil {
		code := http.StatusBadRequest
		if err == appointment.ErrNotFound {
			code = http.StatusNotFound
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.AppointmentToResponse(appt))
}

func (h *AppointmentHandler) UpdateFinalPrice(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}
	var req dto.UpdateFinalPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	appt, err := h.finalPriceUC.Execute(appointmentapp.UpdateFinalPriceInput{
		ID:           id,
		ExtrasAmount: req.ExtrasAmount,
		ExtrasNote:   req.ExtrasNote,
	})
	if err != nil {
		code := http.StatusBadRequest
		if err == appointment.ErrNotFound {
			code = http.StatusNotFound
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.AppointmentToResponse(appt))
}

func (h *AppointmentHandler) UpdateDeposit(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}
	var req dto.UpdateDepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	appt, err := h.updateDepositUC.Execute(appointmentapp.UpdateDepositInput{
		ID:            id,
		DepositAmount: req.DepositAmount,
	})
	if err != nil {
		code := http.StatusBadRequest
		if err == appointment.ErrNotFound {
			code = http.StatusNotFound
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.AppointmentToResponse(appt))
}

func parseID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	return uint(id), err
}
