# Backend

Fiber + GORM API for the hotel management system.

## Stack

- Go
- Fiber
- GORM
- SQLite or PostgreSQL from `.env`
- JWT auth

## Database selection

The backend can use either SQLite or PostgreSQL.

### SQLite

```env
PORT=8080
DB_DRIVER=sqlite
SQLITE_PATH=./hotel_management.db
JWT_SECRET=super-secret-hotel-key
```

### PostgreSQL

```env
PORT=8080
DB_DRIVER=postgres
DATABASE_URL=host=localhost user=postgres password=postgres dbname=hotel_management port=5432 sslmode=disable
JWT_SECRET=super-secret-hotel-key
```

`DB_DRIVER` accepts:

- `sqlite`
- `postgres`
- `postgresql`

For SQLite, the app uses a pure Go driver, so it works on common Windows 32-bit and 64-bit setups without requiring a separate native SQLite installation.

## Run locally

```bash
go mod tidy
go run .
```

Default API URL:

- `http://localhost:8080`

Health check:

- `GET /api/health`

## What the backend does

- Auth and JWT login
- Guests, rooms, and bookings
- Hotel settings storage
- Hotel image upload and replacement
- Income categories
- Income records and receipt metadata
- Automatic schema migration and seed data

## Important folders

- `config/` app configuration and `.env` loading
- `database/` database connection and seed data
- `handlers/` Fiber route handlers
- `middleware/` JWT protection helpers
- `models/` GORM models
- `uploads/` runtime hotel branding uploads

## Notes

- SQLite database files are ignored by `backend/.gitignore`.
- Uploaded hotel images are stored under `backend/uploads/`.
- On startup the app runs auto-migrations and seeds starter records if the database is empty.
