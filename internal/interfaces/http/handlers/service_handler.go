package handlers

import (
	"net/http"
	"strconv"

	serviceapp "apiGoShei/internal/application/service"
	"apiGoShei/internal/domain/service"
	"apiGoShei/internal/interfaces/http/dto"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	createUC  *serviceapp.CreateServiceUseCase
	listUC    *serviceapp.ListServiceUseCase
	getByIDUC *serviceapp.GetServiceByIDUseCase
	updateUC  *serviceapp.UpdateServiceUseCase
	toggleUC  *serviceapp.ToggleServiceUseCase
}

func NewServiceHandler(
	createUC *serviceapp.CreateServiceUseCase,
	listUC *serviceapp.ListServiceUseCase,
	getByIDUC *serviceapp.GetServiceByIDUseCase,
	updateUC *serviceapp.UpdateServiceUseCase,
	toggleUC *serviceapp.ToggleServiceUseCase,
) *ServiceHandler {
	return &ServiceHandler{
		createUC:  createUC,
		listUC:    listUC,
		getByIDUC: getByIDUC,
		updateUC:  updateUC,
		toggleUC:  toggleUC,
	}
}

// POST /services
func (h *ServiceHandler) Create(c *gin.Context) {
	var req dto.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s, err := h.createUC.Execute(serviceapp.CreateServiceInput{
		Name:             req.Name,
		Description:      req.Description,
		DurationMinutes:  req.DurationMinutes,
		BasePrice:        req.BasePrice,
		RequiresDeposit:  req.RequiresDeposit,
		SuggestedDeposit: req.SuggestedDeposit,
		Color:            req.Color,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear servicio"})
		return
	}
	c.JSON(http.StatusCreated, dto.ServiceToResponse(s))
}

// GET /services
func (h *ServiceHandler) List(c *gin.Context) {
	activeOnly := c.Query("active") != "false"
	services, err := h.listUC.Execute(activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al listar servicios"})
		return
	}
	resp := make([]dto.ServiceResponse, len(services))
	for i, s := range services {
		s := s
		resp[i] = dto.ServiceToResponse(&s)
	}
	c.JSON(http.StatusOK, resp)
}

// GET /services/:id
func (h *ServiceHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}
	s, err := h.getByIDUC.Execute(uint(id))
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener servicio"})
		return
	}
	c.JSON(http.StatusOK, dto.ServiceToResponse(s))
}

// PATCH /services/:id
func (h *ServiceHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}
	var req dto.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s, err := h.updateUC.Execute(serviceapp.UpdateServiceInput{
		ID:               uint(id),
		Name:             req.Name,
		Description:      req.Description,
		DurationMinutes:  req.DurationMinutes,
		BasePrice:        req.BasePrice,
		RequiresDeposit:  req.RequiresDeposit,
		SuggestedDeposit: req.SuggestedDeposit,
		Color:            req.Color,
	})
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar servicio"})
		return
	}
	c.JSON(http.StatusOK, dto.ServiceToResponse(s))
}

// PATCH /services/:id/toggle
func (h *ServiceHandler) Toggle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}
	s, err := h.toggleUC.Execute(uint(id))
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al cambiar estado del servicio"})
		return
	}
	c.JSON(http.StatusOK, dto.ServiceToResponse(s))
}
