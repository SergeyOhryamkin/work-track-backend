# Work Track API – Frontend Integration Guide

Updated for the current codebase (chi router, JWT auth, SQLite). Everything is under `/api`.

## Base URL
```
http://localhost:8080/api
```

## Auth Basics
- JWT is returned on register/login.
- Send on protected routes: `Authorization: Bearer <token>`.
- Sessions are automatically tracked with device, platform, and IP information.

## Endpoints

### Auth

**POST** `/auth/register`
- Body:
```json
{
  "login": "johndoe",
  "password": "secret123"
}
```
- Session metadata is automatically captured from request headers (User-Agent, IP address).
- Responses:
  - `201 Created`
  ```json
  {
    "token": "<jwt>",
    "user": {
      "id": 1,
      "login": "johndoe",
      "first_name": "",
      "last_name": "",
      "avatar": "",
      "created_at": "2024-01-20T10:00:00Z",
      "updated_at": "2024-01-20T10:00:00Z"
    },
    "session_id": 1
  }
  ```
  - `409 Conflict` if login already exists.
  - `400 Bad Request` for missing/invalid fields (password ≥ 6 chars).

**POST** `/auth/login`
- Body:
```json
{ "login": "johndoe", "password": "secret123" }
```
- Session metadata is automatically captured from request headers.
- Responses:
  - `200 OK`
  ```json
  {
    "token": "<access_jwt>",
    "refresh_token": "<refresh_jwt>",
    "user": {
      "id": 1,
      "login": "johndoe",
      ...
    },
    "session_id": 2
  }
  ```
  - `401 Unauthorized` if credentials are wrong.

**POST** `/auth/refresh`
- Body:
```json
{
  "refresh_token": "<refresh_jwt>"
}
```
- Exchanges a valid refresh token for a new access token and a new refresh token.
- Responses:
  - `200 OK`
  ```json
  {
    "token": "<new_access_jwt>",
    "refresh_token": "<new_refresh_jwt>",
    "user": { ... },
    "session_id": 2
  }
  ```
  - `401 Unauthorized` if refresh token is invalid, expired, or session is closed.

**POST** `/auth/logout` (Protected)
- Requires: `Authorization: Bearer <token>`
- Body:
```json
{
  "session_id": 1
}
```
- Completes the session by recording logout time and calculating session duration.
- Responses:
  - `200 OK`
  ```json
  {
    "message": "Logged out successfully"
  }
  ```
  - `404 Not Found` if session doesn't exist or already closed.
  - `401 Unauthorized` if token is missing or invalid.

### Track Items (all require `Authorization: Bearer <token>`)

**POST** `/track-items`
- Body (all required):
```json
{
  "type": "regular",
  "emergency_call": false,
  "holiday_call": false,
  "working_hours": 8.5,
  "working_shifts": 1.0,
  "date": "2024-01-20T09:00:00Z" // RFC3339
}
```
- Responses:
  - `201 Created` with created item.
  - `400 Bad Request` on validation/date parse errors.

**GET** `/track-items`
- Optional query: `start_date=YYYY-MM-DD&end_date=YYYY-MM-DD` (end is inclusive, extended to end-of-day).
- Responses:
  - `200 OK` with array of caller’s items.

**GET** `/track-items/{id}`
- Responses:
  - `200 OK` with the item.
  - `404 Not Found` if missing.
  - `403 Forbidden` if the item is not owned by caller.

**PUT** `/track-items/{id}`
- Body: all fields optional; send only what you change.
```json
{
  "type": "overtime",
  "emergency_call": true,
  "holiday_call": false,
  "working_hours": 10,
  "working_shifts": 1,
  "date": "2024-01-21T09:00:00Z"
}
```
- Responses:
  - `200 OK` with updated item.
  - `400 Bad Request` on validation/date parse errors.
  - `404 Not Found` if missing.
  - `403 Forbidden` if not owned by caller.

**DELETE** `/track-items/{id}`
- Responses:
  - `204 No Content` on success.
  - `404 Not Found` if missing.
  - `403 Forbidden` if not owned by caller.

## Response Models
```json
// User
{
  "id": 1,
  "first_name": "John",
  "last_name": "Doe",
  "avatar": "https://example.com/avatar.jpg",
  "login": "johndoe",
  "created_at": "2024-01-20T10:00:00Z",
  "updated_at": "2024-01-20T10:00:00Z"
}

// TrackItem
{
  "id": 1,
  "user_id": 1,
  "type": "regular",
  "emergency_call": false,
  "holiday_call": false,
  "working_hours": 8.5,
  "working_shifts": 1.0,
  "date": "2024-01-20T09:00:00Z",
  "created_at": "2024-01-20T10:00:00Z",
  "updated_at": "2024-01-20T10:00:00Z"
}
```

## Error Shapes (examples)
```json
{ "error": "Invalid request body" }
{ "error": "Authorization header required" }
{ "error": "Invalid or expired token" }
{ "error": "Access denied" }
{ "error": "Track item not found" }
{ "error": "Login already exists" }
```

## curl Quickstart
```bash
# register (returns token and session_id)
RESPONSE=$(curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"johndoe","password":"secret123"}')

TOKEN=$(echo $RESPONSE | jq -r .token)
SESSION_ID=$(echo $RESPONSE | jq -r .session_id)

# login (creates a new session)
RESPONSE=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"johndoe","password":"secret123"}')

TOKEN=$(echo $RESPONSE | jq -r .token)
SESSION_ID=$(echo $RESPONSE | jq -r .session_id)

# create track item
curl -X POST http://localhost:8080/api/track-items \
  -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" \
  -d '{"type":"regular","emergency_call":false,"holiday_call":false,"working_hours":8,"working_shifts":1,"date":"2024-01-20T09:00:00Z"}'

# logout (completes the session)
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" \
  -d "{\"session_id\":$SESSION_ID}"
```

## Minimal Frontend (fetch)
```js
const BASE = 'http://localhost:8080/api';

async function register(payload) {
  const res = await fetch(`${BASE}/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });
  if (!res.ok) throw await res.json();
  return res.json(); // { token, user, session_id }
}

async function login(payload) {
  const res = await fetch(`${BASE}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });
  if (!res.ok) throw await res.json();
  return res.json(); // { token, user, session_id }
}

async function logout(token, sessionId) {
  const res = await fetch(`${BASE}/auth/logout`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ session_id: sessionId }),
  });
  if (!res.ok) throw await res.json();
  return res.json(); // { message: "Logged out successfully" }
}

async function listTrackItems(token, range) {
  const qs = range?.start && range?.end
    ? `?start_date=${range.start}&end_date=${range.end}` : '';
  const res = await fetch(`${BASE}/track-items${qs}`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) throw await res.json();
  return res.json();
}
```

## Notes
- Dates: send RFC3339 for `date` (`2024-01-20T09:00:00Z`). Date-range queries use `YYYY-MM-DD`.
- Ownership: all track-item routes are scoped to the authenticated user; 403 is returned if accessing another user’s item.
- Config: set `JWT_SECRET` and `DB_PATH` (.env). Default DB is SQLite file.
- Session Tracking: Every login/register creates a session record with:
  - Device type (mobile, tablet, desktop)
  - Platform (ios, android, windows, macos, linux, web)
  - User agent string
  - IP address (handles X-Forwarded-For and X-Real-IP headers)
  - Login timestamp
  - Logout timestamp and session duration (calculated on logout)
- Session ID: Save the `session_id` returned from login/register to use when calling logout.
