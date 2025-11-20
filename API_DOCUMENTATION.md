# Work Track API Documentation

## Overview

The Work Track API allows users to register, authenticate, and track their work items with detailed information about working hours, shifts, and special conditions.

## Base URL

```
http://localhost:8080/api
```

## Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

---

## Endpoints

### Authentication

#### Register a New User

**POST** `/api/auth/register`

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "login": "johndoe",
  "password": "password123",
  "avatar": "https://example.com/avatar.jpg" // optional
}
```

**Response:** `201 Created`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "avatar": "https://example.com/avatar.jpg",
    "login": "johndoe",
    "created_at": "2024-01-20T10:00:00Z",
    "updated_at": "2024-01-20T10:00:00Z"
  }
}
```

#### Login

**POST** `/api/auth/login`

**Request Body:**
```json
{
  "login": "johndoe",
  "password": "password123"
}
```

**Response:** `200 OK`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "avatar": "https://example.com/avatar.jpg",
    "login": "johndoe",
    "created_at": "2024-01-20T10:00:00Z",
    "updated_at": "2024-01-20T10:00:00Z"
  }
}
```

---

### Track Items

All track item endpoints require authentication.

#### Create a Track Item

**POST** `/api/track-items`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "type": "regular",
  "emergency_call": false,
  "holiday_call": false,
  "working_hours": 8.5,
  "working_shifts": 1.0,
  "date": "2024-01-20T09:00:00Z"
}
```

**Response:** `201 Created`
```json
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

#### List All Track Items

**GET** `/api/track-items`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
[
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
]
```

#### List Track Items by Date Range

**GET** `/api/track-items?start_date=2024-01-01&end_date=2024-01-31`

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `start_date` (string, required): Start date in YYYY-MM-DD format
- `end_date` (string, required): End date in YYYY-MM-DD format

**Response:** `200 OK`
```json
[
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
]
```

#### Get a Specific Track Item

**GET** `/api/track-items/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
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

#### Update a Track Item

**PUT** `/api/track-items/:id`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:** (all fields optional)
```json
{
  "type": "overtime",
  "emergency_call": true,
  "working_hours": 10.0
}
```

**Response:** `200 OK`
```json
{
  "id": 1,
  "user_id": 1,
  "type": "overtime",
  "emergency_call": true,
  "holiday_call": false,
  "working_hours": 10.0,
  "working_shifts": 1.0,
  "date": "2024-01-20T09:00:00Z",
  "created_at": "2024-01-20T10:00:00Z",
  "updated_at": "2024-01-20T11:30:00Z"
}
```

#### Delete a Track Item

**DELETE** `/api/track-items/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:** `204 No Content`

---

## Data Models

### User

| Field | Type | Description |
|-------|------|-------------|
| `id` | integer | Unique user identifier |
| `first_name` | string | User's first name |
| `last_name` | string | User's last name |
| `avatar` | string | URL to user's avatar image (optional) |
| `login` | string | Unique login username |
| `created_at` | timestamp | Account creation time |
| `updated_at` | timestamp | Last update time |

### Track Item

| Field | Type | Description |
|-------|------|-------------|
| `id` | integer | Unique track item identifier |
| `user_id` | integer | ID of the user who owns this item |
| `type` | string | Type of work (e.g., "regular", "overtime", "remote") |
| `emergency_call` | boolean | Whether this was an emergency call |
| `holiday_call` | boolean | Whether this was a holiday call |
| `working_hours` | float | Number of hours worked |
| `working_shifts` | float | Number of shifts worked |
| `date` | timestamp | Date and time of the work (ISO 8601 format) |
| `created_at` | timestamp | Record creation time |
| `updated_at` | timestamp | Last update time |

---

## Error Responses

### 400 Bad Request
```json
{
  "error": "Invalid request body"
}
```

### 401 Unauthorized
```json
{
  "error": "Authorization header required"
}
```

### 403 Forbidden
```json
{
  "error": "Access denied"
}
```

### 404 Not Found
```json
{
  "error": "Track item not found"
}
```

### 409 Conflict
```json
{
  "error": "Login already exists"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error"
}
```

---

## Example Usage with curl

### Register and Login
```bash
# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "login": "johndoe",
    "password": "password123"
  }'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "johndoe",
    "password": "password123"
  }'

# Save the token
TOKEN="your-jwt-token-here"
```

### Create and Manage Track Items
```bash
# Create a track item
curl -X POST http://localhost:8080/api/track-items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "type": "regular",
    "emergency_call": false,
    "holiday_call": false,
    "working_hours": 8.0,
    "working_shifts": 1.0,
    "date": "2024-01-20T09:00:00Z"
  }'

# List all track items
curl http://localhost:8080/api/track-items \
  -H "Authorization: Bearer $TOKEN"

# Get track items for a date range
curl "http://localhost:8080/api/track-items?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer $TOKEN"

# Update a track item
curl -X PUT http://localhost:8080/api/track-items/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "working_hours": 9.5,
    "emergency_call": true
  }'

# Delete a track item
curl -X DELETE http://localhost:8080/api/track-items/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

## Vue 3 Integration Example

```javascript
// composables/useAuth.js
import { ref } from 'vue'

export function useAuth() {
  const token = ref(localStorage.getItem('token'))
  const user = ref(null)

  const register = async (firstName, lastName, login, password, avatar = '') => {
    const response = await fetch('http://localhost:8080/api/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        first_name: firstName,
        last_name: lastName,
        login,
        password,
        avatar
      })
    })
    const data = await response.json()
    token.value = data.token
    user.value = data.user
    localStorage.setItem('token', data.token)
  }

  const login = async (login, password) => {
    const response = await fetch('http://localhost:8080/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ login, password })
    })
    const data = await response.json()
    token.value = data.token
    user.value = data.user
    localStorage.setItem('token', data.token)
  }

  return { token, user, register, login }
}

// composables/useTrackItems.js
export function useTrackItems() {
  const trackItems = ref([])

  const fetchTrackItems = async (token, startDate = null, endDate = null) => {
    let url = 'http://localhost:8080/api/track-items'
    if (startDate && endDate) {
      url += `?start_date=${startDate}&end_date=${endDate}`
    }
    
    const response = await fetch(url, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    trackItems.value = await response.json()
  }

  const createTrackItem = async (token, itemData) => {
    const response = await fetch('http://localhost:8080/api/track-items', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(itemData)
    })
    return await response.json()
  }

  const updateTrackItem = async (token, id, updates) => {
    const response = await fetch(`http://localhost:8080/api/track-items/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(updates)
    })
    return await response.json()
  }

  const deleteTrackItem = async (token, id) => {
    await fetch(`http://localhost:8080/api/track-items/${id}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token}` }
    })
  }

  return {
    trackItems,
    fetchTrackItems,
    createTrackItem,
    updateTrackItem,
    deleteTrackItem
  }
}
```

---

## Database Schema

### users table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    avatar VARCHAR(500),
    login VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### track_items table
```sql
CREATE TABLE track_items (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(100) NOT NULL,
    emergency_call BOOLEAN NOT NULL DEFAULT FALSE,
    holiday_call BOOLEAN NOT NULL DEFAULT FALSE,
    working_hours DECIMAL(10, 2) NOT NULL DEFAULT 0,
    working_shifts DECIMAL(10, 2) NOT NULL DEFAULT 0,
    date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```
