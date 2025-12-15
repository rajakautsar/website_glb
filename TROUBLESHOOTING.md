# 🔧 TROUBLESHOOTING GUIDE

## Quick Diagnostics

Run this checklist if something doesn't work:

```
✓ Is Go installed?          go version
✓ Is Node.js installed?     node -v && npm -v
✓ Is backend running?       http://localhost:8080
✓ Is frontend running?      http://localhost:5173
✓ Database exists?          backend/glb.db
✓ Uploads folder exists?    backend/uploads/
```

---

## Common Issues & Solutions

### ❌ "Port 8080 already in use"

**Symptom:** Backend won't start, error about port 8080

**Solutions:**
1. Kill the existing process:
   ```bash
   # Windows
   netstat -ano | findstr :8080
   taskkill /PID <PID> /F
   
   # Mac/Linux
   lsof -i :8080
   kill -9 <PID>
   ```

2. Use different port:
   - Edit `backend/main.go`
   - Change: `router.Run(":8080")`
   - To: `router.Run(":9090")`

3. Check what's using it:
   - Windows: Task Manager → Networking tab
   - Mac: Activity Monitor
   - Linux: `sudo netstat -tlnp`

---

### ❌ "Port 5173 already in use"

**Symptom:** Frontend shows different port in terminal

**Solution:** This is normal! Vite auto-selects next available port
- Check terminal output for actual port
- Open that URL in browser

---

### ❌ "go: command not found"

**Symptom:** Terminal says "go: command not found"

**Causes & Solutions:**
1. Go not installed:
   ```bash
   go version  # Check if installed
   # If not: https://golang.org/dl/
   ```

2. PATH not updated:
   - Windows: Restart cmd/PowerShell after installing Go
   - Mac: `export PATH=$PATH:/usr/local/go/bin`
   - Linux: Usually automatic

3. Wrong installation path:
   - Windows: Reinstall Go
   - Mac/Linux: Check `/usr/local/go/bin/go`

---

### ❌ "npm: command not found"

**Symptom:** Terminal says "npm: command not found"

**Solutions:**
1. Node.js not installed:
   ```bash
   node -v
   npm -v
   # If not: https://nodejs.org/
   ```

2. PATH issue:
   - Windows: Restart after installing Node.js
   - Mac/Linux: Check `which node`

3. Multiple Node versions:
   ```bash
   which node
   which npm
   node -v
   ```

---

### ❌ "CORS error" in browser console

**Symptom:** Error: "Access to XMLHttpRequest... blocked by CORS policy"

**Causes:**
1. Backend not running
2. Wrong API URL in frontend/src/api.js
3. Backend CORS not configured

**Solutions:**
1. Check backend is running:
   ```bash
   curl http://localhost:8080/api/models
   # Should return JSON, not error
   ```

2. Check CORS headers:
   ```bash
   curl -i http://localhost:8080/api/models
   # Look for: Access-Control-Allow-Origin: *
   ```

3. Frontend API URL:
   - Open `frontend/src/api.js`
   - Verify: `const API_URL = 'http://localhost:8080/api'`
   - Should match your backend port

4. Restart both services

---

### ❌ "Database locked" or "database is corrupted"

**Symptom:** SQLite error when running backend

**Solutions:**
1. Delete old database:
   ```bash
   cd backend
   rm glb.db
   # Restart backend to create new one
   ```

2. Check file permissions:
   ```bash
   # Windows: Right-click → Properties → Security
   # Mac/Linux: chmod 644 glb.db
   ```

3. Close other connections:
   - Make sure no other process using glb.db
   - Check: `lsof backend/glb.db` (Mac/Linux)

---

### ❌ "Login fails" or "Invalid credentials"

**Symptom:** Can't login, even with correct credentials

**Causes & Solutions:**
1. User doesn't exist:
   - Seed database: `go run main.go seed.go`
   - Or register new user

2. Password wrong:
   - Check credentials: admin@test.com / admin123
   - Passwords are case-sensitive

