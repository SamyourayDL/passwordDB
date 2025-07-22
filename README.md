# Simple Password Storage

Simple password storage with functionality to add and delete users and their passwords.

# TechStack 
- Go
- slog
- sqlx
- PostgreSQL(Docker container)
- AES-256
- chi

## 📘 API

Документация доступна:
- [Swagger-file](https://samyouraydl.github.io/passwordDB/)

## Methods
- `POST /user/{user_name}` — register user
- `POST /password/{user_name}` — add password to a user
- `GET /user/{user_name}` — get all user's passwords
- `GET /password/{user_name}` — get user's passwords from a specific service
- `DELETE /user/{user_name}` — delete user and all his passwords
- `DELETE /password/{user_name}` — delete user's password
