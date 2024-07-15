# TODO List API

REST API for managing tasks (TODO list) built with Go and PostgreSQL.

## Technologies

- Go 1.25.4
- PostgreSQL 16
- Chi Router
- Docker & Docker Compose

## Features

- Create tasks
- Get all tasks
- Get task by ID
- Full task update
- Toggle completion status
- Delete tasks
- Input validation

## API Endpoints
```
GET    /todos           - Get all tasks
POST   /todos           - Create new task
GET    /todos/{id}      - Get task by ID
PUT    /todos/{id}      - Update task
PATCH  /todos/{id}/status - Toggle completion status
DELETE /todos/{id}      - Delete task
```

## Project Structure
```
todoapp/
├── config/          - Configuration (.env handling)
├── database/        - PostgreSQL operations
├── handlers/        - HTTP handlers
├── models/          - Data structures
├── validation/      - Input validation
├── main.go          - Entry point
├── .env             - Environment variables
├── Dockerfile       - Application Docker image
└── docker-compose.yml - Container orchestration
```

## Running the Project

### Using Docker (recommended)
```bash
docker-compose up --build
```

Application will be available at `http://localhost:8080`

### Local Development

1. Install PostgreSQL and create database
2. Create `.env` file with connection settings
3. Run application:
```bash
go run main.go
```

## Request Examples

### Create task
```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy milk", "description": "At the corner store"}'
```

### Get all tasks
```bash
curl http://localhost:8080/todos
```

### Toggle status
```bash
curl -X PATCH http://localhost:8080/todos/1/status \
  -H "Content-Type: application/json" \
  -d '{"completed": true}'
```

## Database Schema
```sql
CREATE TABLE todos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```