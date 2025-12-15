# 🪟 Windows Setup Instructions

## Prerequisites

### 1. Install Go
- Download: https://golang.org/dl/
- Pilih `Windows` → ambil versi terbaru `.msi`
- Install dengan default settings
- Verify: Buka cmd, jalankan:
  ```cmd
  go version
  ```

### 2. Install Node.js
- Download: https://nodejs.org/
- Pilih LTS version
- Install dengan default settings
- Verify: Buka cmd, jalankan:
  ```cmd
  node --version
  npm --version
  ```

## 📦 Setup Backend

### Langkah 1: Navigate ke Backend
```cmd
cd d:\website_glb\backend
```

### Langkah 2: Download Dependencies
```cmd
go mod download
```

### Langkah 3: Jalankan Backend
```cmd
go run main.go
```

**Output yang diharapkan:**
```
🚀 Server running on http://localhost:8080
```

✅ **Backend siap!** Jangan tutup terminal ini.

---

## 🎨 Setup Frontend

### Langkah 1: Buka Terminal Baru
- Tekan `Ctrl + Shift + Esc` atau buka Command Prompt baru

### Langkah 2: Navigate ke Frontend
```cmd
cd d:\website_glb\frontend
```

### Langkah 3: Install Dependencies
```cmd
npm install
```

Proses ini memakan waktu 1-2 menit.

### Langkah 4: Jalankan Development Server
```cmd
npm run dev
```

**Output yang diharapkan:**
```
  ➜  Local:   http://localhost:5173/
  ➜  press h to show help
```

Browser akan otomatis membuka. Jika tidak, buka: `http://localhost:5173`

✅ **Frontend siap!**

---

## ✅ Verifikasi Setup

Pastikan kedua terminal menunjukkan:

**Terminal 1 (Backend):**
```
🚀 Server running on http://localhost:8080
```

**Terminal 2 (Frontend):**
```
➜  Local:   http://localhost:5173/
```

---

## 🧪 Test Flow

1. Buka browser ke: `http://localhost:5173`
2. Klik "Register"
3. Isi form:
   - Email: `test@example.com`
   - Password: `password123`
4. Klik Register
5. Login dengan credentials yang sama
6. Upload file GLB (test file)
7. View model di 3D viewer

---

## 🐛 Common Windows Issues

### Error: Port 8080 sudah terpakai
```cmd
# Matikan process yang pake port 8080
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Atau ganti port di main.go:
# Ganti: router.Run(":8080")
# Menjadi: router.Run(":9090")
```

### Error: npm command not found
- Restart cmd setelah install Node.js
- Atau tambahkan ke PATH manual

### Error: go command not found
- Restart cmd setelah install Go
- Atau tambahkan ke PATH manual

### File tidak bisa diupload
- Pastikan folder `backend/uploads` ada
- Windows Explorer: `d:\website_glb\backend\uploads`

---

## 🚀 Optimization Tips

### Develop Mode
```cmd
# Frontend auto-reload saat edit file
npm run dev

# Backend butuh restart manual saat edit Go
# Atau install: go get github.com/cosmtrek/air
# Jalankan: air
```

### Build untuk Production
```cmd
# Frontend
cd frontend
npm run build

# Output: frontend/dist/
```

---

## 📝 Notes untuk Windows Users

- Gunakan `\` atau `/` untuk path (keduanya work)
- PowerShell atau Command Prompt sama saja
- Jangan close terminal backend saat develop
- Untuk stop server: `Ctrl + C`

---

**Selamat! Setup Windows selesai! 🎉**
