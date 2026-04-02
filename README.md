# Task Manager Backend

Backend para aplicación de gestión de tareas con arquitectura de microservicios, organizado por proyectos con sistema de usuarios y roles.

## Arquitectura

```
Cliente → API Gateway (Node.js :8080) → Task Service (Go :8081) → Firestore
```

- **API Gateway (Node.js/Express)**: Expone API pública, valida inputs, orquesta llamadas al servicio Go.
- **Task Service (Go)**: Maneja lógica de negocio, autenticación JWT, acceso a Firestore.
- **Base de datos**: Google Cloud Firestore (NoSQL).

## Estructura del proyecto

```
├── api-gateway/              # Node.js - API Gateway
│   ├── src/
│   │   ├── features/         # Rutas por dominio (auth, tasks, projects, users)
│   │   ├── shared/           # Middleware, config, utils
│   │   └── app.js            # Entry point
│   ├── Dockerfile
│   └── package.json
├── task-service/             # Go - Task Service
│   ├── cmd/server/           # Entry point
│   ├── internal/             # Lógica por dominio (auth, tasks, projects, users)
│   ├── pkg/                  # Paquetes compartidos (middleware, models, firestore)
│   ├── config/               # Configuración
│   ├── Dockerfile
│   └── go.mod
├── postman/                  # Colección Postman
├── docker-compose.yml
└── README.md
```

## Requisitos previos

- Go 1.23+
- Node.js 20+
- Proyecto en GCP con Firestore habilitado (Native mode)
- Service Account con rol "Cloud Datastore User"
- Docker y Docker Compose (opcional, para contenedores)

## Configuracion local

### 1. Clonar el repositorio

```bash
git clone <url-del-repo>
cd task-manager-backend
```

### 2. Configurar variables de entorno

**Task Service:**
```bash
cd task-service
cp .env.example .env
```

Editar `task-service/.env`:
```
PORT=8081
GCP_PROJECT_ID=tu-project-id
JWT_SECRET=tu-clave-secreta
JWT_EXPIRATION=24h
```

**API Gateway:**
```bash
cd api-gateway
cp .env.example .env
```

Editar `api-gateway/.env`:
```
PORT=8080
TASK_SERVICE_URL=http://localhost:8081
JWT_SECRET=tu-clave-secreta
```

> El JWT_SECRET debe ser el mismo en ambos servicios.

### 3. Configurar credenciales de GCP

```bash
export GOOGLE_APPLICATION_CREDENTIALS="/ruta/a/tu/service-account.json"
```

### 4. Iniciar los servicios

**Opcion A - Manual:**

Terminal 1 (Task Service):
```bash
cd task-service
go run cmd/server/main.go
```

Terminal 2 (API Gateway):
```bash
cd api-gateway
npm install
npm start
```

**Opcion B - Docker Compose:**
```bash
docker-compose up --build
```

## Endpoints de la API

### Auth (publicos)
| Metodo | Ruta | Descripcion |
|--------|------|-------------|
| POST | `/api/auth/register` | Registrar usuario |
| POST | `/api/auth/login` | Iniciar sesion |

### Tasks (autenticado)
| Metodo | Ruta | Descripcion |
|--------|------|-------------|
| GET | `/api/projects/:projectId/tasks` | Listar tareas (paginado) |
| GET | `/api/projects/:projectId/tasks/:taskId` | Obtener tarea |
| POST | `/api/projects/:projectId/tasks` | Crear tarea |
| PUT | `/api/projects/:projectId/tasks/:taskId` | Actualizar tarea |
| DELETE | `/api/projects/:projectId/tasks/:taskId` | Eliminar tarea |

### Projects (autenticado)
| Metodo | Ruta | Descripcion |
|--------|------|-------------|
| GET | `/api/projects` | Listar proyectos del usuario |
| GET | `/api/projects/:projectId` | Obtener proyecto |
| POST | `/api/projects` | Crear proyecto |
| PUT | `/api/projects/:projectId` | Actualizar proyecto |
| DELETE | `/api/projects/:projectId` | Eliminar proyecto |
| POST | `/api/projects/:projectId/members` | Agregar miembro |
| DELETE | `/api/projects/:projectId/members/:userId` | Quitar miembro |

### Users (solo admin)
| Metodo | Ruta | Descripcion |
|--------|------|-------------|
| GET | `/api/users` | Listar usuarios |
| PUT | `/api/users/:userId/role` | Cambiar rol |

### Parametros de paginacion (virtual scroll)
```
GET /api/projects/:projectId/tasks?limit=20&lastId=<id-ultima-tarea>
```

## Roles

| Rol | Permisos |
|-----|----------|
| admin | CRUD proyectos, gestionar miembros, todas las tareas, gestionar usuarios |
| member | CRUD tareas propias, crear proyectos, ver proyectos asignados |
| viewer | Solo lectura |

## Despliegue en GCP (Cloud Run)

### 1. Configurar gcloud

```bash
gcloud auth login
gcloud config set project tu-project-id
```

### 2. Desplegar Task Service

```bash
cd task-service
gcloud run deploy task-service \
  --source . \
  --region us-central1 \
  --set-env-vars "GCP_PROJECT_ID=tu-project-id,JWT_SECRET=tu-clave,JWT_EXPIRATION=24h" \
  --allow-unauthenticated
```

### 3. Desplegar API Gateway

```bash
cd api-gateway
gcloud run deploy api-gateway \
  --source . \
  --region us-central1 \
  --set-env-vars "TASK_SERVICE_URL=https://task-service-xxx.run.app,JWT_SECRET=tu-clave" \
  --allow-unauthenticated
```

> Reemplazar `https://task-service-xxx.run.app` con la URL que devuelve el paso anterior.

## Tecnologias

- **Go** - Task Service (lógica de negocio)
- **Node.js / Express** - API Gateway
- **Google Cloud Firestore** - Base de datos NoSQL
- **JWT** - Autenticación
- **bcrypt** - Hash de contraseñas
- **Docker** - Contenedores
- **GCP Cloud Run** - Despliegue
