// Package mail contiene la interfaz Mailer y sus implementaciones.
// La interfaz es implementada tanto por SMTPMailer (producción) como por
// MockMailer (desarrollo), permitiendo intercambiarlas sin cambiar la lógica
// de negocio.
package mail

// AppointmentData agrupa los datos necesarios para los templates de mail de turno.
type AppointmentData struct {
	ClientName    string
	ClientEmail   string
	ClientPhone   string
	ServiceName   string
	Date          string
	StartTime     string
	EndTime       string
	Status        string
	Notes         string
	BasePrice     float64
	ExtrasAmount  float64
	ExtrasNote    string
	FinalPrice    float64
	DepositAmount float64
}

// Mailer define el contrato para el envío de notificaciones por mail.
// Cualquier implementación (SMTP, SendGrid, mock, etc.) debe satisfacer esta interfaz.
type Mailer interface {
	SendNewAppointmentToClient(data AppointmentData) error
	SendNewAppointmentToAdmin(data AppointmentData) error
}
