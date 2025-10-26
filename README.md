# fiber-micro
*Microservice API applying clean architecture, security, OpenApi, tracing, etc*

## Dependencies
Golang, Docker, Make, [Swag tool](https://github.com/swaggo/swag)

## Features
- Architecture
    - Clean architecture (adapters and ports)
    - Custom Messages and Errors
    - Pagination and Ordering
    - DB Migrator (scripts)
- Go 1.25
- Libraries
    - Web: Fiber
    - ORM: Gorm
    - Security: JWT
    - Validations: Go Playground Validator
    - Unit Test: Testify
    - DB: Postgres
    - Tracing: Opentelemetry
    - Test: Testcontainers
    - OpenAPI: Fiber Swagger
    - Environment: Godot
- Distributed tracing
    - OpenTelemetry and Jaeger
- Swagger
    - Swaggo & Fiber Swagger
    - Auditory (using Gorm)
- Database
    - Postgres
    - Testcontainers

## Files
- [Dockerfile](https://github.com/javiorfo/fiber-micro/tree/master/Dockerfile)

## Enviroment variables
```bash
# SQL
DATABASE_URL=postgres://admin:admin@localhost:5432/db_users
SHOW_SQL_INFO=false

# Jaeger tracing endpoint
TRACING_HOST=localhost:4318

# Swagger
SWAGGER_ENABLED=true

# Security JWT
SECURITY_ENABLED=true
JWT_DURATION=300
JWT_SECRET_KEY=UUID
JWT_ISSUER=https://acme.com
JWT_AUDIENCE=https://acme.com
```

## Usage
- Executing `make help` all the available commands will be listed. 
- Also the standard Go commands could be used, like `go run main.go`

## Services
- **Create users** POST: /users
- **Get users** GET: /users
- **Login** POST: /users/login

---

### Donate
- **Bitcoin** [(QR)](https://raw.githubusercontent.com/javiorfo/img/master/crypto/bitcoin.png)  `1GqdJ63RDPE4eJKujHi166FAyigvHu5R7v`
- [Paypal](https://www.paypal.com/donate/?hosted_button_id=FA7SGLSCT2H8G)
