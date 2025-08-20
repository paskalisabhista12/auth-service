# Auth Service

---

## 📌 About

Auth Service written in **Golang**

---

## ✨ Features

-   RESTful API with [Gin](https://github.com/gin-gonic/gin)
-   JWT Authentication
-   Database integration
-   Dockerized for deployment
-   Configurable via `.env`

---

## 🛠 Tech Stack

-   **Language**: Go 1.24.5
-   **Framework**: Gin
-   **Database**: PostgreSQL
-   **Cache**: Redis

---

## 🗄️ Database Migrations

We use [golang-migrate/migrate](https://github.com/golang-migrate/migrate) to manage schema changes.

---

### ⬆️ Migrate Up

Apply migrations:

```sh
# Apply all pending migrations
migrate -path ./db/migrations -database "$DATABASE_URL" up

# Apply only the next migration
migrate -path ./db/migrations -database "$DATABASE_URL" up 1

# Rollback the last migration
migrate -path ./db/migrations -database "$DATABASE_URL" down 1

# Rollback multiple steps
migrate -path ./db/migrations -database "$DATABASE_URL" down 3

# Rollback all migrations (⚠️ dangerous)
migrate -path ./db/migrations -database "$DATABASE_URL" down
```
