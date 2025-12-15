# 📡 API Documentation

Base URL: `http://localhost:8080/api`

## Authentication Endpoints

### 1. Register User
**Endpoint:** `POST /auth/register`

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (201 Created):**
```json
{
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "role": "user"
  }
}
```

**Error (409 Conflict):**
```json
{
  "error": "Email already exists"
}
```

---

### 2. Login User
**Endpoint:** `POST /auth/login`

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "role": "user"
  },
  "role": "user"
}
```

**Error (401 Unauthorized):**
```json
{
  "error": "Invalid email or password"
}
```

---

### 3. Get User Profile
**Endpoint:** `GET /user/profile`

**Headers:**
```
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "message": "User profile",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "role": "user"
  }
}
```

---

## Model Endpoints

### 1. Get All Models
**Endpoint:** `GET /models`

**Response (200 OK):**
```json
{
  "message": "Models retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Sample Model",
      "description": "A 3D model",
      "file_url": "/uploads/1701234567_model.glb",
      "file_name": "1701234567_model.glb",
      "file_size": 5242880,
      "uploaded_by": "admin@test.com",
      "created_at": "2024-12-05 10:30:15"
    }
  ]
}
```

---

### 2. Upload GLB Model
**Endpoint:** `POST /models/upload`

**Headers:**
```
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**Form Data:**
| Field | Type | Required |
|-------|------|----------|
| file | File (.glb, .gltf) | Yes |
| name | String | Yes |
| description | String | No |

**Response (201 Created):**
```json
{
  "message": "Model uploaded successfully",
  "data": {
    "id": 2,
    "name": "New Model",
    "description": "My awesome model",
    "file_url": "/uploads/1701234567_model.glb",
    "file_name": "1701234567_model.glb",
    "file_size": 8388608
  }
}
```

**Error (403 Forbidden):**
```json
{
  "error": "Only admin can upload models"
}
```

**Error (400 Bad Request):**
```json
{
  "error": "Only .glb and .gltf files are allowed"
}
```

---

### 3. Delete Model
**Endpoint:** `DELETE /models/:id`

**Headers:**
```
Authorization: Bearer {token}
```

**Path Parameters:**
- `id` (Integer) - Model ID

**Response (200 OK):**
```json
{
  "message": "Model deleted successfully"
}
```

**Error (403 Forbidden):**
```json
{
  "error": "Only admin can delete models"
}
```

**Error (404 Not Found):**
```json
{
  "error": "Model not found"
}
```

---

## Static File Access

### Access Uploaded Models
**URL:** `GET /uploads/{filename}`

**Example:**
```
http://localhost:8080/uploads/1701234567_model.glb
```

Returns the binary GLB file (binary/octet-stream)

---

## Error Codes

| Code | Meaning |
|------|---------|
| 200 | OK - Success |
| 201 | Created - Resource created |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Missing/invalid token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource not found |
| 409 | Conflict - Resource already exists |
| 500 | Server Error |

---

## Authentication

All protected endpoints require JWT token in header:

```
Authorization: Bearer <token>
```

Token expires in 24 hours. After expiry, user must login again to get new token.

### Token Format
```
Header.Payload.Signature
```

**Payload contains:**
- `user_id` - User ID
- `email` - User email
- `role` - User role (admin/user)
- `exp` - Expiration time

---

## CORS Headers

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: POST, GET, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
```

---

## Testing with cURL

### Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### Get Models
```bash
curl http://localhost:8080/api/models
```

### Upload Model
```bash
curl -X POST http://localhost:8080/api/models/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@model.glb" \
  -F "name=My Model" \
  -F "description=Description"
```

---

## Rate Limiting

Currently no rate limiting. Production deployment should implement:
- IP-based rate limiting
- User-based throttling
- Request queuing

---

## Database Schema

### Users Table
```sql
CREATE TABLE users (
  id INTEGER PRIMARY KEY,
  email VARCHAR UNIQUE NOT NULL,
  password VARCHAR NOT NULL,
  role VARCHAR NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
)
```

### GLB Models Table
```sql
CREATE TABLE glb_models (
  id INTEGER PRIMARY KEY,
  name VARCHAR NOT NULL,
  description TEXT,
  file_url VARCHAR NOT NULL,
  file_name VARCHAR NOT NULL,
  uploaded_by INTEGER,
  file_size INTEGER,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (uploaded_by) REFERENCES users(id)
)
```

---

**Dokumentasi API selesai!**
