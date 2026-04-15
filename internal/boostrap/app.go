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

func BuildApp() (*gin.Engine, string) {
	cfg := config.Load()

	db := postgres.NewConnection(cfg)

	clientRepo := repositories.NewClientRepository(db)
	adminRepo := repositories.NewAdminRepository(db)
	serviceRepo := repositories.NewServiceRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	appointmentRepo := repositories.NewAppointmentRepository(db)

	hasher := security.NewBcryptHasher()
	tokenGen := security.NewJWTGenerator(cfg.JWTSecret)

	var mailer inframail.Mailer
	if cfg.SMTPUser != "" && cfg.SMTPPass != "" {
		mailer = inframail.NewSMTPMailer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.AdminEmail)
	} else {
		mailer = inframail.NewMockMailer()
	}

	registerUC := authapp.NewRegisterUseCase(clientRepo, hasher, tokenGen)
	loginUC := authapp.NewLoginUseCase(clientRepo, hasher, tokenGen)
	adminLoginUC := authapp.NewAdminLoginUseCase(adminRepo, hasher, tokenGen)

	createServiceUC := serviceapp.NewCreateServiceUseCase(serviceRepo)
	listServiceUC := serviceapp.NewListServiceUseCase(serviceRepo)
	getServiceUC := serviceapp.NewGetServiceByIDUseCase(serviceRepo)
	updateServiceUC := serviceapp.NewUpdateServiceUseCase(serviceRepo)
	toggleServiceUC := serviceapp.NewToggleServiceUseCase(serviceRepo)

	listClientsUC := clientapp.NewListClientsUseCase(clientRepo, appointmentRepo)

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

	getWeeklyUC := scheduleapp.NewGetWeeklyScheduleUseCase(scheduleRepo)
	updateWeeklyUC := scheduleapp.NewUpdateWeeklyScheduleUseCase(scheduleRepo)
	blockSlotUC := scheduleapp.NewBlockTimeSlotUseCase(scheduleRepo)
	availabilityUC := scheduleapp.NewGetAvailabilityUseCase(scheduleRepo, appointmentRepo)

	getMetricsUC := dashboardapp.NewGetMetricsUseCase(appointmentRepo)

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
