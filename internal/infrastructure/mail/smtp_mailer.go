package mail

import (
	"fmt"
	"net/smtp"
	"strings"
)

// SMTPMailer envía mails reales usando autenticación SMTP plain.
// Para Gmail se recomienda usar una "App Password" (no la contraseña principal).
type SMTPMailer struct {
	host       string
	port       string
	user       string
	password   string
	adminEmail string
	appName    string
	ctaURL     string
}

func NewSMTPMailer(host, port, user, password, adminEmail string) *SMTPMailer {
	return &SMTPMailer{
		host:       host,
		port:       port,
		user:       user,
		password:   password,
		adminEmail: adminEmail,
		appName:    "Nails Deni",
		ctaURL:     "",
	}
}

func (m *SMTPMailer) SendNewAppointmentToClient(data AppointmentData) error {
	emailData := AppointmentEmailData{
		ClientName:    data.ClientName,
		ServiceName:   data.ServiceName,
		Date:          data.Date,
		StartTime:     data.StartTime,
		EndTime:       data.EndTime,
		BasePrice:     fmt.Sprintf("%.0f", data.BasePrice),
		ExtrasAmount:  formatIfNonZero(data.ExtrasAmount),
		ExtrasNote:    data.ExtrasNote,
		FinalPrice:    fmt.Sprintf("%.0f", data.FinalPrice),
		DepositAmount: formatIfNonZero(data.DepositAmount),
		Status:        data.Status,
		Notes:         data.Notes,
		AppName:       m.appName,
		CTAUrl:        m.ctaURL,
	}
	html, err := BuildClientAppointmentEmail(emailData)
	if err != nil {
		return err
	}
	return m.sendHTML(data.ClientEmail, "✅ Tu turno fue registrado — "+m.appName, html)
}

func (m *SMTPMailer) SendNewAppointmentToAdmin(data AppointmentData) error {
	emailData := AppointmentEmailData{
		ClientName:    data.ClientName,
		ServiceName:   data.ServiceName,
		Date:          data.Date,
		StartTime:     data.StartTime,
		EndTime:       data.EndTime,
		BasePrice:     fmt.Sprintf("%.0f", data.BasePrice),
		ExtrasAmount:  formatIfNonZero(data.ExtrasAmount),
		ExtrasNote:    data.ExtrasNote,
		FinalPrice:    fmt.Sprintf("%.0f", data.FinalPrice),
		DepositAmount: formatIfNonZero(data.DepositAmount),
		Status:        data.Status,
		Notes:         data.Notes,
		AppName:       m.appName,
		CTAUrl:        m.ctaURL,
	}
	html, err := BuildAdminAppointmentEmail(emailData)
	if err != nil {
		return err
	}
	return m.sendHTML(m.adminEmail, "🔔 Nuevo turno — "+data.ClientName+" | "+m.appName, html)
}

func (m *SMTPMailer) sendHTML(to, subject, htmlBody string) error {
	auth := smtp.PlainAuth("", m.user, m.password, m.host)
	msg := strings.Join([]string{
		"From: " + m.appName + " <" + m.user + ">",
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		"",
		htmlBody,
	}, "\r\n")
	return smtp.SendMail(m.host+":"+m.port, auth, m.user, []string{to}, []byte(msg))
}

func formatIfNonZero(v float64) string {
	if v == 0 {
		return ""
	}
	return fmt.Sprintf("%.0f", v)
}
