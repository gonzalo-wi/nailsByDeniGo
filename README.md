# Shei Nails — API de Gestión de Turnos

API REST desarrollada en Go para la gestión completa de turnos de un local de uñas. Permite a los clientes registrarse, ver disponibilidad y reservar turnos, mientras el administrador gestiona la agenda, los servicios y el negocio desde un panel dedicado.

## Tecnologías

| | |
|---|---|
| **Lenguaje** | Go 1.23 |
| **Framework HTTP** | Gin |
| **ORM** | GORM |
| **Base de datos** | PostgreSQL |
| **Autenticación** | JWT (golang-jwt/jwt v5) |
| **Passwords** | bcrypt (golang.org/x/crypto) |
| **Email** | SMTP (Hostinger / configurable) |
| **Config** | godotenv |

## Arquitectura

El proyecto sigue **arquitectura hexagonal** (puertos y adaptadores) con separación en capas:

```
cmd/
  api/          → punto de entrada del servidor
  seed/         → script de carga de datos iniciales
internal/
  application/  → casos de uso (lógica de negocio)
  domain/       → entidades y contratos de repositorio
  infrastructure/
    config/     → carga de variables de entorno
    logger/     → logger estructurado
    mail/       → mailer SMTP + mock
    persistence/postgres/  → modelos GORM + repositorios
    security/   → JWT + bcrypt
  interfaces/http/
    handlers/   → controladores HTTP
    middleware/  → autenticación y autorización por rol
    dto/         → request/response structs
    router.go   → definición de rutas
  boostrap/     → composition root (inyección de dependencias)
migrations/     → archivos SQL de referencia y migraciones incrementales
```

## Requisitos

- Go 1.23+
- PostgreSQL 14+

## Configuración

Crear un archivo `.env` en la raíz del proyecto:

```env
# Base de datos
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=tu_password
DB_NAME=turnos_db
DB_SSLMODE=disable

# Aplicación
APP_ENV=development
APP_PORT=8080

# JWT (usar un string aleatorio seguro en producción)
JWT_SECRET=tu_jwt_secret_seguro

# Email (opcional — si no se configura usa un mock que loguea en consola)
SMTP_HOST=smtp.hostinger.com
SMTP_PORT=587
SMTP_USER=noreply@tudominio.com
SMTP_PASS=tu_password_smtp
ADMIN_EMAIL=admin@tudominio.com
```

## Instalación y uso

```bash
# Clonar el repositorio
git clone https://github.com/tuusuario/apiGoShei.git
cd apiGoShei

# Instalar dependencias
go mod tidy

# Iniciar el servidor
go run cmd/api/main.go

# (Opcional) Cargar datos de prueba
go run cmd/seed/main.go
```

La API quedará disponible en `http://localhost:8080`.

## Migraciones

Las tablas se crean automáticamente via `AutoMigrate` de GORM al iniciar el servidor. Para aplicar migraciones incrementales (índices, cambios de esquema):

```bash
$env:PGPASSWORD='tu_password'
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -h localhost -U postgres -d turnos_db -f migrations/002_add_composite_indexes.sql
```

---

## Endpoints

Base URL: `/service-nails`

### Auth — público

| Método | Ruta | Descripción |
|--------|------|-------------|
| `POST` | `/auth/register` | Registro de cliente |
| `POST` | `/auth/login` | Login de cliente |
| `POST` | `/auth/admin/login` | Login de administrador |

**Login de cliente — respuesta:**
```json
{
  "token": "eyJ...",
  "client_id": 5,
  "first_name": "María",
  "last_name": "López",
  "email": "maria@example.com",
  "phone": "1122334455"
}
```

---

### Servicios — público

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/services` | Listar servicios activos |
| `GET` | `/services/:id` | Detalle de un servicio |

---

### Disponibilidad — público

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/schedule/weekly` | Horario semanal del local |
| `GET` | `/schedule/availability?date=2026-03-20` | Slots disponibles para una fecha |

---

### Turnos — cliente y admin `🔒`

