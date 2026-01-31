# Backend

This repository contains the backend for the GoBetterAuth Playground, a demo project showcasing authentication features for GoBetterAuth.

---

### Getting Started

`Prerequisites`:

- Go 1.25+
- PostgreSQL (used as the database for this playground)

---

### Installation

```bash
$ git clone https://github.com/GoBetterAuth/go-better-auth-playground.git
$ cd go-better-auth-playground/backend
$ go mod download && go mod tidy
```

---

### Configuration

Copy `.env.example` to `.env` and update environment variables as needed.

- To send emails you need to set up an SMTP provider. There are various cloud options you can choose from such as Resend, Sendgrid and more. But you can also use the following open-source project we've built here for local development: [LocalMailer](https://github.com/m-t-a97/localmailer)

---

### Running the Server

```bash
go run main.go
```

---

### Contributing

Contributions are welcome! Please open issues or submit pull requests.

---
