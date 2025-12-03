# Backend

This repository contains the backend for the Go Better Auth Playground, a demo project showcasing authentication features for GoBetterAuth.

---

### Features

- User registration and login
- Email verification
- Password reset & change
- Email change

---

### Getting Started

`Prerequisites`:

- Go 1.25+
- PostgreSQL (used as the database for this playground)

---

### Installation

```bash
git clone https://github.com/OpenSource/GoBetterAuth/go-better-auth-playground.git
cd go-better-auth-playground/backend
go mod tidy
```

---

### Configuration

Copy `.env.example` to `.env` and update environment variables as needed.

- To send emails you need to set up a mailer service. We'd recommend using [LocalMailer](https://github.com/m-t-a97/localmailer) for local development so you can test out this playground with ease.

---

### Running the Server

```bash
go run main.go
```

---

### API Endpoints

| Method | Endpoint              | Description                                      |
| ------ | --------------------- | ------------------------------------------------ |
| POST   | `/sign-up`            | Sign up a new user                               |
| POST   | `/sign-in`            | Sign in a user                                   |
| POST   | `/email-verification` | Verify a user's email                            |
| POST   | `/reset-password`     | Send an email to reset a user's password         |
| POST   | `/change-password`    | Change a user's password                         |
| POST   | `/email-change`       | Send an email to confirm changing a user's email |
| GET    | `/me`                 | Get user profile                                 |

---

### Contributing

Contributions are welcome! Please open issues or submit pull requests.

---
