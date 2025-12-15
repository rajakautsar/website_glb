# ✅ FINAL INSTALLATION CHECKLIST

## Pre-Installation Requirements

- [ ] Windows/Mac/Linux OS
- [ ] Go 1.21 atau lebih tinggi
- [ ] Node.js 16 atau lebih tinggi
- [ ] 200MB free disk space
- [ ] Text editor (VS Code recommended)

## Installation Steps

### Step 1: Backend Setup (Terminal 1)
```bash
cd backend
go mod download
```
- [ ] No errors during download
- [ ] Go modules downloaded

### Step 2: Run Backend
```bash
cd backend
go run main.go
```
- [ ] Server starts successfully
- [ ] Output: "🚀 Server running on http://localhost:8080"
- [ ] No CORS errors
- [ ] Port 8080 is free

### Step 3: Frontend Setup (Terminal 2)
```bash
cd frontend
npm install
```
- [ ] npm install completes
- [ ] No build errors
- [ ] No security warnings (or only low severity)

### Step 4: Run Frontend
```bash
cd frontend
npm run dev
```
- [ ] Dev server starts
- [ ] Output: "http://localhost:5173"
- [ ] Browser opens automatically

## Verification Checklist

### Backend Running
- [ ] Terminal shows: "🚀 Server running on http://localhost:8080"
- [ ] No error messages
- [ ] Database file created (glb.db)
- [ ] Uploads folder created

### Frontend Running
- [ ] Terminal shows dev server URL
- [ ] Browser shows login page
- [ ] No 404 errors in console
- [ ] CSS loads correctly

### API Connection
- [ ] Can see login form
- [ ] Register button works
- [ ] No CORS errors in console

## First-Time Setup

### Seed Database (Optional)
```bash
cd backend
go run main.go seed.go
```
- [ ] Test users created
- [ ] See "✅ Database seeding complete!"

Or manually:
1. Register as: test@example.com / password123
2. Admin: admin@test.com / admin123

### Test Login
- [ ] Login successful
- [ ] Token stored in localStorage
- [ ] Redirected to correct page (admin/viewer)

### Test File Upload (Admin Only)
- [ ] Go to admin dashboard
- [ ] Upload test GLB file
- [ ] File appears in list
- [ ] File accessible from /uploads

### Test 3D Viewer
- [ ] Go to viewer page
- [ ] Model list appears
- [ ] Click model loads 3D
- [ ] Can rotate model (orbit controls)

## Performance Checks

### Backend Performance
- [ ] Response time < 200ms
- [ ] No memory leaks
- [ ] Handles multiple requests
- [ ] File upload works for large files (100MB+)

### Frontend Performance
- [ ] Page loads in < 2 seconds
- [ ] 3D model renders smoothly
- [ ] No lag during orbit controls
- [ ] Responsive on mobile

## Security Checks

### Authentication
- [ ] JWT token generated on login
- [ ] Token expires after 24 hours
- [ ] Invalid credentials rejected
- [ ] Protected endpoints require token

### Authorization
- [ ] Regular users can't upload
- [ ] Regular users can't delete
- [ ] Admin can upload/delete
- [ ] Role validation works

### File Upload
- [ ] Only .glb/.gltf allowed
- [ ] File validation works
- [ ] Files stored securely
- [ ] File paths exposed correctly

## Browser Compatibility

- [ ] Chrome/Edge latest
- [ ] Firefox latest
- [ ] Safari 13+
- [ ] Mobile browsers

## Common Issues Resolved

### Port Already In Use
- [ ] Found alternative port
- [ ] Updated configuration
- [ ] Server running successfully

### Missing Dependencies
- [ ] go mod download completed
- [ ] npm install completed
- [ ] No "not found" errors

### CORS Issues
- [ ] Backend returns correct headers
- [ ] Frontend can call API
- [ ] No cross-origin errors

### Authentication Failed
- [ ] User credentials correct
- [ ] Database contains user
- [ ] Token generation works

## Documentation Reviewed

- [ ] README.md read ✅
- [ ] API_DOCS.md reviewed ✅
- [ ] QUICKSTART.md understood ✅
- [ ] WINDOWS_SETUP.md followed (if Windows) ✅

## Deployment Ready

- [ ] Code complete
- [ ] All features working
- [ ] Database initialized
- [ ] CORS configured
- [ ] Error handling implemented
- [ ] Logging implemented

## Final Status

**Overall Setup:** ✅ COMPLETE

**Backend Status:** ✅ RUNNING
**Frontend Status:** ✅ RUNNING

**Ready for Development:** ✅ YES

---

## Next Steps

1. ✅ Start developing features
2. ✅ Add more 3D models
3. ✅ Customize styling
4. ✅ Deploy to production

---

**Checklist Last Updated:** December 5, 2024
**Project Version:** 1.0
