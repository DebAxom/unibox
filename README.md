<img src="./logo.svg" width="96" height="96" alt="Unibox" />

# Unibox

**The Unified Complaint Box for Universities**

A smart, production-ready campus issue management system that routes every student complaint to exactly the right department — automatically.

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![Fiber](https://img.shields.io/badge/Fiber-v3-00ACD7?style=flat-square)](https://gofiber.io)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=flat-square&logo=postgresql&logoColor=white)](https://postgresql.org)
[![Redis](https://img.shields.io/badge/Redis-7.2-DC382D?style=flat-square&logo=redis&logoColor=white)](https://redis.io)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind-v4-06B6D4?style=flat-square&logo=tailwindcss&logoColor=white)](https://tailwindcss.com)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)

---

## Overview

Universities have a complaint problem. Students rely on physical registers, informal emails, and word-of-mouth to report infrastructure failures, academic grievances, or welfare concerns. Complaints vanish. The wrong person sees them. Nothing gets tracked.

**Unibox** solves this end-to-end:

- Students file one complaint through a single, clean form
- A weighted keyword engine automatically routes it to the correct department or hostel block
- The right admin sees it immediately on their role-scoped dashboard
- Status updates flow back to the student in real time

No manual triage. No lost tickets. Full accountability.

---

## Features

### For Students
| Feature | Description |
|---|---|
| **Centralized Filing** | One form covers all categories — hostel, academics, mess, finance, student welfare |
| **Image Evidence Upload** | Attach photo proof up to 3 MB; stored server-side with UUID filenames |
| **Live Status Tracking** | Personal dashboard shows every ticket: `pending` → `progress` → `resolved` / `rejected` |
| **Real-time Notifications** | In-app notification centre updates the moment an admin acts on your complaint |
| **OTP-Verified Registration** | Email OTP (7-minute TTL via Redis) gates account creation |

### For Admins
| Feature | Description |
|---|---|
| **Smart Auto-Routing** | Complaints are classified and dispatched automatically — no manual assignment |
| **Department-Scoped Views** | Each admin only sees issues belonging to their own department |
| **Status Management** | Mark issues In-Progress, Resolved, or Rejected — with a reason sent to the student |
| **Hostel Block Isolation** | BH8 warden sees BH8 issues only; GH3 warden sees GH3 issues only |

---

## How the Routing Engine Works

The core of Unibox lives in [`api/utils/route.go`](api/utils/route.go).

When a student submits a complaint, `RouteComplain` scores the title and description against a weighted keyword map. Title text is doubled in weight to reflect its importance. Each matching keyword contributes its score to a department bucket. The highest-scoring bucket wins — as long as it clears a **confidence threshold of 2**.

```
"My wifi in BH8 is broken"
  → wifi  (+3, hostel)
  → BH8   (+2, hostel)
  → score: hostel = 5  ✓
  → hostel dept → hostel-bh8  (resolved via student's registered hostel block)
```

**Departments and example keywords:**

| Department | High-weight keywords |
|---|---|
| `hostel` | `wifi`(3), `mess`(3), `fan`(2), `water`(2), `food`(2) |
| `academic` | `exam`(3), `cgpa`(3), `sgpa`(3), `faculty`(2), `marks`(2) |
| `accounts` | `fee`(3), `payment`(3), `refund`(3) |
| `sw` (Student Welfare) | `harassment`(3), `ragging`(3), `scholarship`(3), `mental`(2) |

**Hostel block mapping** (from student's registered hostel):

```
abh → hostel-abh
bh8 → hostel-bh8
gh1 → hostel-gh1
gh3 → hostel-gh3
```

Typo-tolerance handles common misspellings (`wfi`, `wi-fi` → `wifi`). If no department scores above the threshold, the complaint routes to `other` for manual review.

---

## Tech Stack

| Layer | Technology |
|---|---|
| **Frontend** | [Nijor](https://nijor.dev) + Tailwind CSS v4 |
| **Backend** | Go 1.26 + [Fiber v3](https://gofiber.io) |
| **Database** | PostgreSQL 16 via `pgxpool` |
| **Cache / OTP** | Redis 7.2 |
| **Authentication** | JWT (HS256) — access tokens (15 min) + refresh tokens (15 days) + token version rotation |
| **Email** | SMTP via [`go-mail`](https://github.com/wneessen/go-mail) (Gmail) |
| **Infrastructure** | Docker Compose |

---

## Project Structure

```
unibox/
├── api/                        # Go backend
│   ├── main.go                 # Server bootstrap, routes
│   ├── handlers/
│   │   ├── api.go              # Profile, notifications endpoints
│   │   ├── auth.go             # OTP, register, login, logout, refresh
│   │   └── issues.go           # Create, get, status management
│   ├── middlewares/
│   │   └── auth.go             # JWT validation + token version check
│   ├── db/
│   │   ├── user.go             # User queries
│   │   ├── admin.go            # Admin queries
│   │   └── issue.go            # Issue queries
│   ├── models/
│   │   ├── user.go
│   │   ├── admin.go
│   │   └── issue.go
│   └── utils/
│       ├── route.go            # ★ Smart routing engine
│       ├── otp.go              # SMTP OTP sender
│       └── notification.go     # Notification writer
│
├── client/                     # Nijor + Tailwind frontend
│   ├── src/
│   │   ├── pages/
│   │   │   ├── auth/           # Login, register, OTP flow
│   │   │   ├── app/            # Student: dashboard, complaint form, notifications
│   │   │   └── dashboard/      # Admin: issues, analytics, past issues
│   │   ├── components/
│   │   │   ├── app/            # Student card, header, notification, issue
│   │   │   └── dashboard/      # Admin issue card, header, sidebar
│   │   └── layouts/            # app, auth, dashboard, default
│   ├── middlewares/
│   │   └── auth.js             # Server-side route protection (SSR guard)
│   └── assets/
│       ├── images/             # logo.svg, logo.png, auth.png
│       └── uploads/            # User-uploaded complaint images
│
├── sql/                        # Database schema
│   ├── users.sql
│   ├── admins.sql
│   ├── issues.sql
│   └── notifications.sql
│
└── docker-compose.yml          # PostgreSQL + Redis
```

---

## API Reference

### Auth — `/auth`

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/auth/otp` | Request a 4-digit OTP to the provided email |
| `POST` | `/auth/register` | Register a new student (requires valid OTP) |
| `POST` | `/auth/login` | Student login — returns access token, sets refresh cookie |
| `POST` | `/auth/login/admin` | Admin login — returns access token, sets refresh cookie |
| `POST` | `/auth/refresh` | Rotate access token using refresh cookie |
| `POST` | `/auth/logout` | Clear refresh cookie for current device |
| `POST` | `/auth/logoutall` | Increment `token_version` — invalidates all active sessions |

### API — `/api` *(JWT required)*

| Method | Endpoint | Role | Description |
|---|---|---|---|
| `GET` | `/api/` | both | Current user/admin profile |
| `GET` | `/api/notifications` | user | All notifications (marks unread as read) |
| `GET` | `/api/new/notification` | user | Check for unread notifications |
| `POST` | `/api/issue` | user | File a new complaint (multipart, image optional) |
| `GET` | `/api/issues` | user | All issues filed by this student |
| `GET` | `/api/issues/unresolved` | admin | Active issues for admin's department |
| `GET` | `/api/issues/resolved` | admin | Past resolved issues for admin's department |
| `GET` | `/api/issues/count` | admin | Count of unresolved issues (used for polling) |
| `PATCH` | `/api/issue/:id/progress` | admin | Mark issue In-Progress + notify student |
| `PATCH` | `/api/issue/:id/resolved` | admin | Mark issue Resolved + notify student |
| `PATCH` | `/api/issue/:id/reject` | admin | Reject issue with reason + notify student |

---

## Database Schema

```sql
-- Students
CREATE TABLE users (
  id           TEXT PRIMARY KEY,
  name         TEXT NOT NULL,
  email        TEXT UNIQUE NOT NULL,
  scholar_id   TEXT NOT NULL,
  password     TEXT NOT NULL,         -- bcrypt cost 12
  gender       TEXT,
  hostel       TEXT,                  -- used for hostel block routing
  token_version INT NOT NULL DEFAULT 0,
  created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Admins
CREATE TABLE admins (
  id            TEXT PRIMARY KEY,
  username      TEXT UNIQUE NOT NULL,
  password      TEXT NOT NULL,        -- bcrypt cost 12
  dept          TEXT NOT NULL,        -- e.g. "hostel-bh8", "academic"
  token_version INT NOT NULL DEFAULT 0,
  created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Issues
CREATE TABLE issues (
  id          TEXT PRIMARY KEY,       -- UUID
  issuer      TEXT,                   -- user.id
  title       TEXT NOT NULL,
  description TEXT NOT NULL,
  img         TEXT,                   -- filename or "null"
  status      TEXT NOT NULL,          -- pending | progress | resolved | rejected
  dept        TEXT NOT NULL,          -- routing result
  updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Notifications
CREATE TABLE notifications (
  id         SERIAL PRIMARY KEY,
  user_id    TEXT NOT NULL,
  issue_id   TEXT,
  dept       TEXT,
  title      TEXT NOT NULL,
  message    TEXT NOT NULL,
  read       BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## Getting Started

### Prerequisites

- [Go 1.26+](https://go.dev/dl/)
- [Docker + Docker Compose](https://docs.docker.com/get-docker/)
- [Bun](https://bun.sh) (for the Nijor frontend)
- A Gmail account with an [App Password](https://support.google.com/accounts/answer/185833)

### 1. Clone the repository

```bash
git clone https://github.com/your-org/unibox.git
cd unibox
```

### 2. Start the database services

```bash
docker compose up -d
```

This starts:
- **PostgreSQL 16** on `localhost:5001` — auto-runs all SQL scripts in `/sql`
- **Redis 7.2** on `localhost:5002`

### 3. Configure the backend

Create `api/.env`:

```env
DATABASE_URL=postgresql://admin:joiaiaxom@localhost:5001/unibox
REDIS_URL=localhost:5002

ACCESS_SECRET=your_access_secret_here
REFRESH_SECRET=your_refresh_secret_here

GMAIL_ADDRESS=your@gmail.com
GMAIL_APP_PASS=your_app_password_here
```

### 4. Run the backend

```bash
cd api
go mod download
go run .
```

The API server starts on **`:5000`**.

### 5. Configure the frontend

Create `client/.env`:

```env
JWT_SECRET=your_refresh_secret_here   # must match REFRESH_SECRET above
```

### 6. Run the frontend

```bash
cd client
bun install
bun run tw &          # Tailwind watcher
nijor dev        # Nijor dev server on :3000 ; use bun x nijor dev if using bun
```

Open [http://localhost:3000](http://localhost:3000).

---

## Pre-configured Admin Accounts

The SQL seed data in [`sql/admins.sql`](sql/admins.sql) creates the following admin accounts. Default password : **`joiaiaxom`** (bcrypt hashed).

| Username | Department |
|---|---|
| `dean.academic` | `academic` |
| `registrar` | `accounts` |
| `dean.sw` | `sw` |
| `warden.abh` | `hostel-abh` |
| `supervisor.abh` | `hostel-abh` |
| `warden.bh8` | `hostel-bh8` |
| `supervisor.bh8` | `hostel-bh8` |
| `warden.gh1` | `hostel-gh1` |
| `warden.gh3` | `hostel-gh3` |

> **Change all passwords before any production deployment.**

---

## Security

| Mechanism | Implementation |
|---|---|
| **Password hashing** | bcrypt at cost factor 12 |
| **OTP verification** | Cryptographically random 4-digit OTP, stored in Redis with 7-minute TTL |
| **JWT access tokens** | HS256, 15-minute expiry |
| **Refresh token rotation** | HTTPOnly cookie, 15-day expiry, rotated on every refresh |
| **Token versioning** | `token_version` in DB — `logoutall` increments it, instantly invalidating every active session |
| **Department isolation** | API middleware checks `admin.Department == issue.Dept` before every status mutation — not just at the UI layer |
| **Upload validation** | 3 MB size cap, UUID-named files prevent path traversal |

---

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `DATABASE_URL` | ✅ | PostgreSQL connection string |
| `REDIS_URL` | ✅ | Redis address (`host:port`) |
| `ACCESS_SECRET` | ✅ | JWT signing secret for access tokens |
| `REFRESH_SECRET` | ✅ | JWT signing secret for refresh tokens |
| `GMAIL_ADDRESS` | ✅ | Gmail address used to send OTPs |
| `GMAIL_APP_PASS` | ✅ | Gmail App Password (not your account password) |

---

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/your-feature`
3. Commit your changes: `git commit -m "feat: describe your change"`
4. Push and open a Pull Request

Please keep pull requests focused — one feature or fix per PR.

---

## License

[MIT](LICENSE.md) — free to use, modify, and deploy.

---

<div>
Made with love from Guwahati, Assam
</div>