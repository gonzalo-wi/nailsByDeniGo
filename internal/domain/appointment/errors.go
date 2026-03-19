package appointment

import "errors"

var (
	ErrNotFound         = errors.New("turno no encontrado")
	ErrSlotUnavailable  = errors.New("el horario no está disponible")
	ErrOverlap          = errors.New("ya existe un turno en ese horario")
	ErrInvalidStatus    = errors.New("transición de estado inválida")
	ErrPastDate         = errors.New("no se pueden agendar turnos en fechas pasadas")
	ErrClientHasThatDay = errors.New("ya tenés un turno agendado para ese día")
)
