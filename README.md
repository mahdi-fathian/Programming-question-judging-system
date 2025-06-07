# Programming Question Judging System

A comprehensive system for managing programming contests, problems, and code submissions.

## Features

- User Management (Registration, Authentication, Authorization)
- Problem Management (CRUD operations for programming problems)
- Contest Management (Create contests, register participants)
- Code Submission and Evaluation
- Real-time Results using NATS Message Queue
- Admin Panel for System Management

## Tech Stack

### Backend
- Go (Gin Framework)
- SQLite Database
- GORM ORM
- NATS Message Queue
- JWT Authentication

### Frontend
- React
- TypeScript
- Monaco Editor
- Tailwind CSS

## Project Structure

```
.
├── backend/
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── database/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   ├── models/
│   │   └── services/
│   └── go.mod
└── frontend/
    ├── src/
    │   ├── components/
    │   ├── pages/
    │   ├── services/
    │   └── styles/
    └── package.json
```

## Setup

### Backend

1. Install Go 1.21 or later
2. Create `.env` file in backend directory:
   ```
   DB_PATH=judge.db
   JWT_SECRET=your_jwt_secret
   SERVER_PORT=8080
   NATS_URL=nats://localhost:4222
   ENVIRONMENT=development
   ```
3. Install dependencies:
   ```bash
   cd backend
   go mod tidy
   ```
4. Run the server:
   ```bash
   go run cmd/main.go
   ```

### Frontend

1. Install Node.js and npm
2. Install dependencies:
   ```bash
   cd frontend
   npm install
   ```
3. Start development server:
   ```bash
   npm start
   ```

## API Endpoints

### Authentication
- POST /api/auth/register
- POST /api/auth/login

### Problems
- GET /api/problems
- GET /api/problems/:id
- POST /api/problems
- PUT /api/problems/:id
- DELETE /api/problems/:id

### Contests
- GET /api/contests
- GET /api/contests/:id
- POST /api/contests
- PUT /api/contests/:id
- DELETE /api/contests/:id
- POST /api/contests/:id/register

### Submissions
- POST /api/submissions
- GET /api/submissions
- GET /api/submissions/:id

### Admin
- GET /api/admin/users
- PUT /api/admin/users/:id
- DELETE /api/admin/users/:id

## License

MIT 