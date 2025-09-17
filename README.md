# 🔐 Auth Service

Sistema de autenticación JWT con Go, PostgreSQL y arquitectura hexagonal.

## 🚀 Características

- ✅ Autenticación JWT (Access + Refresh tokens)
- ✅ Registro y login de usuarios
- ✅ Gestión de perfiles
- ✅ Arquitectura hexagonal/Clean Architecture
- ✅ Base de datos PostgreSQL
- ✅ Middleware de autenticación
- ✅ Usuarios de prueba precargados
- ✅ Migraciones automáticas
- ✅ Logs estructurados
- ✅ Validaciones de entrada

## 📋 Requisitos

- Go 1.21+
- PostgreSQL 13+
- Make (opcional)

## 🛠️ Setup rápido

```bash
# 1. Clonar proyecto
git clone <tu-repo>
cd auth-service

# 2. Configurar variables
cp .env.example .env
# Editar .env con tus valores

# 3. Setup completo con Docker
make setup

# 4. Ejecutar
make run
```



## 📡 API Endpoints

### Públicos

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/api/v1/auth/register` | Registrar usuario |
| POST | `/api/v1/auth/login` | Iniciar sesión |
| POST | `/api/v1/auth/refresh` | Refrescar token |

### Protegidos (requieren JWT)

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/api/v1/users/profile` | Ver perfil |
| PUT | `/api/v1/users/profile` | Actualizar perfil |
| GET | `/api/v1/users` | Listar usuarios |
| GET | `/api/v1/users/{id}` | Ver usuario por ID |

## 🧪 Usuarios de prueba

```
admin@example.com / admin123
john@example.com  / john123
jane@example.com  / jane123
bob@example.com   / bob123
alice@example.com / alice123
```

## 📝 Ejemplos de uso

### Registro
```bash
curl -X POST     \
  -H "Content-Type: application/json" \
  -d '{
    "email": "nuevo@example.com",
    "password": "password123",
    "first_name": "Nuevo",
    "last_name": "Usuario",
    "username": "nuevousuario"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

### Ver perfil (requiere token)
```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer <tu-token-aqui>"
```

### Listar usuarios
```bash
curl -X GET "http://localhost:8080/api/v1/users?page=1&limit=10" \
  -H "Authorization: Bearer <tu-token-aqui>"
```

## 🗂️ Estructura del proyecto

```
auth-service/
├── cmd/server/main.go              # Punto de entrada
├── internal/
│   ├── config/config.go            # Configuración
│   ├── core/                       # Lógica de negocio
│   │   ├── domain/user.go          # Entidades
│   │   ├── ports/repositories.go   # Interfaces
│   │   └── services/               # Servicios
│   └── adapters/                   # Adaptadores
│       ├── http/                   # HTTP (handlers, middleware, routes)
│       └── database/postgres/      # PostgreSQL
├── pkg/                           # Utilidades reutilizables
│   ├── jwt/jwt.go                 # JWT utilities
│   └── hash/bcrypt.go             # Password hashing
├── scripts/                       # Scripts útiles
└── Makefile                       # Comandos automatizados
```

## 🔧 Comandos disponibles

```bash
make help         # Ver todos los comandos
make setup        # Setup inicial completo
make run          # Ejecutar aplicación
make test         # Ejecutar tests
make build        # Compilar aplicación
make migrate-up   # Ejecutar migraciones
make migrate-down # Revertir migraciones
make clean        # Limpiar archivos generados
```

## 🔒 Seguridad

- Contraseñas hasheadas con bcrypt
- Tokens JWT con expiración
- Refresh tokens para renovación segura
- Middleware de autenticación
- Validación de entrada
- Soft delete de usuarios

## 🔄 Flujo de autenticación

1. Usuario se registra o hace login
2. Sistema genera Access Token (24h) + Refresh Token (7 días)
3. Cliente usa Access Token para peticiones autenticadas
4. Cuando Access Token expira, usar Refresh Token para generar nuevos
5. Middleware valida tokens automáticamente

## 🌍 Variables de entorno

Ver `.env.example` para todas las configuraciones disponibles.

## 🐛 Troubleshooting

### Error de conexión a DB
- Verificar que PostgreSQL esté corriendo
- Revisar credenciales en `.env`
- Ejecutar `make docker-up` para levantar PostgreSQL

### Error de migraciones
- Ejecutar `make migrate-up`
- Si falla, verificar conexión a DB

### Token inválido
- Verificar que el token no haya expirado
- Usar refresh token para generar uno nuevo

## 📄 Licencia

Cristian Daniel Galindo Jimenez

## 🤝 Contribuir

1. Fork del proyecto
2. Crear feature branch
3. Commit cambios
4. Push al branch
5. Crear Pull Request