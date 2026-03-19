-- Migración 002 — Índices compuestos para optimizar queries frecuentes
-- Aplicar manualmente: psql -d <db> -f migrations/002_add_composite_indexes.sql

-- Cubre: ExistsByClientAndDate → WHERE client_id = ? AND date = ? AND status NOT IN (...)
--        GetAllClientStats     → GROUP BY client_id con filtros de status
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_appointments_client_date_status
    ON appointments(client_id, date, status)
    WHERE deleted_at IS NULL;

-- Cubre: GetMetrics → WHERE date >= ? AND status = 'DONE'/'CONFIRMED'
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_appointments_date_status
    ON appointments(date, status)
    WHERE deleted_at IS NULL;
