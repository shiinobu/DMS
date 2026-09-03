<div align="center">

# Device Management System (DMS)

</div>

A simple Device Management System (DMS) for registering devices, receiving
periodic heartbeats, monitoring device availability in real time, sending
online/offline notifications, and generating basic reports.

## Features

- Admin authentication using JWT
- Device registration and CRUD management
- Device simulator for 5 devices
- Heartbeat every 10 seconds
- Automatic ONLINE/OFFLINE detection
- Real-time status updates using WebSocket
- Browser notifications for status transitions
- Device monitoring summary
- CSV device export
- Docker Compose deployment

## Documentation

- [1. Project Overview](#1-project-overview)
- [2. Architecture](#2-architecture)
- [3. Requirements](#3-requirements)
- [4. Installation](#4-installation)
- [5. Configuration](#5-configuration)
- [6. Database](#6-database)
- [7. Admin User](#7-admin-user)
- [8. Run the Application](#8-run-the-application)
- [9. Simulator](#9-simulator)
- [10. Device Status & Realtime](#10-device-status--realtime)
- [11. API Reference](#11-api-reference)
- [12. Reports & CSV Export](#12-reports--csv-export)
- [13. Docker Deployment](#13-docker-deployment)
- [14. Troubleshooting](#14-troubleshooting)
- [License](#license)

---

# 1. Project Overview

The DMS monitors the availability of registered devices.

Each device sends a heartbeat every 10 seconds. When a heartbeat is received,
the device becomes `ONLINE`.

If no heartbeat is received for more than 30 seconds, the backend changes the
device to `OFFLINE`.

Status changes are broadcast through WebSocket to the Next.js dashboard.

```text
Device / Simulator
       |
       | Heartbeat
       v
  Go Backend
       |
       +----> PostgreSQL
       |
       +----> Status Monitor
       |
       +----> WebSocket
                  |
                  v
           Next.js Dashboard
                  |
                  v
             Notification
```

# 2. Architecture

The backend uses a simple layered architecture:

```text
Handler
   |
   v
Service
   |
   v
Repository
   |
   v
PostgreSQL
```

### Technology Stack

**Backend**

- Go 1.24+
- Gin
- pgxpool
- JWT
- bcrypt
- Gorilla WebSocket

**Frontend**

- Next.js 16
- React 19
- TypeScript
- Tailwind CSS

**Database**

- PostgreSQL 16+

### Project Structure

```text
DMS/
├── backend/
│   ├── cmd/
│   │   ├── server/
│   │   └── seed-admin/
│   ├── internal/
│   │   ├── auth/
│   │   ├── config/
│   │   ├── database/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   ├── models/
│   │   ├── repositories/
│   │   ├── services/
│   │   └── websocket/
│   ├── migrations/
│   ├── .env
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── src/
│   │   ├── app/
│   │   ├── components/
│   │   ├── lib/
│   │   └── types/
│   ├── .env.local
│   ├── Dockerfile
│   └── package.json
│
├── simulator/
│   ├── main.go
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
├── docker-compose.yml
└── README.md
```

# 3. Requirements

For local development:

- Go 1.24+
- Node.js 22+
- npm
- PostgreSQL 16+
- Git

Optional:

- Docker Desktop
- Docker Compose

Check versions:

```powershell
go version
node --version
npm --version
psql --version
docker --version
docker compose version
```

# 4. Installation

## Clone

```powershell
git clone <repository-url>
cd DMS
```

## Backend

```powershell
cd backend
go mod download
```

## Frontend

```powershell
cd ..\frontend
npm install
```

## Simulator

```powershell
cd ..\simulator
go mod download
```

# 5. Configuration

## Backend

Create:

```text
backend/.env
```

### Backend + PostgreSQL running locally

```env
APP_ENV=development
APP_PORT=8080

DATABASE_URL=postgres://postgres:postgres@localhost:5432/dms?sslmode=disable

JWT_SECRET=change-this-secret
JWT_EXPIRATION=24h

HEARTBEAT_TIMEOUT=30s
STATUS_CHECK_INTERVAL=5s

CORS_ORIGIN=http://localhost:3000

ADMIN_USERNAME=admin
ADMIN_PASSWORD=admin123
```

### Backend + PostgreSQL running in Docker

Use the Docker service name `postgres`:

```env
DATABASE_URL=postgres://postgres:postgres@postgres:5432/dms?sslmode=disable
```

The important rule is:

```text
Backend running on Windows  → localhost:5432
Backend running in Docker   → postgres:5432
```

Do not use the Docker hostname `postgres` from a backend running directly
on Windows.

## Frontend

Create:

```text
frontend/.env.local
```

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_WS_URL=ws://localhost:8080/api/v1/ws
```

These URLs are used by browser-side JavaScript, so keep `localhost` when
accessing the application through the host machine.

Do not use:

```text
http://backend:8080
ws://backend:8080
```

for `NEXT_PUBLIC_*` variables.

## Simulator

When running locally:

```env
DMS_API_URL=http://localhost:8080/api/v1
```

When running in Docker:

```env
DMS_API_URL=http://backend:8080/api/v1
```

Optional credentials:

```env
DMS_USERNAME=admin
DMS_PASSWORD=admin123
```

### Docker networking rule

```text
Host / Browser → localhost
Container → Docker service name
```

Examples:

```text
Browser → Backend       localhost:8080
Backend → PostgreSQL   postgres:5432
Simulator → Backend    backend:8080
```

# 6. Database

Create the database:

```sql
CREATE DATABASE dms;
```

Run migrations in numerical order:

```powershell
cd backend

psql -U postgres -d dms -f migrations/001_create_devices.sql
psql -U postgres -d dms -f migrations/002_create_users.sql
```

If additional migration files exist, execute them in numerical order.

The `devices` table contains:

```text
id
device_id
device_name
serial_number
os_version
ip_address
location
status
last_seen
last_online_at
last_offline_at
created_at
updated_at
```

Allowed device statuses:

```text
ONLINE
OFFLINE
```

Verify:

```sql
\dt

SELECT * FROM devices;

SELECT id, username, role FROM users;
```

# 7. Admin User

Create the initial admin:

```powershell
cd backend
go run ./cmd/seed-admin
```

Default development credentials:

```text
Username: admin
Password: admin123
```

The password is stored using bcrypt.

> Do not use the example credentials or JWT secret in production.

# 8. Run the Application

## Backend

Start PostgreSQL first, then:

```powershell
cd backend
go run ./cmd/server
```

Backend:

```text
http://localhost:8080
```

Health check:

```text
GET http://localhost:8080/health
```

## Frontend

```powershell
cd frontend
npm run dev
```

Open:

```text
http://localhost:3000
```

Login with:

```text
Username: admin
Password: admin123
```

## Simulator

```powershell
cd simulator
go run .
```

The simulator automatically:

1. Logs in using the admin account.
2. Checks whether each device already exists.
3. Registers missing devices.
4. Starts heartbeat for all devices.

---

# 9. Simulator

The simulator represents five devices:

```text
DMS-001  Office PC 001
DMS-002  Office PC 002
DMS-003  Office PC 003
DMS-004  Warehouse PC 001
DMS-005  IT Administrator PC
```

Device information:

```text
  DMS-001
  Serial:  SN-DMS-001
  OS:      Windows 11 Pro
  IP:      192.168.1.101
  Location: Office - Floor 1

  DMS-002
  Serial:  SN-DMS-002
  OS:      Windows 11 Pro
  IP:      192.168.1.102
  Location: Office - Floor 2

  DMS-003
  Serial:  SN-DMS-003
  OS:      Windows 10 Pro
  IP:      192.168.1.103
  Location: Office - Floor 3

  DMS-004
  Serial:  SN-DMS-004
  OS:      Windows 11 Pro
  IP:      192.168.1.104
  Location: Warehouse

  DMS-005
  Serial:  SN-DMS-005
  OS:      Windows 11 Pro
  IP:      192.168.1.105
  Location: IT Room
```

## Registration

The simulator first authenticates:

```http
POST /api/v1/auth/login
```

It receives a JWT and uses it for device management requests.

For every device:

```text
GET /api/v1/devices/{device_id}
        |
        +-- 200 → already exists → skip
        |
        +-- 404 → POST /api/v1/devices
```

This makes the simulator safe to restart without creating duplicate
device records.

## Heartbeat

After registration, every simulated device sends:

```http
POST /api/v1/devices/{device_id}/heartbeat
```

The heartbeat is sent immediately and then every:

```text
10 seconds
```

Stop the simulator with:

```text
Ctrl + C
```

When heartbeat stops, the backend eventually changes the devices to
`OFFLINE`.

# 10. Device Status & Realtime

## Heartbeat

When a heartbeat is received:

```text
status    = ONLINE
last_seen = NOW()
```

A normal heartbeat:

```text
ONLINE → ONLINE
```

does not generate a notification.

If the previous state was `OFFLINE`:

```text
OFFLINE → ONLINE
```

a status-change event is generated.

## Offline Detection

The backend runs a background status monitor every:

```text
5 seconds
```

A device becomes `OFFLINE` when:

```text
current_time - last_seen > 30 seconds
```

The monitor checks the database because a disconnected device cannot report
its own offline state.

```text
OFFLINE
   |
   | heartbeat
   v
ONLINE
   |
   | no heartbeat > 30s
   v
OFFLINE
```

## WebSocket

Endpoint:

```text
ws://localhost:8080/api/v1/ws
```

Event:

```text
DEVICE_STATUS_CHANGED
```

Example:

```json
{
  "type": "DEVICE_STATUS_CHANGED",
  "device_id": "DMS-003",
  "status": "OFFLINE",
  "device": {
    "device_id": "DMS-003",
    "device_name": "Office PC 003",
    "status": "OFFLINE"
  }
}
```

The dashboard uses the event to:

1. Update the affected device.
2. Update dashboard statistics.
3. Display a notification.
4. Keep the WebSocket connection open.
5. Reconnect after a disconnection.

Reconnect delay:

```text
3 seconds
```

# 11. API Reference

Base URL:

```text
http://localhost:8080/api/v1
```

## Authentication

```http
POST /auth/login
```

Request:

```json
{
  "username": "admin",
  "password": "admin123"
}
```

## Device Management

Protected by:

```http
Authorization: Bearer <JWT>
```

Endpoints:

```text
POST   /devices
GET    /devices
GET    /devices/:device_id
PUT    /devices/:device_id
DELETE /devices/:device_id
```

Only `OFFLINE` devices can be deleted.

## Heartbeat

```http
POST /devices/:device_id/heartbeat
```

Example:

```json
{
  "ip_address": "192.168.1.101"
}
```

## Reports

```http
GET /reports/summary
GET /reports/devices/export
```

Both report endpoints require admin authentication.

# 12. Reports & CSV Export

## Report Summary

The report contains:

```text
Total Devices
Online Devices
Offline Devices
Last Online
Last Offline
```

Endpoint:

```http
GET /api/v1/reports/summary
```

The frontend displays timestamps using:

```text
Asia/Jakarta
```

## CSV Export

Endpoint:

```http
GET /api/v1/reports/devices/export
```

Authentication:

```http
Authorization: Bearer <JWT>
```

CSV columns:

```text
ID
Device ID
Device Name
Serial Number
OS Version
IP Address
Location
Status
Last Seen
Last Online
Last Offline
Created At
Updated At
```

Timestamps are formatted for readability, for example:

```text
2026-09-03 21:30:25
```

# 13. Docker Deployment

The Docker Compose stack contains:

```text
PostgreSQL
Backend
Frontend
Simulator
```

Start:

```powershell
docker compose up -d --build
```

Check:

```powershell
docker compose ps
```

Logs:

```powershell
docker compose logs -f
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f simulator
```

Stop:

```powershell
docker compose down
```

Restart:

```powershell
docker compose restart
```

PostgreSQL data is stored in:

```text
postgres_data
```

Avoid:

```powershell
docker compose down -v
```

unless you intentionally want to remove the PostgreSQL data.

### Docker addresses

```text
Backend → PostgreSQL
postgres:5432

Simulator → Backend
backend:8080

Browser → Backend
localhost:8080

Browser → WebSocket
localhost:8080
```

The rule is simple:

> Docker service names are for container-to-container communication.
> `localhost` is used by the host/browser through published ports.

# 14. Troubleshooting

## Backend cannot connect to PostgreSQL

Check where the backend runs.

Local backend:

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/dms?sslmode=disable
```

Docker backend:

```env
DATABASE_URL=postgres://postgres:postgres@postgres:5432/dms?sslmode=disable
```

## Frontend CORS error

Backend:

```env
CORS_ORIGIN=http://localhost:3000
```

Frontend:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

Restart the backend and Next.js after changing environment variables.

## Login returns 401

Verify the admin:

```sql
SELECT id, username, role
FROM users
WHERE username = 'admin';
```

If necessary:

```powershell
cd backend
go run ./cmd/seed-admin
```

## Dashboard has no realtime updates

Verify:

```text
ws://localhost:8080/api/v1/ws
```

Check the browser console and backend logs.

The dashboard must connect to WebSocket when mounted.

## Devices remain OFFLINE

Run the simulator:

```powershell
cd simulator
go run .
```

Check:

```sql
SELECT
    device_id,
    status,
    last_seen
FROM devices
ORDER BY device_id;
```

## Devices do not become OFFLINE

Verify:

```env
HEARTBEAT_TIMEOUT=30s
STATUS_CHECK_INTERVAL=5s
```

Stop the simulator and wait approximately 30–35 seconds.

## Docker simulator cannot connect to backend

For a local simulator:

```env
DMS_API_URL=http://localhost:8080/api/v1
```

For a Docker simulator:

```env
DMS_API_URL=http://backend:8080/api/v1
```

Check:

```powershell
docker compose ps backend
docker compose logs simulator
```

## Docker frontend cannot connect to backend

Browser-side configuration must remain:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_WS_URL=ws://localhost:8080/api/v1/ws
```

Do not use `backend:8080` in `NEXT_PUBLIC_*`.

## CSV export returns 401

Make sure the request contains:

```http
Authorization: Bearer <JWT>
```

The JWT must have the `ADMIN` role.

## PostgreSQL data disappeared

Check Docker volumes:

```powershell
docker volume ls
```

Avoid:

```powershell
docker compose down -v
```

because `-v` removes the PostgreSQL volume.

---

# License

This project is intended as a simple Device Management System project for
development, learning, testing, and demonstration purposes.
