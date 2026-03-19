// Package bootstrap ensambla todas las dependencias del sistema siguiendo
// el patrón de inyección de dependencias manual (Composition Root).
// Ningún otro paquete debería conocer la estructura global de dependencias.
package bootstrap

import (
	appointmentapp "apiGoShei/internal/application/appointment"
	authapp "apiGoShei/internal/application/auth"
	clientapp "apiGoShei/internal/application/client"
	dashboardapp "apiGoShei/internal/application/dashboard"
	scheduleapp "apiGoShei/internal/application/schedule"
	serviceapp "apiGoShei/internal/application/service"
	"apiGoShei/internal/infrastructure/config"
	inframail "apiGoShei/internal/infrastructure/mail"
	"apiGoShei/internal/infrastructure/persistence/postgres"
	"apiGoShei/internal/infrastructure/persistence/postgres/repositories"
	"apiGoShei/internal/infrastructure/security"
	httprouter "apiGoShei/internal/interfaces/http"
	"apiGoShei/internal/interfaces/http/handlers"

	"github.com/gin-gonic/gin"
)

// BuildApp carga la configuración, conecta a la base de datos, construye
// todos los casos de uso y devuelve el router Gin listo para escuchar junto
// con el puerto configurado.
func BuildApp() (*gin.Engine, string) {
	cfg := config.Load()

	// ─── Infraestructura ─────────────────────────────────────────────────────
	db := postgres.NewConnection(cfg)

	// Repositorios
	clientRepo := repositories.NewClientRepository(db)
	adminRepo := repositories.NewAdminRepository(db)
	serviceRepo := repositories.NewServiceRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	appointmentRepo := repositories.NewAppointmentRepository(db)

	// Seguridad
	hasher := security.NewBcryptHasher()
	tokenGen := security.NewJWTGenerator(cfg.JWTSecret)

	// Mailer: SMTPMailer si hay credenciales configuradas, MockMailer si no.
	var mailer inframail.Mailer
	if cfg.SMTPUser != "" && cfg.SMTPPass != "" {
		mailer = inframail.NewSMTPMailer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.AdminEmail)
	} else {
		mailer = inframail.NewMockMailer()
	}

	// ─── Casos de uso — Auth ─────────────────────────────────────────────────
	registerUC := authapp.NewRegisterUseCase(clientRepo, hasher, tokenGen)
	loginUC := authapp.NewLoginUseCase(clientRepo, hasher, tokenGen)
	adminLoginUC := authapp.NewAdminLoginUseCase(adminRepo, hasher, tokenGen)

	// ─── Casos de uso — Servicios ─────────────────────────────────────────────
	createServiceUC := serviceapp.NewCreateServiceUseCase(serviceRepo)
	listServiceUC := serviceapp.NewListServiceUseCase(serviceRepo)
	getServiceUC := serviceapp.NewGetServiceByIDUseCase(serviceRepo)
	updateServiceUC := serviceapp.NewUpdateServiceUseCase(serviceRepo)
	toggleServiceUC := serviceapp.NewToggleServiceUseCase(serviceRepo)

	// ─── Casos de uso — Clientes ─────────────────────────────────────────────
	listClientsUC := clientapp.NewListClientsUseCase(clientRepo, appointmentRepo)

	// ─── Casos de uso — Turnos ────────────────────────────────────────────────
	createApptUC := appointmentapp.NewCreateAppointmentUseCase(appointmentRepo, clientRepo, serviceRepo, scheduleRepo, mailer)
	cancelApptUC := appointmentapp.NewCancelAppointmentUseCase(appointmentRepo)
	confirmApptUC := appointmentapp.NewConfirmAppointmentUseCase(appointmentRepo)
	completeApptUC := appointmentapp.NewCompleteAppointmentUseCase(appointmentRepo)
	listApptUC := appointmentapp.NewListAppointmentUseCase(appointmentRepo)
	getApptUC := appointmentapp.NewGetAppointmentByIDUseCase(appointmentRepo)
	calendarApptUC := appointmentapp.NewListCalendarAppointmentUseCase(appointmentRepo)
	finalPriceUC := appointmentapp.NewUpdateFinalPriceUseCase(appointmentRepo)
	updateDepositUC := appointmentapp.NewUpdateDepositUseCase(appointmentRepo)
	nextApptUC := appointmentapp.NewNextAppointmentUseCase(appointmentRepo)

	// ─── Casos de uso — Horarios ──────────────────────────────────────────────
	getWeeklyUC := scheduleapp.NewGetWeeklyScheduleUseCase(scheduleRepo)
	updateWeeklyUC := scheduleapp.NewUpdateWeeklyScheduleUseCase(scheduleRepo)
	blockSlotUC := scheduleapp.NewBlockTimeSlotUseCase(scheduleRepo)
	availabilityUC := scheduleapp.NewGetAvailabilityUseCase(scheduleRepo, appointmentRepo)

	// ─── Casos de uso — Dashboard ─────────────────────────────────────────────
	getMetricsUC := dashboardapp.NewGetMetricsUseCase(appointmentRepo)

	// ─── Handlers HTTP ────────────────────────────────────────────────────────
	clientHandler := handlers.NewClientHandler(listClientsUC)
	authHandler := handlers.NewAuthHandler(registerUC, loginUC, adminLoginUC)
	appointmentHandler := handlers.NewAppointmentHandler(
		createApptUC, cancelApptUC, confirmApptUC, completeApptUC,
		listApptUC, getApptUC, calendarApptUC, finalPriceUC, updateDepositUC, nextApptUC,
	)
	serviceHandler := handlers.NewServiceHandler(createServiceUC, listServiceUC, getServiceUC, updateServiceUC, toggleServiceUC)
	scheduleHandler := handlers.NewScheduleHandler(getWeeklyUC, updateWeeklyUC, blockSlotUC, availabilityUC)
	dashboardHandler := handlers.NewDashboardHandler(getMetricsUC)

	return httprouter.SetupRouter(authHandler, appointmentHandler, serviceHandler, scheduleHandler, dashboardHandler, clientHandler, cfg.JWTSecret), cfg.AppPort
}