3. Backend not running:
   - Check: `curl http://localhost:8080/api/models`
   - Should work without login

4. Database corrupted:
   - Delete `backend/glb.db`
   - Restart backend
   - Re-seed data

---

### ❌ "Model not uploading"

**Symptom:** Upload button doesn't work, or "Only admin can upload"

**Causes & Solutions:**
1. Not logged in as admin:
   - Login as: admin@test.com
   - Check localStorage in DevTools

2. Not admin role:
   - Check user role in database
   - Delete user and re-register

3. File type wrong:
   - Only .glb or .gltf allowed
   - File size limit: 100MB (in config)

4. Uploads folder missing:
   ```bash
   cd backend
   mkdir uploads
   # Restart backend
   ```

5. Permissions issue:
   ```bash
   # Mac/Linux
   chmod 755 uploads/
   chmod 644 uploads/*
   ```

---

### ❌ "Model not showing in 3D viewer"

**Symptom:** Blank screen or "Loading model..." forever

**Causes & Solutions:**
1. File path incorrect:
   - Check network tab in DevTools
   - File URL should be: `/uploads/1234567890_filename.glb`
   - Should respond with 200 status

2. GLB file corrupted:
   - Try uploading different file
   - Check file size: not 0 bytes

3. GLTFLoader not working:
   ```javascript
   // Check browser console for errors
   // Three.js should load from CDN
   ```

4. Model not rendered:
   - Check DevTools → Console for errors
   - Rotate with mouse to see if rendering
   - Zoom in/out to check scale

---

### ❌ "npm install takes forever"

**Symptom:** `npm install` stuck or very slow

**Solutions:**
1. Kill and retry:
   ```bash
   npm install --force
   # or
   npm ci  # Cleaner install
   ```

2. Clear cache:
   ```bash
   npm cache clean --force
   npm install
   ```

3. Use different registry:
   ```bash
   npm install -r https://registry.npmjs.org/
   ```

4. Network issue:
   - Check internet connection
   - Try different WiFi/network

---

### ❌ "go mod download fails"

**Symptom:** Error downloading Go dependencies

**Solutions:**
1. Retry:
   ```bash
   go mod download
   go mod tidy
   ```

2. Clear cache:
   ```bash
   go clean -modcache
   go mod download
   ```

3. Network issue:
   - Check internet
   - Try: `go env -w GOPROXY=direct`

4. Proxy issue:
   ```bash
   # Reset to default
   go env -w GOPROXY=https://proxy.golang.org,direct
   ```

---

### ❌ "404 Not Found" on pages

**Symptom:** Browser shows 404 when visiting pages

**Causes:**
1. Frontend not running
   - Check terminal for: `http://localhost:5173`
   - Make sure `npm run dev` is running

2. Wrong URL:
   - Should be: http://localhost:5173
   - Not: http://localhost:8080

3. Vite proxy issue:
   - Edit `vite.config.js`
   - Clear browser cache

---

### ❌ "Token expired" or "Unauthorized"

**Symptom:** Suddenly logged out or "401 Unauthorized"

**Causes & Solutions:**
1. Token expired (24 hours):
   - Login again
   - Token will be refreshed

2. Token not saved:
   - Check localStorage in DevTools
   - Key should be: `token`

3. Token invalid:
   - Clear localStorage: 
     ```javascript
     localStorage.clear()
     ```
   - Login again

4. Backend JWT secret changed:
   - All tokens become invalid
   - All users must login again

---

### ❌ Browser console shows errors

**Symptom:** DevTools → Console has red error messages

**Common Errors:**

1. **Uncaught TypeError: API is undefined**
   - Missing `import` statement
   - Check: `import { getModels } from '/src/api.js'`

2. **Failed to fetch resource: the server responded with a status of 404**
   - Backend endpoint not found
   - Check API URL
   - Check backend is running

