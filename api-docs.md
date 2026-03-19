# API Documentation — Service Nails

> Instrucciones para el frontend. Pegar este archivo como contexto en GitHub Copilot del proyecto frontend.

## Base URL

```
/service-nails
```

---

## Autenticación

Todas las rutas protegidas requieren el siguiente header:

```
Authorization: Bearer <token>
```

El token JWT contiene:
| Campo | Tipo | Descripción |
|-------|------|-------------|
| `user_id` | `uint` | ID del usuario (cliente o admin) |
| `role` | `string` | `client` \| `admin` \| `superadmin` |

- **Expiración:** 72 horas
- **Algoritmo:** HS256

### Roles y acceso
| Rol | Acceso |
|-----|--------|
| *(sin token)* | Solo rutas de auth |
| `client` | Sus propios turnos |
| `admin` / `superadmin` | Gestión completa (servicios, turnos, horarios, dashboard) |

---

## Tipos de respuesta comunes

### `AppointmentResponse`
```json
{
  "id": 1,
  "client_id": 1,
  "service_id": 1,
  "professional_id": 2,
  "date": "2026-03-16",
  "start_time": "14:00",
  "end_time": "15:00",
  "base_price": 1500.00,
  "extras_amount": 200.00,
  "deposit_amount": 500.00,
  "final_price": 1700.00,
  "status": "PENDING",
  "notes": "Quiero uñas largas"
}
```

> `professional_id` y `notes` pueden ser `null` / omitidos.

**Valores posibles de `status`:**
- `PENDING` — Pendiente de confirmación
- `CONFIRMED` — Confirmado
- `DONE` — Completado
- `CANCELLED` — Cancelado
- `ABSENT` — Cliente ausente

---

### `ServiceResponse`
```json
{
  "id": 1,
  "name": "Manicura",
  "description": "Descripción del servicio",
  "duration_minutes": 60,
  "base_price": 1500.00,
  "requires_deposit": true,
  "suggested_deposit": 500.00,
  "color": "#FF0000",
  "active": true
}
```

---

### `WeeklyScheduleResponse`
```json
{
  "id": 1,
  "day_of_week": 1,
  "day_name": "Lunes",
  "enabled": true,
  "opening_time": "09:00",
  "closing_time": "19:00",
  "slot_duration_min": 30
}
```

> `day_of_week`: 0=Domingo, 1=Lunes, 2=Martes, 3=Miércoles, 4=Jueves, 5=Viernes, 6=Sábado

---

### `BlockedSlotResponse`
```json
{
  "id": 1,
  "date": "2026-03-16",
  "start_time": "14:00",
  "end_time": "15:00",
  "reason": "Feriado",
  "permanent": false
}
```

---

## Respuestas de error

Todos los errores devuelven un JSON con la clave `error`:

```json
{ "error": "mensaje descriptivo" }
```

| Código HTTP | Significado |
|-------------|-------------|
| 400 | Request inválido / validación fallida |
| 401 | Token ausente, inválido o expirado |
| 403 | Rol sin permiso para esta ruta |
| 404 | Recurso no encontrado |
| 409 | Conflicto (ej: email ya registrado) |
| 500 | Error interno del servidor |

---

## Endpoints de Auth (Público — sin token)

### `POST /service-nails/auth/register`
Registrar un nuevo cliente.

**Request body:**
```json
{
  "first_name": "María",
  "last_name": "López",
  "email": "maria@example.com",
  "phone": "1122334455",
  "password": "mipassword123"
}
```

| Campo | Tipo | Requerido | Reglas |
|-------|------|-----------|--------|
| `first_name` | string | ✅ | — |
| `last_name` | string | ✅ | — |
| `email` | string | ✅ | formato email |
| `phone` | string | ❌ | — |
| `password` | string | ✅ | mínimo 6 caracteres |

**Response `201`:**
```json
{
  "token": "eyJhbGciOiJIUzI1...",
  "client_id": 5
}
```

**Errores:**
- `400` — Campos inválidos
- `409` — El email ya está registrado

---

### `POST /service-nails/auth/login`
Login de cliente.

**Request body:**
```json
{
  "email": "maria@example.com",
  "password": "mipassword123"
}
```

**Response `200`:**
```json
{
  "token": "eyJhbGciOiJIUzI1...",
  "client_id": 5
}
```

**Errores:**
- `400` — Campos inválidos
- `401` — Credenciales incorrectas

