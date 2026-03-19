package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

// AppointmentEmailData contiene los datos necesarios para construir el email de turno.
type AppointmentEmailData struct {
	ClientName    string
	ServiceName   string
	Date          string
	StartTime     string
	EndTime       string
	Professional  string // opcional, vacío si no aplica
	BasePrice     string
	ExtrasAmount  string
	ExtrasNote    string // opcional
	FinalPrice    string
	DepositAmount string
	Status        string
	Notes         string // opcional
	AppName       string
	CTAUrl        string
}

const appointmentHTMLTemplate = `<!DOCTYPE html>
<html lang="es">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>{{.AppName}}</title>
</head>
<body style="margin:0;padding:0;background-color:#f4f4f7;font-family:'Segoe UI',Arial,sans-serif;">

  <!-- Wrapper -->
  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f7;padding:40px 0;">
    <tr>
      <td align="center">

        <!-- Card -->
        <table width="600" cellpadding="0" cellspacing="0"
          style="max-width:600px;width:100%;background-color:#ffffff;border-radius:12px;
                 box-shadow:0 4px 24px rgba(0,0,0,0.08);overflow:hidden;">

          <!-- Header -->
          <tr>
            <td style="background:linear-gradient(135deg,#7c3aed 0%,#a855f7 100%);
                        padding:36px 40px;text-align:center;">
              <p style="margin:0;font-size:13px;color:rgba(255,255,255,0.75);
                         letter-spacing:2px;text-transform:uppercase;font-weight:600;">
                💅 {{.AppName}}
              </p>
              <h1 style="margin:10px 0 0;font-size:26px;font-weight:700;color:#ffffff;line-height:1.3;">
                {{.Title}}
              </h1>
            </td>
          </tr>

          <!-- Greeting -->
          <tr>
            <td style="padding:32px 40px 8px;">
              <p style="margin:0;font-size:16px;color:#374151;line-height:1.6;">
                {{.Greeting}}
              </p>
            </td>
          </tr>

          <!-- Appointment card -->
          <tr>
            <td style="padding:16px 40px 8px;">
              <table width="100%" cellpadding="0" cellspacing="0"
                style="background-color:#faf5ff;border-radius:10px;border:1px solid #e9d5ff;">
                <tr>
                  <td style="padding:24px 28px;">

                    <!-- Fila: Servicio -->
                    <table width="100%" cellpadding="0" cellspacing="0" style="margin-bottom:14px;">
                      <tr>
                        <td width="36" valign="top" style="padding-top:2px;">
                          <span style="font-size:20px;">💅</span>
                        </td>
                        <td>
                          <p style="margin:0;font-size:11px;color:#7c3aed;font-weight:700;
                                     text-transform:uppercase;letter-spacing:1px;">Servicio</p>
                          <p style="margin:3px 0 0;font-size:15px;color:#1f2937;font-weight:600;">
                            {{.ServiceName}}
                          </p>
                        </td>
                      </tr>
                    </table>

                    <!-- Fila: Fecha y hora -->
                    <table width="100%" cellpadding="0" cellspacing="0" style="margin-bottom:14px;">
                      <tr>
                        <td width="36" valign="top" style="padding-top:2px;">
                          <span style="font-size:20px;">📅</span>
                        </td>
                        <td>
                          <p style="margin:0;font-size:11px;color:#7c3aed;font-weight:700;
                                     text-transform:uppercase;letter-spacing:1px;">Fecha y hora</p>
                          <p style="margin:3px 0 0;font-size:15px;color:#1f2937;font-weight:600;">
                            {{.Date}} &nbsp;·&nbsp; {{.StartTime}} – {{.EndTime}}
                          </p>
                        </td>
                      </tr>
                    </table>

                    {{if .Professional}}
                    <!-- Fila: Profesional -->
                    <table width="100%" cellpadding="0" cellspacing="0" style="margin-bottom:14px;">
                      <tr>
                        <td width="36" valign="top" style="padding-top:2px;">
                          <span style="font-size:20px;">👩‍🎨</span>
                        </td>
                        <td>
                          <p style="margin:0;font-size:11px;color:#7c3aed;font-weight:700;
                                     text-transform:uppercase;letter-spacing:1px;">Profesional</p>
                          <p style="margin:3px 0 0;font-size:15px;color:#1f2937;font-weight:600;">
                            {{.Professional}}
                          </p>
                        </td>
                      </tr>
                    </table>
                    {{end}}

                    <!-- Fila: Estado -->
                    <table width="100%" cellpadding="0" cellspacing="0" style="margin-bottom:14px;">
                      <tr>
                        <td width="36" valign="top" style="padding-top:2px;">
                          <span style="font-size:20px;">📋</span>
                        </td>
                        <td>
                          <p style="margin:0;font-size:11px;color:#7c3aed;font-weight:700;
                                     text-transform:uppercase;letter-spacing:1px;">Estado</p>
                          <p style="margin:3px 0 0;font-size:15px;color:#1f2937;font-weight:600;">
                            {{.StatusLabel}}
                          </p>
                        </td>
                      </tr>
                    </table>

                    <!-- Separador -->
                    <hr style="border:none;border-top:1px solid #e9d5ff;margin:18px 0;" />

                    <!-- Precios -->
                    <table width="100%" cellpadding="0" cellspacing="0">
                      <tr>
                        <td style="font-size:13px;color:#6b7280;">Precio base</td>
                        <td align="right" style="font-size:13px;color:#1f2937;font-weight:600;">
                          ${{.BasePrice}}
                        </td>
                      </tr>
                      {{if .ExtrasAmount}}
                      <tr>
                        <td style="font-size:13px;color:#6b7280;padding-top:6px;">
                          Extra{{if .ExtrasNote}}: {{.ExtrasNote}}{{end}}
                        </td>
                        <td align="right" style="font-size:13px;color:#1f2937;font-weight:600;padding-top:6px;">
                          +${{.ExtrasAmount}}
                        </td>
                      </tr>
                      {{end}}
                      <tr>
                        <td colspan="2">
                          <hr style="border:none;border-top:1px solid #e9d5ff;margin:10px 0;" />
                        </td>
                      </tr>
                      <tr>
                        <td style="font-size:15px;color:#1f2937;font-weight:700;">Total</td>
                        <td align="right"
                          style="font-size:18px;color:#7c3aed;font-weight:800;">
                          ${{.FinalPrice}}
                        </td>
                      </tr>
                      {{if .DepositAmount}}
                      <tr>
                        <td style="font-size:12px;color:#6b7280;padding-top:6px;">Seña abonada</td>
                        <td align="right" style="font-size:12px;color:#059669;font-weight:600;padding-top:6px;">
                          ${{.DepositAmount}}
                        </td>
                      </tr>
                      {{end}}
                    </table>

                    {{if .Notes}}
                    <!-- Observaciones -->
                    <div style="margin-top:18px;padding:12px 16px;background:#fff;border-radius:8px;
                                border-left:3px solid #a855f7;">
                      <p style="margin:0;font-size:12px;color:#7c3aed;font-weight:700;
                                 text-transform:uppercase;letter-spacing:1px;">Observaciones</p>
                      <p style="margin:4px 0 0;font-size:13px;color:#374151;">{{.Notes}}</p>
                    </div>
                    {{end}}

                  </td>
                </tr>
              </table>
            </td>
          </tr>

          {{if .CTAUrl}}
          <!-- Botón CTA -->
          <tr>
            <td align="center" style="padding:28px 40px 8px;">
              <a href="{{.CTAUrl}}"
                style="display:inline-block;background:linear-gradient(135deg,#7c3aed,#a855f7);
                        color:#ffffff;text-decoration:none;font-size:15px;font-weight:700;
                        padding:14px 36px;border-radius:50px;
                        box-shadow:0 4px 14px rgba(124,58,237,0.4);">
                Ver mis turnos →
              </a>
            </td>
          </tr>
          {{end}}

          <!-- Mensaje cierre -->
          <tr>
            <td style="padding:24px 40px 12px;">
              <p style="margin:0;font-size:14px;color:#6b7280;line-height:1.6;text-align:center;">
                {{.ClosingMessage}}
              </p>
            </td>
          </tr>

          <!-- Footer -->
          <tr>
            <td style="background-color:#f9fafb;border-top:1px solid #f3e8ff;
                        padding:24px 40px;text-align:center;">
              <p style="margin:0;font-size:13px;color:#9ca3af;">
                © {{.Year}} <strong style="color:#7c3aed;">{{.AppName}}</strong>
                &nbsp;·&nbsp; Este mail fue enviado automáticamente, no respondas a este correo.
              </p>
            </td>
          </tr>

        </table>
      </td>
    </tr>
  </table>

</body>
</html>`

