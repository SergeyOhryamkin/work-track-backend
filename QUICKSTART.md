# Quick Start Guide

## 🚀 Getting Started

### 1. Start the Backend (Docker)

```bash
cd backend

# Start PostgreSQL and API
docker-compose up -d

# Run database migrations
make migrate-up

# Check logs
docker-compose logs -f api
```

The API will be running at `http://localhost:8080`

### 2. Test the API

```bash
# Health check
curl http://localhost:8080/health

# Register a user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Login (save the token from response)
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Create a task (replace TOKEN with your JWT)
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{"title":"Learn Go","description":"Complete backend","status":"in_progress"}'

# List tasks
curl http://localhost:8080/api/tasks \
  -H "Authorization: Bearer TOKEN"
```

## 📋 Common Commands

```bash
# Development
make run              # Run locally without Docker
make build            # Build binary
make test             # Run tests
make fmt              # Format code

# Docker
make docker-up        # Start services
make docker-down      # Stop services
make docker-logs      # View logs
make docker-build     # Build image

# Database
make migrate-up       # Apply migrations
make migrate-down     # Rollback migrations

# Cleanup
make clean            # Remove build artifacts
```

## 🔧 Environment Setup

1. Copy `.env.example` to `.env`
2. Update `JWT_SECRET` with a secure random string
3. Adjust `ALLOWED_ORIGINS` for your Vue 3 frontend URL

## 🌐 Vue 3 Frontend Integration

```javascript
// Example API call from Vue 3
const response = await fetch('http://localhost:8080/api/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ email, password })
})
const { token, user } = await response.json()

// Store token for subsequent requests
localStorage.setItem('token', token)

// Use token in protected requests
const tasks = await fetch('http://localhost:8080/api/tasks', {
  headers: { 'Authorization': `Bearer ${token}` }
})
```

## 📚 API Endpoints

### Public
- `POST /api/auth/register` - Create account
- `POST /api/auth/login` - Get JWT token

### Protected (requires JWT)
- `GET /api/tasks` - List tasks
- `POST /api/tasks` - Create task
- `GET /api/tasks/:id` - Get task
- `PUT /api/tasks/:id` - Update task
- `DELETE /api/tasks/:id` - Delete task

## 🐛 Troubleshooting

**Docker not running?**
```bash
# Start Docker Desktop or Docker daemon
# Then retry: docker-compose up -d
```

**Port 8080 already in use?**
```bash
# Change PORT in .env file
PORT=3001
```

**Database connection failed?**
```bash
# Check PostgreSQL is running
docker-compose ps

# View database logs
docker-compose logs postgres
```

## 📖 Next Steps

1. Review the [README.md](file:///Users/sergey/Documents/Projects/work_track/backend/README.md) for detailed documentation
2. Check [walkthrough.md](file:///Users/sergey/.gemini/antigravity/brain/91e6333b-f4b4-4d08-804a-7c760985e100/walkthrough.md) for learning concepts
3. Start building your Vue 3 frontend!
4. Add tests to the backend
5. Deploy to a cloud platform
