# Estructura de Carpetas

```
Classplanner_2.0/
│
├── 📁 cmd/
│   └── 📁 api/
│       └── 📄 main.go                     # Punto de entrada: inicia Fiber, DB, rutas
│
├── 📁 internal/
│   ├── 📁 model/                          # Entidades del dominio (modelos)
│   │   ├── 📄 book.go                     # struct Book
│   │   ├── 📄 task.go                     # struct Task
│   │   └── 📄 user.go                     # struct User (Email, Password, Role, etc.)
│   │
│   ├── 📁 repository/                     # Capa de persistencia (repositorios SQL)
│   │   ├── 📄 repository.go               # Inicializa DB y agrupa repositorios
│   │   ├── 📄 task_repo.go                # CRUD de tasks
│   │   └── 📄 user_repo.go                # CRUD usuarios, búsqueda por email, etc.
│   │
│   ├── 📁 service/                        # Lógica de negocio (casos de uso)
│   │   ├── 📄 task_service.go             # Lógica de tareas (crear, listar, eliminar)
│   │   └── 📄 user_service.go             # Login, registro, validaciones, roles
│   │
│   ├── 📁 transport/                      # Capa HTTP / transporte
│   │   ├── 📁 users/
│   │   │   ├── 📄 handler.go              # Registro, login, logout
│   │   │   └── 📄 auth.go                 # Validación de sesiones / JWT
│   │   │
│   │   ├── 📁 tasks/
│   │   │   └── 📄 handler.go              # Endpoints CRUD de tareas
│   │   │
│   │   └── 📁 websocket/                  # (Opcional) soporte para WebSocket
│   │       └── 📄 hub.go
│   │
│   ├── 📁 infrastructure/                 # Servicios externos / configuración
│   │   ├── 📁 database/
│   │   │   ├── 📄 postgres.go             # Conexión y setup de la base de datos
│   │   │   └── 📁 migrations/             # Archivos SQL de migraciones
│   │   │
│   │   ├── 📁 logger/
│   │   │   └── 📄 logger.go               # Configuración del logger global
│   │   │
│   │   └── 📁 config/
│   │       └── 📄 config.go               # Carga y validación de variables .env
│   │
│   ├── 📁 middleware/                     # Middlewares HTTP
│   │   ├── 📄 auth.go                     # Valida cookie o token JWT
│   │   └── 📄 logging.go                  # Logs de requests/responses
│   │
│   └── 📁 security/                       # Funciones de seguridad reutilizables
│       ├── 📄 cookies.go                  # Manejo seguro de cookies de sesión
│       └── 📄 hash.go                     # Hashing de contraseñas (bcrypt)
│
├── 📁 pkg/                                # Paquetes utilitarios genéricos
│   ├── 📁 response/
│   │   └── 📄 fiber_response.go           # Respuestas JSON estandarizadas
│   │
│   ├── 📁 errors/
│   │   └── 📄 custom_error.go             # Tipos de error reutilizables
│   │
│   └── 📁 utils/
│       └── 📄 time_utils.go               # Funciones auxiliares (fechas, etc.)
│
├── 📁 tests/                              # Pruebas unitarias y de integración
│   ├── 📁 tasks/
│   │   └── 📄 task_service_test.go
│   └── 📁 users/
│       └── 📄 user_service_test.go
│
├── 📄 .env                                # Variables de entorno (DB_URL, JWT_SECRET, etc.)
├── 📄 .gitignore
├── 📄 main.go                             # Entrypoint (cmd, parametros de entrada, etc.)
├── 📄 go.mod
├── 📄 go.sum
└── 📄 README.md
```
# Bibliotecas utilizadas

* gorm.io/gorm
* gofiber/fiber/v2
* joho/godotenv
* go-deepseek/deepseek
