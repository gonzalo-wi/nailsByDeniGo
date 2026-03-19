package mail

import "apiGoShei/internal/infrastructure/logger"

// MockMailer es la implementación de Mailer para entornos de desarrollo.
// En lugar de enviar mails reales, imprime los datos en el log.
type MockMailer struct{}

func NewMockMailer() *MockMailer { return &MockMailer{} }

func (m *MockMailer) SendNewAppointmentToClient(data AppointmentData) error {
	logger.Info.Printf(
		"[MOCK MAIL → CLIENTE] Para: %s <%s> | Servicio: %s | %s %s-%s",
		data.ClientName, data.ClientEmail, data.ServiceName, data.Date, data.StartTime, data.EndTime,
	)
	return nil
}

func (m *MockMailer) SendNewAppointmentToAdmin(data AppointmentData) error {
	logger.Info.Printf(
		"[MOCK MAIL → ADMIN] Cliente: %s | Teléfono: %s | Servicio: %s | %s %s-%s",
		data.ClientName, data.ClientPhone, data.ServiceName, data.Date, data.StartTime, data.EndTime,
	)
	return nil
}