3. **Uncaught ReferenceError: THREE is not defined**
   - Three.js not loaded
   - Check CDN link in viewer.html
   - Check internet connection

4. **CORS error**
   - See CORS section above

---

### ❌ Slow performance

**Symptom:** App is very slow or unresponsive

**Solutions:**
1. Large 3D model:
   - File size >50MB
   - Try smaller model first
   - Check browser memory in DevTools

2. Many models in list:
   - Backend: Add pagination
   - Frontend: Lazy load models

3. Browser issue:
   - Clear cache: Ctrl+Shift+Del
   - Restart browser
   - Try different browser

4. System resources:
   - Close other apps
   - Check CPU/Memory usage
   - Restart computer

---

### ❌ File permissions error

**Symptom:** "Permission denied" when saving files

**Causes & Solutions:**

1. Windows:
   - Right-click folder → Properties → Security
   - Give full control to your user

2. Mac/Linux:
   ```bash
   sudo chmod -R 755 backend/uploads
   sudo chown -R $USER:$USER backend/
   ```

3. Docker (if using):
   ```dockerfile
   RUN chown -R app:app /app
   ```

---

## Advanced Debugging

### Enable Verbose Logging

**Backend:**
```go
// In main.go, add:
gin.SetMode(gin.DebugMode)
```

**Frontend:**
```javascript
// In api.js, add logging:
console.log('Calling:', url, method);
console.log('Response:', response);
```

### Check Network Traffic

1. Browser DevTools:
   - F12 → Network tab
   - Click on request
   - Check: Headers, Preview, Response

2. Command line:
   ```bash
   # Monitor all requests
   tcpdump -i any -n -s0 -w traffic.pcap host localhost
   ```

### Database Debug

```bash
# Inspect SQLite database
sqlite3 backend/glb.db

# Show tables
.tables

# Query data
SELECT * FROM users;
SELECT * FROM glb_models;

# Exit
.quit
```

### Check Ports

```bash
# Windows
netstat -ano | findstr :8080
netstat -ano | findstr :5173

# Mac/Linux
lsof -i :8080
lsof -i :5173
```

---

## When All Else Fails

### Complete Reset

```bash
# Stop all services (Ctrl+C)

# Backend reset
cd backend
rm glb.db
go mod download
go run main.go

# Frontend reset (new terminal)
cd frontend
rm -rf node_modules package-lock.json
npm install
npm run dev

# Clear browser cache
# DevTools → Network → Disable cache (checkbox)
```

### Check System

```bash
# Check all requirements
echo "Go:" && go version
echo "Node:" && node -v
echo "npm:" && npm -v

# Check disk space
# Windows: C: drive properties
# Mac: About This Mac → Storage
# Linux: df -h
```

### Get Help

1. Check error message again
2. Read relevant documentation
3. Search GitHub issues
4. Check browser console carefully
5. Review this guide again

---

## File Locations Reference

```
Configuration:
  backend/main.go           (Backend config)
  frontend/src/api.js       (API config)
  vite.config.js           (Frontend config)

Data:
  backend/glb.db           (SQLite database)
  backend/uploads/         (Uploaded files)

Logs:
  Terminal output          (Backend logs)
  Browser console          (Frontend logs)

LocalStorage:
  DevTools → Application   (Stored data)
```

---

## Prevention Tips

1. Keep backups of database
2. Test uploads before deployment
3. Use version control (git)
4. Keep dependencies updated
5. Monitor error logs
6. Test on multiple browsers
7. Document customizations

---

## Getting Support

If still stuck:

1. ✅ Read documentation again
2. ✅ Check all error messages
3. ✅ Review this guide
4. ✅ Check browser console
5. ✅ Check terminal output
6. ✅ Try fresh install
7. ✅ Check system requirements

**Remember:** Most issues are either:
- Port already in use
- Dependencies not installed
- Backend not running
- CORS misconfiguration
- Browser cache issues

---

**Last Updated:** December 5, 2024
**Project Version:** 1.0
