-- =============================================================================
--  Migración inicial — Shei Nails Turnos
--  Nota: AutoMigrate de GORM crea estas tablas automáticamente al iniciarse
--  la aplicación. Este archivo es la referencia formal del esquema y sirve
--  como base para migraciones futuras con herramientas como golang-migrate.
-- =============================================================================

-- Clientes (clientas que reservan turnos)
CREATE TABLE IF NOT EXISTS clients (
    id            SERIAL PRIMARY KEY,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ,
    first_name    VARCHAR(100) NOT NULL,
    last_name     VARCHAR(100) NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    phone         VARCHAR(50),
    password_hash VARCHAR(255) NOT NULL,
    active        BOOLEAN NOT NULL DEFAULT TRUE
);
CREATE INDEX IF NOT EXISTS idx_clients_email ON clients(email);

-- Admins / usuarios internos
CREATE TABLE IF NOT EXISTS admins (
    id            SERIAL PRIMARY KEY,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ,
    name          VARCHAR(150) NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(20) NOT NULL DEFAULT 'admin',
    active        BOOLEAN NOT NULL DEFAULT TRUE
);

-- Servicios ofrecidos
CREATE TABLE IF NOT EXISTS services (
    id                SERIAL PRIMARY KEY,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ,
    name              VARCHAR(150) NOT NULL,
    description       TEXT,
    duration_minutes  INT NOT NULL,
    base_price        NUMERIC(10,2) NOT NULL,
    requires_deposit  BOOLEAN NOT NULL DEFAULT FALSE,
    suggested_deposit NUMERIC(10,2) NOT NULL DEFAULT 0,
    color             VARCHAR(20) NOT NULL DEFAULT '#ffffff',
    active            BOOLEAN NOT NULL DEFAULT TRUE
);

-- Configuración horaria semanal
-- day_of_week: 0=Domingo … 6=Sábado (convención time.Weekday de Go)
CREATE TABLE IF NOT EXISTS weekly_schedules (
    id               SERIAL PRIMARY KEY,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ,
    day_of_week      INT NOT NULL UNIQUE CHECK (day_of_week BETWEEN 0 AND 6),
    enabled          BOOLEAN NOT NULL DEFAULT TRUE,
    opening_time     VARCHAR(5) NOT NULL,   -- "09:00"
    closing_time     VARCHAR(5) NOT NULL,   -- "19:00"
    slot_duration_min INT NOT NULL DEFAULT 30
);

-- Bloqueos de agenda (días/rangos horarios no disponibles)
CREATE TABLE IF NOT EXISTS blocked_slots (
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    date       VARCHAR(10) NOT NULL,   -- "2026-03-16"
    start_time VARCHAR(5) NOT NULL,    -- "14:00"
    end_time   VARCHAR(5) NOT NULL,    -- "15:00"
    reason     TEXT,
    permanent  BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_blocked_slots_date ON blocked_slots(date);
CREATE INDEX IF NOT EXISTS idx_blocked_slots_date_times ON blocked_slots(date, start_time, end_time);

-- Turnos
CREATE TABLE IF NOT EXISTS appointments (
    id              SERIAL PRIMARY KEY,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,
    client_id       INT NOT NULL REFERENCES clients(id),
    service_id      INT NOT NULL REFERENCES services(id),
    professional_id INT,                -- nullable: soporte futuro para profesionales
    date            VARCHAR(10) NOT NULL,   -- "2026-03-16"
    start_time      VARCHAR(5) NOT NULL,    -- "14:00"
    end_time        VARCHAR(5) NOT NULL,    -- "15:00"
    base_price      NUMERIC(10,2) NOT NULL,
    extras_amount   NUMERIC(10,2) NOT NULL DEFAULT 0,
    deposit_amount  NUMERIC(10,2) NOT NULL DEFAULT 0,
    final_price     NUMERIC(10,2) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'PENDING'
                    CHECK (status IN ('PENDING','CONFIRMED','DONE','CANCELLED','ABSENT')),
    notes           TEXT
);
-- Índices recomendados para las consultas más frecuentes
CREATE INDEX IF NOT EXISTS idx_appointments_date              ON appointments(date);
CREATE INDEX IF NOT EXISTS idx_appointments_date_times        ON appointments(date, start_time, end_time);
CREATE INDEX IF NOT EXISTS idx_appointments_client_id         ON appointments(client_id);
CREATE INDEX IF NOT EXISTS idx_appointments_service_id        ON appointments(service_id);
CREATE INDEX IF NOT EXISTS idx_appointments_status            ON appointments(status);
-- Índices compuestos para las queries más críticas:
-- ExistsByClientAndDate: WHERE client_id = ? AND date = ? AND status NOT IN (...)
CREATE INDEX IF NOT EXISTS idx_appointments_client_date_status ON appointments(client_id, date, status);
-- Métricas: WHERE date >= ? AND status = 'DONE' / CONFIRMED
CREATE INDEX IF NOT EXISTS idx_appointments_date_status        ON appointments(date, status);
