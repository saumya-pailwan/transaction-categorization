# Autonomous Transaction Categorization Agent

- **backend/** — Go (chi router) API with categorization engine skeleton  
- **frontend/** — Next.js 14 (App Router) + TypeScript UI  

## Features

- Go backend with routing, handlers, DB stubs
- Categorization engine architecture (rules -> embeddings -> LLM -> validation)
- Next.js frontend showing transactions & rules
- Ready to integrate Plaid + LLMs

## Structure

```
backend/   # Go code
frontend/  # Next.js and TS UI
```

## Getting Started

### Backend

```bash
cd backend
go mod tidy
export DB_DSN="postgres://user:pass@localhost:5432/txagent?sslmode=disable"
go run ./cmd/api
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Navigate to: http://localhost:3000

---