> Requiere `Authorization: Bearer <token>`

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/appointments` | Lista de turnos (cliente: solo los suyos con historial; admin: últimos 30 días por defecto) |
| `GET` | `/appointments/next` | Próximo turno del cliente autenticado |
| `GET` | `/appointments/:id` | Detalle de un turno |
| `POST` | `/appointments` | Crear turno |
| `PATCH` | `/appointments/:id/cancel` | Cancelar turno |

**Crear turno — body:**
```json
{
  "client_id": 5,
  "service_id": 2,
  "date": "2026-03-25",
  "start_time": "14:00",
  "notes": "Quiero el diseño floral"
}
```

**Reglas de negocio:**
- No se pueden crear turnos en fechas pasadas
- Un cliente puede tener solo un turno activo por día
- El cliente siempre crea el turno para sí mismo (el `client_id` del token tiene prioridad); el admin puede crearlo para cualquier cliente

---

### Turnos — solo admin `🔒`

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/appointments/calendar?from=2026-03-01&to=2026-03-31` | Vista calendario por rango de fechas |
| `PATCH` | `/appointments/:id/confirm` | Confirmar turno |
| `PATCH` | `/appointments/:id/complete` | Completar turno (marcar como DONE) |
| `PATCH` | `/appointments/:id/final-price` | Actualizar precio final con extras |
| `PATCH` | `/appointments/:id/deposit` | Registrar seña recibida |

**Estados de un turno:**
```
PENDING → CONFIRMED → DONE
       ↘            ↗
        CANCELLED / ABSENT
```

---

### Servicios — solo admin `🔒`

| Método | Ruta | Descripción |
|--------|------|-------------|
| `POST` | `/services` | Crear servicio |
| `PATCH` | `/services/:id` | Actualizar servicio |
| `PATCH` | `/services/:id/toggle` | Activar/desactivar servicio |

---

### Horarios — solo admin `🔒`

| Método | Ruta | Descripción |
|--------|------|-------------|
| `PUT` | `/schedule/weekly` | Actualizar horario semanal |
| `POST` | `/schedule/blocked-slots` | Bloquear un horario específico |

---

### Clientes — solo admin `🔒`

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/clients` | Listar clientes con acumulado de turnos y total gastado |

**Respuesta:**
```json
[
  {
    "id": 3,
    "first_name": "María",
    "last_name": "López",
    "email": "maria@example.com",
    "phone": "1122334455",
    "active": true,
    "appointment_count": 8,
    "total_spent": 24500.00
  }
]
```

---

### Dashboard — solo admin `🔒`

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/dashboard/metrics` | Métricas del negocio |

**Respuesta — desglose por período:**
```json
{
  "today":  { "total": 3, "completed": 1, "cancelled": 0, "pending": 1, "confirmed": 1, "revenue": 3500.00, "deposits": 1500.00 },
  "week":   { ... },
  "month":  { ... },
  "year":   { ... }
}
```

---

## Roles

| Rol | Acceso |
|-----|--------|
| `client` | Registro, login, ver/crear/cancelar sus propios turnos, ver próximo turno |
| `admin` | Todo lo anterior + gestión completa de agenda, servicios, clientes y dashboard |
| `superadmin` | Igual que `admin` |

---

## Notificaciones por email

Al crear un turno se envían automáticamente (de forma asíncrona):
- **Email al cliente** con los datos del turno
- **Email al administrador** con los datos del cliente y el turno

Si las variables SMTP no están configuradas, los emails se loguean en consola (mock).

---

## Variables de entorno — referencia completa

| Variable | Descripción | Default |
|----------|-------------|---------|
| `DB_HOST` | Host de PostgreSQL | `localhost` |
| `DB_PORT` | Puerto de PostgreSQL | `5432` |
| `DB_USER` | Usuario de PostgreSQL | `postgres` |
| `DB_PASSWORD` | Contraseña de PostgreSQL | — |
| `DB_NAME` | Nombre de la base de datos | `shei_turnos` |
| `DB_SSLMODE` | Modo SSL (`disable` / `require`) | `disable` |
| `APP_ENV` | Entorno (`development` / `production`) | `development` |
| `APP_PORT` | Puerto HTTP | `8080` |
| `JWT_SECRET` | Clave secreta para firmar tokens | `changeme_in_production` |
| `SMTP_HOST` | Host del servidor SMTP | `smtp.gmail.com` |
| `SMTP_PORT` | Puerto SMTP | `587` |
| `SMTP_USER` | Usuario SMTP | — |
| `SMTP_PASS` | Contraseña SMTP | — |
| `ADMIN_EMAIL` | Email destino de notificaciones al admin | `admin@shei.com` |