// templateData es el struct interno que se pasa al template HTML.
type templateData struct {
	AppointmentEmailData
	Title          string
	Greeting       template.HTML
	StatusLabel    string
	ClosingMessage template.HTML
	Year           int
}

// BuildClientAppointmentEmail genera el HTML del mail al cliente al crear un turno.
func BuildClientAppointmentEmail(data AppointmentEmailData) (string, error) {
	td := templateData{
		AppointmentEmailData: data,
		Title:                "¡Tu turno fue registrado!",
		Greeting:             template.HTML(fmt.Sprintf("Hola <strong>%s</strong>, tu turno ha sido registrado con éxito. Te esperamos 💜", template.HTMLEscapeString(data.ClientName))),
		StatusLabel:          translateStatus(data.Status),
		ClosingMessage:       "Si necesitás cancelar o modificar tu turno, contactanos con anticipación. ¡Gracias por elegirnos!",
		Year:                 time.Now().Year(),
	}
	return renderTemplate(td)
}

// BuildAdminAppointmentEmail genera el HTML del mail al admin al crearse un turno.
func BuildAdminAppointmentEmail(data AppointmentEmailData) (string, error) {
	td := templateData{
		AppointmentEmailData: data,
		Title:                "Nuevo turno registrado",
		Greeting:             template.HTML(fmt.Sprintf("Se registró un nuevo turno para <strong>%s</strong>.", template.HTMLEscapeString(data.ClientName))),
		StatusLabel:          translateStatus(data.Status),
		ClosingMessage:       "Revisá el panel de administración para gestionar el turno.",
		Year:                 time.Now().Year(),
	}
	return renderTemplate(td)
}

func renderTemplate(data templateData) (string, error) {
	tmpl, err := template.New("email").Parse(appointmentHTMLTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func translateStatus(status string) string {
	switch status {
	case "PENDING":
		return "⏳ Pendiente de confirmación"
	case "CONFIRMED":
		return "✅ Confirmado"
	case "DONE":
		return "🎉 Completado"
	case "CANCELLED":
		return "❌ Cancelado"
	case "ABSENT":
		return "🚫 Ausente"
	default:
		return status
	}
}