---

### `POST /service-nails/auth/admin/login`
Login de administrador.

**Request body:**
```json
{
  "email": "admin@example.com",
  "password": "adminpass123"
}
```

**Response `200`:**
```json
{
  "token": "eyJhbGciOiJIUzI1...",
  "admin_id": 1
}
```

**Errores:**
- `400` — Campos inválidos
- `401` — Credenciales incorrectas

---

## Endpoints de Turnos — Cliente (`role: client`)

Requieren `Authorization: Bearer <token>` con rol `client`.

---

### `POST /service-nails/appointments`
Crear un nuevo turno.

**Request body:**
```json
{
  "client_id": 5,
  "service_id": 2,
  "professional_id": 1,
  "date": "2026-03-20",
  "start_time": "14:00",
  "notes": "Quiero diseño floral"
}
```

| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `client_id` | uint | ✅ | ID del cliente |
| `service_id` | uint | ✅ | ID del servicio |
| `professional_id` | uint | ❌ | ID del profesional |
| `date` | string | ✅ | Formato `YYYY-MM-DD` |
| `start_time` | string | ✅ | Formato `HH:MM` |
| `notes` | string | ❌ | Notas adicionales |

**Response `201`:** [`AppointmentResponse`](#appointmentresponse)

**Errores:**
- `400` — Campos inválidos, horario no disponible, turno superpuesto

---

### `GET /service-nails/appointments`
Listar mis turnos (solo los del cliente autenticado).

**Query params:**
| Param | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `status` | string | ❌ | Filtrar por estado: `PENDING`, `CONFIRMED`, `DONE`, `CANCELLED`, `ABSENT` |

**Response `200`:** Array de [`AppointmentResponse`](#appointmentresponse)

---

### `GET /service-nails/appointments/:id`
Obtener un turno por ID.

**Response `200`:** [`AppointmentResponse`](#appointmentresponse)

**Errores:**
- `404` — Turno no encontrado

---

### `PATCH /service-nails/appointments/:id/cancel`
Cancelar un turno propio.

**Body:** *(ninguno)*

**Response `200`:** [`AppointmentResponse`](#appointmentresponse) con `status: "CANCELLED"`

**Errores:**
- `400` — Transición de estado inválida
- `404` — Turno no encontrado

---

## Endpoints de Servicios — Admin (`role: admin | superadmin`)

Requieren `Authorization: Bearer <token>` con rol `admin` o `superadmin`.

---

### `POST /service-nails/services`
Crear un nuevo servicio.

**Request body:**
```json
{
  "name": "Manicura Semipermanente",
  "description": "Esmalte que dura hasta 3 semanas",
  "duration_minutes": 60,
  "base_price": 2500.00,
  "requires_deposit": true,
  "suggested_deposit": 800.00,
  "color": "#E91E8C"
}
```

| Campo | Tipo | Requerido | Reglas |
|-------|------|-----------|--------|
| `name` | string | ✅ | — |
| `description` | string | ❌ | — |
| `duration_minutes` | int | ✅ | mínimo 1 |
| `base_price` | float64 | ✅ | mínimo 0 |
| `requires_deposit` | bool | ❌ | default `false` |
| `suggested_deposit` | float64 | ❌ | — |
| `color` | string | ❌ | Ej: `#FF0000` |

**Response `201`:** [`ServiceResponse`](#serviceresponse)

---

### `GET /service-nails/services`
Listar servicios.

**Query params:**
| Param | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `active` | string | ❌ | Por defecto solo activos. Pasar `false` para ver todos |

**Response `200`:** Array de [`ServiceResponse`](#serviceresponse)

---

### `GET /service-nails/services/:id`
Obtener servicio por ID.

**Response `200`:** [`ServiceResponse`](#serviceresponse)

**Errores:**
- `404` — Servicio no encontrado

---

### `PATCH /service-nails/services/:id`
Actualizar un servicio.

**Request body:** Igual que `POST /services`

**Response `200`:** [`ServiceResponse`](#serviceresponse)

**Errores:**
- `404` — Servicio no encontrado

---

### `PATCH /service-nails/services/:id/toggle`
Activar / desactivar un servicio (alterna el campo `active`).

**Body:** *(ninguno)*

**Response `200`:** [`ServiceResponse`](#serviceresponse) con el nuevo valor de `active`

**Errores:**
- `404` — Servicio no encontrado

---

## Endpoints de Turnos — Admin (`role: admin | superadmin`)

---

### `GET /service-nails/appointments`
Listar todos los turnos (vista admin).

**Query params:**
| Param | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `status` | string | ❌ | `PENDING` \| `CONFIRMED` \| `DONE` \| `CANCELLED` \| `ABSENT` |
| `client_id` | uint | ❌ | Filtrar por cliente |
| `date_from` | string | ❌ | Desde fecha `YYYY-MM-DD` |
| `date_to` | string | ❌ | Hasta fecha `YYYY-MM-DD` |

**Response `200`:** Array de [`AppointmentResponse`](#appointmentresponse)

---

### `GET /service-nails/appointments/calendar`
Obtener turnos en rango de fechas (vista calendario).

**Query params:**
| Param | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `from` | string | ✅ | Fecha inicio `YYYY-MM-DD` |
| `to` | string | ✅ | Fecha fin `YYYY-MM-DD` |

**Response `200`:** Array de [`AppointmentResponse`](#appointmentresponse)

---

### `GET /service-nails/appointments/:id`
Obtener turno por ID.

**Response `200`:** [`AppointmentResponse`](#appointmentresponse)

---

### `PATCH /service-nails/appointments/:id/confirm`
Confirmar un turno (`PENDING` → `CONFIRMED`).

**Body:** *(ninguno)*

**Response `200`:** [`AppointmentResponse`](#appointmentresponse)

**Errores:**
- `400` — Estado inválido para esta transición
- `404` — Turno no encontrado

---

### `PATCH /service-nails/appointments/:id/cancel`
Cancelar un turno.

**Body:** *(ninguno)*

**Response `200`:** [`AppointmentResponse`](#appointmentresponse)

---

### `PATCH /service-nails/appointments/:id/complete`
Marcar turno como completado (`CONFIRMED` → `DONE`).

**Body:** *(ninguno)*

**Response `200`:** [`AppointmentResponse`](#appointmentresponse)

---

### `PATCH /service-nails/appointments/:id/final-price`
Registrar precio final (puede incluir extras).

**Request body:**
```json
{
  "final_price": 2800.00,
  "extras_amount": 300.00
}
```

| Campo | Tipo | Requerido | Reglas |
|-------|------|-----------|--------|
| `final_price` | float64 | ✅ | mínimo 0 |
| `extras_amount` | float64 | ❌ | — |

**Response `200`:** [`AppointmentResponse`](#appointmentresponse)

---

### `PATCH /service-nails/appointments/:id/deposit`
Registrar seña / depósito de un turno.

**Request body:**
```json
{
  "deposit_amount": 800.00
}
```

| Campo | Tipo | Requerido | Reglas |
|-------|------|-----------|--------|
| `deposit_amount` | float64 | ✅ | mínimo 0 |

**Response `200`:** [`AppointmentResponse`](#appointmentresponse)

---

## Endpoints de Horarios — Admin (`role: admin | superadmin`)

---

### `GET /service-nails/schedule/weekly`
Obtener la configuración de horarios semanales.

**Response `200`:** Array de [`WeeklyScheduleResponse`](#weeklyscheduleresponse) (7 elementos, uno por día)

---

### `PUT /service-nails/schedule/weekly`
Actualizar la configuración de horarios semanales.

**Request body:**
```json
{
  "schedule": [
    {
      "day_of_week": 1,
      "enabled": true,
      "opening_time": "09:00",
      "closing_time": "19:00",
      "slot_duration_min": 30
    },
    {
      "day_of_week": 0,
      "enabled": false,
      "opening_time": "09:00",
      "closing_time": "13:00",
      "slot_duration_min": 30
    }
  ]
}
```

| Campo | Tipo | Requerido | Reglas |
|-------|------|-----------|--------|
| `day_of_week` | int | ✅ | 0–6 |
| `enabled` | bool | ✅ | — |
| `opening_time` | string | ✅ | Formato `HH:MM` |
| `closing_time` | string | ✅ | Formato `HH:MM` |
| `slot_duration_min` | int | ✅ | mínimo 10 |

**Response `200`:** Array de [`WeeklyScheduleResponse`](#weeklyscheduleresponse)

---

### `POST /service-nails/schedule/blocked-slots`
Bloquear un horario específico.

**Request body:**
```json
{
  "date": "2026-03-25",
  "start_time": "14:00",
  "end_time": "16:00",
  "reason": "Feriado nacional",
  "permanent": false
}
```

| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `date` | string | ✅ | Formato `YYYY-MM-DD` |
| `start_time` | string | ✅ | Formato `HH:MM` |
| `end_time` | string | ✅ | Formato `HH:MM` |
| `reason` | string | ❌ | Motivo del bloqueo |
| `permanent` | bool | ❌ | Si aplica todos los años en esa fecha |

**Response `201`:** [`BlockedSlotResponse`](#blockedslotresponse)

---

### `GET /service-nails/schedule/availability`
Consultar turnos disponibles para una fecha.

**Query params:**
| Param | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `date` | string | ✅ | Formato `YYYY-MM-DD` |
| `duration` | int | ❌ | Duración en minutos. Si se omite, usa `slot_duration_min` del día |

**Response `200`:**
```json
{
  "date": "2026-03-20",
  "opening_time": "09:00",
  "closing_time": "19:00",
  "available_slots": [
    { "start_time": "09:00", "end_time": "09:30" },
    { "start_time": "10:00", "end_time": "10:30" },
    { "start_time": "14:30", "end_time": "15:00" }
  ]
}
```

> Si el día no está habilitado o está completamente bloqueado, `available_slots` será `[]`.

---

## Endpoints de Dashboard — Admin (`role: admin | superadmin`)

---

### `GET /service-nails/dashboard/metrics`
Obtener métricas generales del negocio.

**Response `200`:**
```json
{
  "today_appointments": 5,
  "week_appointments": 23,
  "today_revenue": 7500.00,
  "month_revenue": 45000.00,
  "appointments_by_status": {
    "PENDING": 3,
    "CONFIRMED": 10,
    "DONE": 100,
    "CANCELLED": 5,
    "ABSENT": 2
  },
  "total_deposits": 15000.00
}
```

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `today_appointments` | int64 | Turnos de hoy |
| `week_appointments` | int64 | Turnos de esta semana |
| `today_revenue` | float64 | Ingresos de hoy |
| `month_revenue` | float64 | Ingresos del mes |
| `appointments_by_status` | map[string]int64 | Conteo por estado |
| `total_deposits` | float64 | Total de señas recibidas |

---

## Resumen de endpoints

| Método | Ruta | Auth | Rol |
|--------|------|------|-----|
| POST | `/service-nails/auth/register` | ❌ | — |
| POST | `/service-nails/auth/login` | ❌ | — |
| POST | `/service-nails/auth/admin/login` | ❌ | — |
| POST | `/service-nails/appointments` | ✅ | client |
| GET | `/service-nails/appointments` | ✅ | client |
| GET | `/service-nails/appointments/:id` | ✅ | client |
| PATCH | `/service-nails/appointments/:id/cancel` | ✅ | client |
| POST | `/service-nails/services` | ✅ | admin/superadmin |
| GET | `/service-nails/services` | ✅ | admin/superadmin |
| GET | `/service-nails/services/:id` | ✅ | admin/superadmin |
| PATCH | `/service-nails/services/:id` | ✅ | admin/superadmin |
| PATCH | `/service-nails/services/:id/toggle` | ✅ | admin/superadmin |
| GET | `/service-nails/appointments` | ✅ | admin/superadmin |
| GET | `/service-nails/appointments/calendar` | ✅ | admin/superadmin |
| GET | `/service-nails/appointments/:id` | ✅ | admin/superadmin |
| PATCH | `/service-nails/appointments/:id/confirm` | ✅ | admin/superadmin |
| PATCH | `/service-nails/appointments/:id/cancel` | ✅ | admin/superadmin |
| PATCH | `/service-nails/appointments/:id/complete` | ✅ | admin/superadmin |
| PATCH | `/service-nails/appointments/:id/final-price` | ✅ | admin/superadmin |
| PATCH | `/service-nails/appointments/:id/deposit` | ✅ | admin/superadmin |
| GET | `/service-nails/schedule/weekly` | ✅ | admin/superadmin |
| PUT | `/service-nails/schedule/weekly` | ✅ | admin/superadmin |
| POST | `/service-nails/schedule/blocked-slots` | ✅ | admin/superadmin |
| GET | `/service-nails/schedule/availability` | ✅ | admin/superadmin |
| GET | `/service-nails/dashboard/metrics` | ✅ | admin/superadmin |
