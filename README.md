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
SEED_ADMIN_NAME=System Admin
SEED_ADMIN_EMAIL=admin@hotel.local
SEED_ADMIN_PASSWORD=admin123
SEED_ADMIN_ROLE=admin
```

### PostgreSQL

```env
PORT=8080
DB_DRIVER=postgres
DATABASE_URL=host=localhost user=postgres password=postgres dbname=hotel_management port=5432 sslmode=disable
JWT_SECRET=super-secret-hotel-key
SEED_ADMIN_NAME=System Admin
SEED_ADMIN_EMAIL=admin@hotel.local
SEED_ADMIN_PASSWORD=admin123
SEED_ADMIN_ROLE=admin
```

`DB_DRIVER` accepts:

- `sqlite`
- `postgres`
- `postgresql`

For SQLite, the app uses a pure Go driver, so it works on common Windows 32-bit and 64-bit setups without requiring a separate native SQLite installation.

The backend only auto-seeds one admin user now. The seeded admin account comes entirely from the `SEED_ADMIN_*` values in `.env`.

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
- Automatic schema migration
- Admin-only bootstrap seed from `.env`

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
- On startup the app runs auto-migrations and, if there are no users yet, creates only the admin user from `.env`.
