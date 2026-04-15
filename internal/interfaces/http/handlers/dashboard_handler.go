package handlers

import (
	"net/http"

	dashboardapp "apiGoShei/internal/application/dashboard"
	"apiGoShei/internal/interfaces/http/dto"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	getMetricsUC *dashboardapp.GetMetricsUseCase
}

func NewDashboardHandler(getMetricsUC *dashboardapp.GetMetricsUseCase) *DashboardHandler {
	return &DashboardHandler{getMetricsUC: getMetricsUC}
}

func (h *DashboardHandler) GetMetrics(c *gin.Context) {
	metrics, err := h.getMetricsUC.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener metricas"})
		return
	}
	c.JSON(http.StatusOK, dto.MetricsToResponse(metrics))
}
