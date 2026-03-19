package clientapp

import (
	"apiGoShei/internal/domain/appointment"
	"apiGoShei/internal/domain/client"
)

// ClientWithStats combina los datos del cliente con sus acumulados de turnos.
type ClientWithStats struct {
	client.Client
	AppointmentCount int64
	TotalSpent       float64
}

type ListClientsUseCase struct {
	repo     client.Repository
	apptRepo appointment.Repository
}

func NewListClientsUseCase(repo client.Repository, apptRepo appointment.Repository) *ListClientsUseCase {
	return &ListClientsUseCase{repo: repo, apptRepo: apptRepo}
}

func (uc *ListClientsUseCase) Execute() ([]ClientWithStats, error) {
	clients, err := uc.repo.FindAll()
	if err != nil {
		return nil, err
	}
	stats, err := uc.apptRepo.GetAllClientStats()
	if err != nil {
		return nil, err
	}
	result := make([]ClientWithStats, len(clients))
	for i, c := range clients {
		s := stats[c.ID]
		result[i] = ClientWithStats{
			Client:           c,
			AppointmentCount: s.AppointmentCount,
			TotalSpent:       s.TotalSpent,
		}
	}
	return result, nil
}
