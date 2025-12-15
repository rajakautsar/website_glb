# 🎨 GLB 3D Model Management System

Aplikasi web lengkap untuk mengelola dan menampilkan model 3D dalam format GLB/GLTF dengan backend Go dan frontend Vite.

## 📋 Struktur Folder

```
website_glb/
├── backend/
│   ├── controllers/
│   │   ├── auth.go              # Auth handlers
│   │   └── models.go            # Model upload/list handlers
│   ├── middleware/
│   │   └── auth.go              # JWT middleware
│   ├── models/
│   │   └── models.go            # Data models
│   ├── utils/
│   │   └── jwt.go               # JWT utilities
│   ├── uploads/                 # File storage (auto-created)
│   ├── main.go                  # Entry point
│   ├── go.mod                   # Dependencies
│   └── glb.db                   # SQLite database (auto-created)
│
└── frontend/
    ├── src/
    │   ├── api.js               # API client
    │   └── style.css            # Global styles
    ├── index.html               # Login page
    ├── admin.html               # Admin dashboard
    ├── viewer.html              # 3D viewer
    ├── package.json             # Dependencies
    ├── vite.config.js           # Vite config
    └── public/                  # Static files
```

## 🚀 Persiapan & Instalasi

### Backend Setup

1. **Install Go** (versi 1.21+)
   - Download: https://golang.org/dl/

2. **Navigate ke backend folder:**
   ```bash
   cd backend
   ```

3. **Download dependencies:**
   ```bash
   go mod download
   ```

4. **Run backend:**
   ```bash
   go run main.go
   ```
   
   Backend akan berjalan di: `http://localhost:8080`

### Frontend Setup

1. **Install Node.js** (versi 16+)
   - Download: https://nodejs.org/

2. **Navigate ke frontend folder:**
   ```bash
   cd frontend
   ```

3. **Install dependencies:**
   ```bash
   npm install
   ```

4. **Run development server:**
   ```bash
   npm run dev
   ```
   
   Frontend akan berjalan di: `http://localhost:5173`

## 🔑 Fitur Utama

### Backend (Go + Gin + SQLite)

#### Authentication
- **POST** `/api/auth/register` - Register user baru
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```
- **POST** `/api/auth/login` - Login & dapatkan JWT token
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```
- **GET** `/api/user/profile` - Dapatkan profile user (protected)

#### Model Management
- **GET** `/api/models` - Dapatkan daftar semua model 3D (public)
- **POST** `/api/models/upload` - Upload file GLB (admin only)
  - Form-data: `file`, `name`, `description`
- **DELETE** `/api/models/:id` - Hapus model (admin only)
- **Static** `/uploads` - Akses file GLB yang sudah diupload

### Frontend (Vite + Vanilla JS + Three.js)

#### Pages

**1. Login Page** (`index.html`)
- Register user baru
- Login dengan email & password
- JWT token disimpan di localStorage
- Auto-redirect ke dashboard sesuai role

**2. Admin Dashboard** (`admin.html`)
- Upload file GLB/GLTF
- Lihat daftar semua model
- Delete model
- View/preview model

**3. 3D Viewer** (`viewer.html`)
- Menampilkan model 3D menggunakan Three.js
- Sidebar daftar model
- Click model untuk melihat detail
- Orbit controls untuk manipulasi 3D
- Lighting & shadows
- Info panel

## 📝 User Roles

### Admin
- Bisa upload file GLB
- Bisa hapus model
- Akses admin dashboard
- Bisa melihat semua model

### User (Regular)
- Tidak bisa upload
- Hanya bisa view model 3D
- Redirect ke viewer page

## 🔐 Keamanan

- JWT token dengan expiry 24 jam
- Password di-hash menggunakan bcrypt
- Admin-only middleware untuk upload/delete
- CORS enabled untuk frontend
- Authorization header validation

## 📦 Dependencies

### Backend
```
github.com/gin-gonic/gin v1.9.1           # Web framework
github.com/golang-jwt/jwt/v5 v5.0.0       # JWT
github.com/mattn/go-sqlite3 v1.14.17      # SQLite driver
golang.org/x/crypto v0.15.0                # Password hashing
gorm.io/driver/sqlite v1.5.2              # GORM SQLite
gorm.io/gorm v1.25.4                      # ORM
```

### Frontend
```
three: ^r128                                # 3D graphics
vite: ^5.0.0                               # Build tool
axios: ^1.6.0                              # HTTP client
```

## 🧪 Testing

### Buat Admin User (Testing)

Perlu edit `main.go` untuk seed admin user pertama kali:

```go
// Add di function main() sebelum router.Run():
db.Create(&User{
    Email: "admin@test.com",
    Password: bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost),
    Role: "admin",
})
```

### Test Flow

1. **Register User:**
   ```
   Email: user@test.com
   Password: password123
   ```

2. **Login Admin:**
   ```
   Email: admin@test.com
   Password: admin123
   ```

3. **Upload Model:**
   - Go to admin dashboard
   - Upload GLB file
   - Model akan tersimpan di `/uploads`

4. **View Model:**
   - Go to viewer page
   - Pilih model dari sidebar
   - Lihat 3D model dengan Three.js

## 🔧 Troubleshooting

### CORS Error
Pastikan backend di `http://localhost:8080` dan frontend di `http://localhost:5173`

### Token Expired
JWT token berlaku 24 jam. User harus login ulang setelah token expired.

### File Not Uploading
- Pastikan folder `uploads/` ada (auto-created)
- Pastikan file extension `.glb` atau `.gltf`
- Pastikan user adalah admin

### Model Not Showing
- Pastikan file URL benar: `/uploads/filename.glb`
- Check browser console untuk error
- Pastikan model tidak corrupt

## 📱 Browser Support

- Chrome/Edge: ✅ Full support
- Firefox: ✅ Full support
- Safari: ✅ Full support (iOS 13+)
- Mobile: ⚠️ Responsive (sidebar collapse)

## 🎨 Customization

### Change JWT Secret
Edit `backend/utils/jwt.go`:
```go
const JWTSecret = "your-new-secret-key"
```

### Change API URL
Edit `frontend/src/api.js`:
```javascript
const API_URL = 'http://your-backend-url/api';
```

### Change Port
Backend: Edit `main.go` → `router.Run(":YOUR_PORT")`
Frontend: Edit `vite.config.js` → `port: YOUR_PORT`

## 🚢 Deployment

### Backend (Linux/Windows)
1. Build: `go build -o glb-server`
2. Run: `./glb-server`

### Frontend
1. Build: `npm run build`
2. Output di `dist/` folder
3. Deploy ke static hosting (Vercel, Netlify, etc.)

## 📞 Support

Jika ada error atau pertanyaan, check:
- Browser console untuk error messages
- Backend logs di terminal
- CORS headers di network tab
- JWT token validity

---

**Project dibuat dengan ❤️ untuk GLB 3D Model Management**
