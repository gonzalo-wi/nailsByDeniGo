package handlers

import (
	"net/http"

	clientapp "apiGoShei/internal/application/client"
	"apiGoShei/internal/interfaces/http/dto"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	listUC *clientapp.ListClientsUseCase
}

func NewClientHandler(listUC *clientapp.ListClientsUseCase) *ClientHandler {
	return &ClientHandler{listUC: listUC}
}

// GET /clients
func (h *ClientHandler) List(c *gin.Context) {
	clients, err := h.listUC.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al listar clientes"})
		return
	}
	resp := make([]dto.ClientResponse, len(clients))
	for i, cl := range clients {
		resp[i] = dto.ClientResponse{
			ID:               cl.ID,
			FirstName:        cl.FirstName,
			LastName:         cl.LastName,
			Email:            cl.Email,
			Phone:            cl.Phone,
			Active:           cl.Active,
			AppointmentCount: cl.AppointmentCount,
			TotalSpent:       cl.TotalSpent,
		}
	}
	c.JSON(http.StatusOK, resp)
}
