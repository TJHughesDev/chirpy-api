# Chirpy API

A robust backend service for a microblogging platform, featuring secure user authentication, chirp management, and administrative tracking.

---

## 🚀 Getting Started

### Prerequisites

- **Go**: 1.22+
- **PostgreSQL**: 15+
- **Goose**: For database migrations
- **SQLC**: For type-safe SQL generation

### Environment Setup

Create a `.env` file in the root directory:

```env
DB_URL=postgres://user:password@localhost:5432/pulse_db?sslmode=disable
JWT_SECRET=your_super_secret_random_string
PLATFORM=dev
POLKA_KEY=your_api_key_for_webhooks
```

### Installation & Run

1. **Migrations**:
   ```bash
   cd sql/schema
   goose postgres <YOUR_DB_URL> up
   ```
2. **Build and Run**:
   ```bash
   go build -o pulse && ./pulse
   ```

---

## 🔑 Authentication

Pulse uses a dual-token system for security and user persistence:

- **Access Token (JWT)**: Short-lived (1 hour). Sent via `Authorization: Bearer <TOKEN>`.
- **Refresh Token**: Long-lived (60 days). Used to generate new Access Tokens.
- **API Key**: Used for Polka webhooks via `Authorization: ApiKey <KEY>`.

---

## 📡 API Reference

### 👤 Users

| Method | Endpoint     | Auth | Description                                     |
| :----- | :----------- | :--- | :---------------------------------------------- |
| `POST` | `/api/users` | None | Create a new user account.                      |
| `PUT`  | `/api/users` | JWT  | Update email and password for the current user. |

### 🔐 Auth & Sessions

| Method | Endpoint       | Auth    | Description                                          |
| :----- | :------------- | :------ | :--------------------------------------------------- |
| `POST` | `/api/login`   | None    | Authenticate and receive Access & Refresh tokens.    |
| `POST` | `/api/refresh` | Refresh | Provide a Refresh Token to receive a new Access JWT. |
| `POST` | `/api/revoke`  | Refresh | Revoke the provided Refresh Token (Logout).          |

### 🐤 Chirps

| Method   | Endpoint           | Auth | Description                                                   |
| :------- | :----------------- | :--- | :------------------------------------------------------------ |
| `POST`   | `/api/chirps`      | JWT  | Create a chirp (Max 140 chars. Profanity filtered).           |
| `GET`    | `/api/chirps`      | None | List all chirps. Supports `?author_id=UUID` and `?sort=desc`. |
| `GET`    | `/api/chirps/{id}` | None | Retrieve a specific chirp by ID.                              |
| `DELETE` | `/api/chirps/{id}` | JWT  | Delete a chirp (Must be the author).                          |

### 🛠 Admin & Webhooks

| Method | Endpoint              | Auth   | Description                                       |
| :----- | :-------------------- | :----- | :------------------------------------------------ |
| `GET`  | `/api/healthz`        | None   | Service health check.                             |
| `GET`  | `/admin/metrics`      | None   | View total fileserver hits (HTML).                |
| `POST` | `/admin/reset`        | Local  | Reset all users and hit counters (Dev mode only). |
| `POST` | `/api/polka/webhooks` | ApiKey | Upgrade a user to "Chirpy Red" status.            |

---

## 📦 Data Structures

### User Object

```json
{
  "id": "uuid",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "email": "string",
  "is_chirpy_red": "boolean"
}
```

### Chirp Object

```json
{
  "id": "uuid",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "user_id": "uuid",
  "body": "string"
}
```

---

## 🛠 Internal Architecture

- **Password Hashing**: Argon2id via `internal/auth`.
- **Database Layer**: Automated SQL mapping using `sqlc`.
- **Validation**: Chirps are automatically scrubbed for "bad words" (`kerfuffle`, `sharbert`, `fornax`) before persistence.
