# Quick Start Guide

## ⚡ Cara Jalankan Project

### Step 1: Backend

Buka terminal pertama:

```bash
# Navigate ke backend folder
cd backend

# Download dependencies Go
go mod download

# Run backend server
go run .
```

Tunggu sampai muncul pesan:
```
🚀 Server running on http://localhost:8080
```

### Step 2: Frontend

Buka terminal baru:

```bash
# Navigate ke frontend folder
cd frontend

# Install npm dependencies
npm install

# Run development server
npm run dev
```npm

Akan terbuka browser otomatis ke:
```
http://localhost:5173
```

## 🧑‍💻 Test Credentials

### Admin Account
```
Email: admin@test.com
Password: admin123
```

### Regular User
```
Email: user@test.com
Password: password123
```

> **Note:** Kedua akun harus didaftarkan terlebih dahulu atau di-seed di database. Untuk testing awal, silakan register akun baru.

## 📝 Workflow Testing

1. **Register** → Halaman login
2. **Login** sebagai admin → Admin dashboard muncul
3. **Upload GLB** → Pilih file `.glb`
4. **View Models** → Pergi ke viewer page
5. **Lihat 3D** → Rotate, zoom, pan model

## 🛑 Troubleshooting

| Error | Solusi |
|-------|--------|
| Port 8080 already in use | Change port di `main.go` |
| Port 5173 already in use | Vite akan auto-pick port |
| CORS error | Pastikan backend running dulu |
| Model tidak load | Check file URL di browser console |
| Login fail | Pastikan email/password benar |

## 📂 File Locations

```
Backend database: backend/glb.db
Uploaded models: backend/uploads/
```

---

✨ **Project siap dijalankan!**
