package http

import (
	"apiGoShei/internal/interfaces/http/handlers"
	"apiGoShei/internal/interfaces/http/middleware"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configura todas las rutas del servidor HTTP.
// Los handlers y el secreto JWT se reciben desde bootstrap, manteniendo el
// router desacoplado de la construcción de dependencias.
func SetupRouter(
	authHandler *handlers.AuthHandler,
	appointmentHandler *handlers.AppointmentHandler,
	serviceHandler *handlers.ServiceHandler,
	scheduleHandler *handlers.ScheduleHandler,
	dashboardHandler *handlers.DashboardHandler,
	clientHandler *handlers.ClientHandler,
	jwtSecret string,
) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"https://nailsdeni.com",
			"https://www.nailsdeni.com",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	base := r.Group("/service-nails")

	// ─── Health check ────────────────────────────────────────────────────────
	base.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// ─── Rutas públicas ───────────────────────────────────────────────────────
	auth := base.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/admin/login", authHandler.AdminLogin)
	}

	// Servicios: lectura pública (sin token)
	pub := base.Group("/services")
	{
		pub.GET("", serviceHandler.List)
		pub.GET("/:id", serviceHandler.GetByID)
	}

	// Horarios: lectura pública (sin token) — el cliente los necesita para reservar
	pubSch := base.Group("/schedule")
	{
		pubSch.GET("/weekly", scheduleHandler.GetWeekly)
		pubSch.GET("/availability", scheduleHandler.GetAvailability)
	}

	// ─── Rutas compartidas (cliente + admin) ───────────────────────────────────
	// GET /appointments usa ListByRole: cliente ve solo los suyos, admin ve todos.
	shared := base.Group("")
	shared.Use(middleware.AuthMiddleware(jwtSecret))
	shared.Use(middleware.RequireRole("client", "admin", "superadmin"))
	{
		shared.GET("/appointments", appointmentHandler.ListByRole)
		shared.GET("/appointments/next", appointmentHandler.NextAppointment)
		shared.GET("/appointments/:id", appointmentHandler.GetByID)
		shared.PATCH("/appointments/:id/cancel", appointmentHandler.Cancel)
		shared.POST("/appointments", appointmentHandler.Create)
	}

	// ─── Rutas de admin (role: admin, superadmin) ─────────────────────────────
	admin := base.Group("")
	admin.Use(middleware.AuthMiddleware(jwtSecret))
	admin.Use(middleware.RequireRole("admin", "superadmin"))
	{
		// Services (solo escritura/admin)
		svc := admin.Group("/services")
		{
			svc.POST("", serviceHandler.Create)
			svc.PATCH("/:id", serviceHandler.Update)
			svc.PATCH("/:id/toggle", serviceHandler.Toggle)
		}

		// Appointments (gestión exclusiva de admin)
		appt := admin.Group("/appointments")
		{
			appt.GET("/calendar", appointmentHandler.Calendar)
			appt.PATCH("/:id/confirm", appointmentHandler.Confirm)
			appt.PATCH("/:id/complete", appointmentHandler.Complete)
			appt.PATCH("/:id/final-price", appointmentHandler.UpdateFinalPrice)
			appt.PATCH("/:id/deposit", appointmentHandler.UpdateDeposit)
		}

		// Schedule (solo escritura/admin)
		sch := admin.Group("/schedule")
		{
			sch.PUT("/weekly", scheduleHandler.UpdateWeekly)
			sch.POST("/blocked-slots", scheduleHandler.BlockSlot)
		}

		// Clients
		admin.GET("/clients", clientHandler.List)

		// Dashboard
		admin.GET("/dashboard/metrics", dashboardHandler.GetMetrics)
	}

	return r
}
