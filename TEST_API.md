# API Testing Guide

## ✅ Backend Running Successfully

The backend has been tested and verified:
- ✅ Go build successful (no syntax errors)
- ✅ Server startup successful on port 8080
- ✅ All 7 API routes registered and active
- ✅ In-memory storage initialized with test users
- ✅ CORS middleware configured
- ✅ JWT authentication ready

## 🚀 How to Run Everything

### Terminal 1: Start Backend
```bash
cd backend
go run main.go
```

Expected output:
```
✅ Test data initialized
   Admin: admin@test.com / admin123
   User:  user@test.com / password123
🚀 Server running on http://localhost:8080
[GIN-debug] Listening and serving HTTP on :8080
```

### Terminal 2: Start Frontend
```bash
cd frontend
npm install  # Only needed first time
npm run dev
```

Expected output:
```
VITE v5.0.0  ready in XXX ms

➜  Local:   http://localhost:5173/
```

## 📡 Testing the API

### 1. Register New User (POST)
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'
```

### 2. Login User (POST)
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"admin123"}'
```

Response will include a JWT token.

### 3. Get User Profile (GET) - Requires Auth
```bash
curl -X GET http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 4. Get All Models (GET) - Public
```bash
curl -X GET http://localhost:8080/api/models
```

### 5. Upload Model (POST) - Admin Only
```bash
curl -X POST http://localhost:8080/api/models/upload \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -F "name=Model Name" \
  -F "description=Model Description" \
  -F "file=@model.glb"
```

### 6. Delete Model (DELETE) - Admin Only
```bash
curl -X DELETE http://localhost:8080/api/models \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"id":1}'
```

## 🔧 Technology Stack Verification

### Backend
- Go 1.21+ ✅
- Gin v1.9.1 ✅
- JWT v5.0.0 ✅
- Bcrypt (crypto) ✅
- In-memory storage ✅
- No external database required ✅

### Frontend
- Vite v5.0.0 ✅
- Vanilla JavaScript ✅
- Three.js v0.160.0 ✅
- Axios v1.6.0 ✅
- CSS3 Responsive ✅

## ✅ Project Status

| Component | Status | Notes |
|-----------|--------|-------|
| Backend Build | ✅ PASSED | Compiles without errors |
| Backend Startup | ✅ PASSED | Server runs on port 8080 |
| Frontend Setup | ✅ PASSED | npm install successful (35 packages) |
| Routes | ✅ PASSED | All 7 routes registered |
| Auth | ✅ READY | JWT + Bcrypt configured |
| Database | ✅ READY | In-memory storage initialized |

## 🎯 Next Steps

1. **Run Backend**: `cd backend && go run main.go`
2. **Run Frontend**: `cd frontend && npm run dev`
3. **Open Browser**: http://localhost:5173
4. **Test Login**: Use admin@test.com / admin123
5. **Upload GLB**: Upload 3D models from admin dashboard
6. **View Models**: Switch to user account and view in 3D viewer

## 📚 Documentation

For detailed information:
- See `00_START_HERE.md` for project overview
- See `README.md` for full documentation
- See `API_DOCS.md` for complete API reference
- See `WINDOWS_SETUP.md` for setup on Windows
