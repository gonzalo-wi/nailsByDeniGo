package postgres

import (
	"fmt"

	"apiGoShei/internal/infrastructure/config"
	"apiGoShei/internal/infrastructure/logger"
	"apiGoShei/internal/infrastructure/persistence/postgres/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// NewConnection abre la conexión a PostgreSQL y corre las auto-migraciones.
// Panea si no puede conectar, ya que la app no puede funcionar sin DB.
func NewConnection(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=America/Argentina/Buenos_Aires",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	logLevel := gormlogger.Warn
	if cfg.AppEnv == "development" {
		logLevel = gormlogger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		panic("no se pudo conectar a la base de datos: " + err.Error())
	}

	runMigrations(db)
	logger.Info.Println("Conexión a PostgreSQL establecida")
	return db
}

func runMigrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.ClientModel{},
		&models.AdminModel{},
		&models.ServiceModel{},
		&models.WeeklyScheduleModel{},
		&models.BlockedSlotModel{},
		&models.AppointmentModel{},
	)
	if err != nil {
		panic("error al correr migraciones: " + err.Error())
	}
}
