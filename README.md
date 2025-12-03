# Go Better Auth Playground

This repository showcases a full-stack authentication system using Go for the backend and Next.js for the frontend. It is designed as a playground to explore secure authentication flows, including user registration, login, and protected routes powered by GoBetterAuth.

## Tech Stack

- Frontend: Next.js
- Backend: Go (Golang), Echo (for this example)
- Library: [GoBetterAuth](https://github.com/GoBetterAuth/go-better-auth)

## Getting Started

1. **Clone the repository**
2. **Docker Compose**
   - Ensure Docker and Docker Compose are installed.
   - Copy `docker-compose.env.example` to `docker-compose.env`.
   - Update environment variables in `docker-compose.env` as needed.
   - Start services: `docker compose up -d`
3. **Backend Setup**
   - Copy `backend/.env.example` to `backend/.env`.
   - Update environment variables in `backend/.env`.
   - Install dependencies: `go mod tidy`
   - Start the backend: `go run main.go`
4. **Frontend Setup**
   - Install dependencies: `pnpm install`
   - Copy `frontend/.env.local.example` to `frontend/.env.local`.
   - Update environment variables in `frontend/.env.local`.
   - Start the frontend: `pnpm dev`

This repository will be updated as GoBetterAuth evolves, with additional examples demonstrating integration in various scenarios.

For comprehensive documentation and usage instructions, visit the [GoBetterAuth Docs](https://go-better-auth.vercel.app/docs).

---
