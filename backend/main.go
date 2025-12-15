package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// ============ IN-MEMORY STORAGE ============
var (
	users                 = make(map[uint]*User)
	models                = make(map[uint]*GLBModel)
	archives              = make(map[uint]*Archive)
	userIDCounter    uint = 1
	modelIDCounter   uint = 1
	archiveIDCounter uint = 1
	mu               sync.RWMutex
)

// ============ MODELS ============
type User struct {
	ID           uint      `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type GLBModel struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	FileURL     string    `json:"file_url"`
	FileName    string    `json:"file_name"`
	ArchiveID   uint      `json:"archive_id"`
	UploadedBy  uint      `json:"uploaded_by"`
	FileSize    int64     `json:"file_size"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Archive struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"` // folder name (e.g., ARSIP_001)
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

// ============ REQUEST/RESPONSE STRUCTS ============
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
	Role  string `json:"role"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ============ JWT HELPERS ============
const JWTSecret = "your-super-secret-key-change-in-production"

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func generateToken(userID uint, email string, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

func verifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// generateRandomToken returns a hex encoded random string of length 2*n
func generateRandomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// ============ MIDDLEWARE ============
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := verifyToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// ============ AUTH HANDLERS ============
func registerHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error processing password"})
		return
	}

	mu.Lock()
	for _, u := range users {
		if u.Email == req.Email {
			mu.Unlock()
			c.JSON(409, gin.H{"error": "Email already exists"})
			return
		}
	}

	user := &User{
		ID:           userIDCounter,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	users[userIDCounter] = user
	userIDCounter++
	mu.Unlock()

	c.JSON(201, gin.H{
		"message": "User registered successfully",
		"data": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func loginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	mu.RLock()
	var user *User
	for _, u := range users {
		if u.Email == req.Email {
			user = u
			break
		}
	}
	mu.RUnlock()

	if user == nil {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := generateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(200, AuthResponse{
		Token: token,
		User: User{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
		Role: user.Role,
	})
}

func getUserProfileHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	mu.RLock()
	user, ok := users[userID.(uint)]
	mu.RUnlock()

	if !ok {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, gin.H{
		"message": "User profile",
		"data": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// ============ MODEL HANDLERS ============
func uploadModelHandler(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		c.JSON(403, gin.H{"error": "Only admin can upload models"})
		return
	}

	userID, _ := c.Get("user_id")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}

	fileExt := filepath.Ext(file.Filename)
	if fileExt != ".glb" && fileExt != ".gltf" {
		c.JSON(400, gin.H{"error": "Only .glb and .gltf files are allowed"})
		return
	}

	name := c.PostForm("name")
	description := c.PostForm("description")
	archiveIDStr := c.PostForm("archive_id")

	if name == "" {
		c.JSON(400, gin.H{"error": "Model name is required"})
		return
	}

	// determine destination: default uploads/ unless archive specified
	var destDir string = "uploads"
	var archiveID uint = 0
	var fileURL string
	if archiveIDStr != "" {
		// parse archive id
		var aid uint64
		var err error
		if aid, err = strconv.ParseUint(archiveIDStr, 10, 64); err == nil {
			mu.RLock()
			arch, ok := archives[uint(aid)]
			mu.RUnlock()
			if ok {
				destDir = filepath.Join("model_archives", arch.Name)
				archiveID = uint(aid)
			} else {
				c.JSON(400, gin.H{"error": "Archive not found"})
				return
			}
		} else {
			c.JSON(400, gin.H{"error": "Invalid archive_id"})
			return
		}
	}

	// ensure dest dir exists
	if err := os.MkdirAll(destDir, 0755); err != nil {
		c.JSON(500, gin.H{"error": "Error creating destination directory"})
		return
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filePath := filepath.Join(destDir, fileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(500, gin.H{"error": "Error saving file"})
		return
	}

	if archiveID != 0 {
		// File served via secure archive route
		mu.RLock()
		arch := archives[archiveID]
		mu.RUnlock()
		fileURL = fmt.Sprintf("/api/archives/%s/files/%s", arch.Name, fileName)
	} else {
		fileURL = fmt.Sprintf("/uploads/%s", fileName)
	}

	// insert into SQLite (if available)
	var dbID int64 = 0
	if DB != nil {
		var aid *int64
		var ub *int64
		if archiveID != 0 {
			a := int64(archiveID)
			aid = &a
		}
		if userID != nil {
			if uid, ok := userID.(uint); ok {
				u := int64(uid)
				ub = &u
			}
		}
		id, err := InsertModel(name, description, fileName, fileURL, int64(file.Size), aid, ub)
		if err != nil {
			log.Printf("Warning: failed insert model to sqlite: %v", err)
		} else {
			dbID = id
		}
	}

	// also keep in-memory for compatibility
	mu.Lock()
	assignedID := modelIDCounter
	if dbID > 0 {
		// keep counters in sync
		if uint(dbID) >= modelIDCounter {
			assignedID = uint(dbID)
			modelIDCounter = uint(dbID + 1)
		}
	}
	model := &GLBModel{
		ID:          assignedID,
		Name:        name,
		Description: description,
		FileName:    fileName,
		FileURL:     fileURL,
		ArchiveID:   archiveID,
		UploadedBy: func() uint {
			if userID != nil {
				if uid, ok := userID.(uint); ok {
					return uid
				}
			}
			return 0
		}(),
		FileSize:  file.Size,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	models[assignedID] = model
	if assignedID >= modelIDCounter {
		modelIDCounter = assignedID + 1
	}
	mu.Unlock()

	c.JSON(201, gin.H{
		"message": "Model uploaded successfully",
		"data": gin.H{
			"id":          model.ID,
			"name":        model.Name,
			"description": model.Description,
			"file_url":    model.FileURL,
			"file_name":   model.FileName,
			"file_size":   model.FileSize,
			"archive_id":  model.ArchiveID,
		},
	})
}

func getModelsHandler(c *gin.Context) {
	// Allow optional filtering by archive. If caller is an archive_user, only return their archive.
	var archiveFilter uint = 0
	// check Authorization header for archive user token
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			if claims, err := verifyToken(parts[1]); err == nil {
				if claims.Role == "archive_user" {
					archiveFilter = claims.UserID
				}
			}
		}
	}

	// query param override (admin can request specific archive)
	if archiveFilter == 0 {
		if aid := c.Query("archive_id"); aid != "" {
			if v, err := strconv.ParseUint(aid, 10, 64); err == nil {
				archiveFilter = uint(v)
			}
		}
	}

	mu.RLock()
	defer mu.RUnlock()

	var response []interface{}
	for _, model := range models {
		if archiveFilter != 0 && model.ArchiveID != archiveFilter {
			continue
		}
		var uploaderEmail string
		if user, ok := users[model.UploadedBy]; ok {
			uploaderEmail = user.Email
		}

		response = append(response, gin.H{
			"id":          model.ID,
			"name":        model.Name,
			"description": model.Description,
			"file_url":    model.FileURL,
			"file_name":   model.FileName,
			"file_size":   model.FileSize,
			"uploaded_by": uploaderEmail,
			"archive_id":  model.ArchiveID,
			"created_at":  model.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(200, gin.H{
		"message": "Models retrieved successfully",
		"data":    response,
	})
}

// ============ ARCHIVE HANDLERS & MIDDLEWARE ============
func createArchiveHandler(c *gin.Context) {
	// only admin
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		c.JSON(403, ErrorResponse{Error: "Only admins can create archives"})
		return
	}

	// optional custom name
	name := c.PostForm("name")
	if name == "" {
		name = fmt.Sprintf("ARSIP_%d", time.Now().Unix())
	}
	// sanitize name: allow alphanum and underscore/dash
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "_")

	token, err := generateRandomToken(16)
	if err != nil {
		c.JSON(500, ErrorResponse{Error: "Failed to generate token"})
		return
	}

	// create folder
	dir := filepath.Join("model_archives", name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.JSON(500, ErrorResponse{Error: "Failed to create archive folder"})
		return
	}

	// write token.txt
	if err := os.WriteFile(filepath.Join(dir, "token.txt"), []byte(token), 0644); err != nil {
		c.JSON(500, ErrorResponse{Error: "Failed to write token file"})
		return
	}

	mu.Lock()
	arch := &Archive{
		ID:        archiveIDCounter,
		Name:      name,
		Token:     token,
		CreatedAt: time.Now(),
	}
	archives[archiveIDCounter] = arch
	archiveIDCounter++
	mu.Unlock()

	c.JSON(201, gin.H{"message": "Archive created", "data": arch})
}

func listArchivesHandler(c *gin.Context) {
	// only admin
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		c.JSON(403, ErrorResponse{Error: "Only admins can list archives"})
		return
	}

	mu.RLock()
	defer mu.RUnlock()

	var resp []interface{}
	for _, a := range archives {
		// count models in archive
		count := 0
		for _, m := range models {
			if m.ArchiveID == a.ID {
				count++
			}
		}
		resp = append(resp, gin.H{
			"id":         a.ID,
			"name":       a.Name,
			"token":      a.Token,
			"count":      count,
			"created_at": a.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(200, gin.H{"message": "Archives retrieved", "data": resp})
}

func deleteArchiveHandler(c *gin.Context) {
	// only admin
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		c.JSON(403, ErrorResponse{Error: "Only admins can delete archives"})
		return
	}

	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid request"})
		return
	}

	mu.Lock()
	arch, ok := archives[req.ID]
	if !ok {
		mu.Unlock()
		c.JSON(404, ErrorResponse{Error: "Archive not found"})
		return
	}

	// delete models in archive
	for id, m := range models {
		if m.ArchiveID == req.ID {
			// remove file
			p := filepath.Join("model_archives", arch.Name, m.FileName)
			if err := os.Remove(p); err != nil {
				log.Printf("Warning: failed to remove file %s: %v", p, err)
			}
			delete(models, id)
		}
	}

	// delete folder
	dir := filepath.Join("model_archives", arch.Name)
	if err := os.RemoveAll(dir); err != nil {
		log.Printf("Warning: failed to remove archive folder %s: %v", dir, err)
	}

	delete(archives, req.ID)
	mu.Unlock()

	c.JSON(200, gin.H{"message": "Archive deleted"})
}

func archiveLoginHandler(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid request"})
		return
	}

	mu.RLock()
	var found *Archive
	for _, a := range archives {
		if a.Token == req.Token {
			found = a
			break
		}
	}
	mu.RUnlock()

	if found == nil {
		c.JSON(401, ErrorResponse{Error: "Invalid token"})
		return
	}

	// create JWT for archive user; reuse Claims.UserID to store archive ID
	tokenStr, err := generateToken(found.ID, found.Name, "archive_user")
	if err != nil {
		c.JSON(500, ErrorResponse{Error: "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{"message": "Login successful", "token": tokenStr, "archive": found})
}

func archiveAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}
		token := parts[1]
		claims, err := verifyToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		if claims.Role != "archive_user" {
			c.JSON(403, gin.H{"error": "Not an archive token"})
			c.Abort()
			return
		}
		// set archive info in context
		c.Set("archive_id", claims.UserID)
		c.Set("archive_name", claims.Email)
		c.Next()
	}
}

func archiveFileHandler(c *gin.Context) {
	archiveName := c.Param("archiveName")
	fileName := c.Param("fileName")

	// ensure requester is archive user and matches archiveName
	aidInterface, ok := c.Get("archive_id")
	if !ok {
		c.JSON(401, ErrorResponse{Error: "Unauthorized"})
		return
	}
	aid := aidInterface.(uint)

	mu.RLock()
	arch, exists := archives[aid]
	mu.RUnlock()
	if !exists || arch.Name != archiveName {
		c.JSON(403, ErrorResponse{Error: "Forbidden"})
		return
	}

	// safe file join
	cleanName := filepath.Clean(fileName)
	if strings.Contains(cleanName, "..") {
		c.JSON(400, ErrorResponse{Error: "Invalid file name"})
		return
	}

	p := filepath.Join("model_archives", archiveName, cleanName)
	// ensure file exists
	if _, err := os.Stat(p); os.IsNotExist(err) {
		c.JSON(404, ErrorResponse{Error: "File not found"})
		return
	}

	c.File(p)
}

func deleteModelHandler(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, ErrorResponse{Error: "Unauthorized"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(401, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	log.Printf("deleteModelHandler: invoked by user interface present: %v", exists)
	// Check if user is admin
	mu.RLock()
	user, userExists := users[userID]
	mu.RUnlock()

	if !userExists || user.Role != "admin" {
		c.JSON(403, ErrorResponse{Error: "Only admins can delete models"})
		return
	}

	var req struct {
		ID uint `json:"id" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		log.Printf("deleteModelHandler: BindJSON error: %v", err)
		c.JSON(400, ErrorResponse{Error: "Invalid request"})
		return
	}
	log.Printf("deleteModelHandler: request to delete model id=%d by user=%d", req.ID, userID)

	// First try to delete from DB and determine file path
	var filePath string
	if DB != nil {
		fid := int64(req.ID)
		// get file info
		fileName, _, archID, err := func() (string, string, sql.NullInt64, error) {
			// reuse GetModelRow which returns file_name, file_url, archive_id
			fn, fu, aid, err := GetModelRow(fid)
			return fn, fu, aid, err
		}()
		if err == nil {
			// determine path
			if archID.Valid {
				// get archive name
				aname, err2 := GetArchiveNameByID(archID.Int64)
				if err2 == nil {
					filePath = filepath.Join("model_archives", aname, fileName)
				} else {
					filePath = filepath.Join("model_archives", fileName)
				}
			} else {
				filePath = filepath.Join("uploads", fileName)
			}
		}
		// delete DB row
		_, _, err = DeleteModelByID(int64(req.ID))
		if err != nil {
			log.Printf("Warning: failed delete model row in sqlite: %v", err)
		}
	}

	// fallback / also remove from in-memory map and filesystem
	mu.Lock()
	model, ok := models[req.ID]
	if ok {
		// compute filePath if not set above
		if filePath == "" {
			if model.ArchiveID != 0 {
				arch := archives[model.ArchiveID]
				filePath = filepath.Join("model_archives", arch.Name, model.FileName)
			} else {
				filePath = filepath.Join("uploads", model.FileName)
			}
		}
		delete(models, req.ID)
		log.Printf("deleteModelHandler: model id=%d deleted from memory", req.ID)
	}
	mu.Unlock()

	if filePath != "" {
		if err := os.Remove(filePath); err != nil {
			log.Printf("Warning: Failed to delete file %s: %v\n", filePath, err)
		} else {
			log.Printf("deleteModelHandler: file removed %s", filePath)
		}
	}

	c.JSON(200, gin.H{"message": "Model deleted successfully"})
}

// ============ INIT DATA ============
func initData() {
	mu.Lock()
	defer mu.Unlock()

	// Admin user
	adminPass, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 10)
	admin := &User{
		ID:           1,
		Email:        "admin@test.com",
		PasswordHash: string(adminPass),
		Role:         "admin",
		CreatedAt:    time.Now(),
	}
	users[1] = admin

	// Regular user
	userPass, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	user := &User{
		ID:           2,
		Email:        "user@test.com",
		PasswordHash: string(userPass),
		Role:         "user",
		CreatedAt:    time.Now(),
	}
	users[2] = user

	fmt.Println("✅ Test data initialized")
	fmt.Println("   Admin: admin@test.com / admin123")
	fmt.Println("   User:  user@test.com / password123")
}

// ============ MAIN ============
func main() {
	// Create uploads directory if not exists
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0755)
	}

	// Initialize SQLite DB (path from env SQLITE_DB_PATH or default)
	dbPath := os.Getenv("SQLITE_DB_PATH")
	if dbPath == "" {
		dbPath = "./3d_db.db"
	}
	if err := InitDB(dbPath); err != nil {
		log.Printf("Warning: failed to open sqlite db %s: %v", dbPath, err)
	}

	// Initialize in-memory data
	initData()

	// Scan uploads directory and populate models map so uploaded files persist across restarts
	func() {
		mu.Lock()
		defer mu.Unlock()
		// ensure modelIDCounter continues from any existing entries
		// walk uploads dir
		filepath.Walk("uploads", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(info.Name()))
			if ext != ".glb" && ext != ".gltf" {
				return nil
			}

			// check if already present in models
			for _, m := range models {
				if m.FileName == info.Name() {
					return nil
				}
			}

			// derive a friendly name from filename (strip timestamp prefix if present)
			name := info.Name()
			if idx := strings.Index(name, "_"); idx != -1 {
				name = name[idx+1:]
			}
			name = strings.TrimSuffix(name, ext)

			model := &GLBModel{
				ID:          modelIDCounter,
				Name:        name,
				Description: "",
				FileURL:     fmt.Sprintf("/uploads/%s", info.Name()),
				FileName:    info.Name(),
				UploadedBy:  1, // unknown, mark as admin
				FileSize:    info.Size(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			models[modelIDCounter] = model
			modelIDCounter++
			return nil
		})
	}()

	// Scan model_archives directory: populate archives and models
	func() {
		mu.Lock()
		defer mu.Unlock()
		if _, err := os.Stat("model_archives"); os.IsNotExist(err) {
			return
		}
		// walk archive folders
		filepath.Walk("model_archives", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() {
				return nil
			}
			// skip root
			if path == "model_archives" || path == "model_archives." {
				return nil
			}
			// folder name
			name := filepath.Base(path)
			// read token if exists
			tokenPath := filepath.Join(path, "token.txt")
			var token string
			if b, err := os.ReadFile(tokenPath); err == nil {
				token = strings.TrimSpace(string(b))
			} else {
				// generate token and write
				t, _ := generateRandomToken(16)
				token = t
				_ = os.WriteFile(tokenPath, []byte(token), 0644)
			}
			// check if archive already exists
			var archID uint = 0
			for _, a := range archives {
				if a.Name == name {
					archID = a.ID
					break
				}
			}
			if archID == 0 {
				arch := &Archive{ID: archiveIDCounter, Name: name, Token: token, CreatedAt: time.Now()}
				archives[archiveIDCounter] = arch
				archID = archiveIDCounter
				archiveIDCounter++
			}
			// now list files inside folder and create model entries for glb/gltf
			files, _ := os.ReadDir(path)
			for _, f := range files {
				if f.IsDir() {
					continue
				}
				ext := strings.ToLower(filepath.Ext(f.Name()))
				if ext != ".glb" && ext != ".gltf" {
					continue
				}
				// skip token.txt
				if f.Name() == "token.txt" {
					continue
				}
				// check if model exists
				exists := false
				for _, m := range models {
					if m.FileName == f.Name() && m.ArchiveID == archID {
						exists = true
						break
					}
				}
				if exists {
					continue
				}
				info, _ := f.Info()
				nameOnly := f.Name()
				if idx := strings.Index(nameOnly, "_"); idx != -1 {
					nameOnly = nameOnly[idx+1:]
				}
				nameOnly = strings.TrimSuffix(nameOnly, ext)
				model := &GLBModel{ID: modelIDCounter, Name: nameOnly, Description: "", FileURL: fmt.Sprintf("/api/archives/%s/files/%s", name, f.Name()), FileName: f.Name(), ArchiveID: archID, UploadedBy: 1, FileSize: info.Size(), CreatedAt: time.Now(), UpdatedAt: time.Now()}
				models[modelIDCounter] = model
				modelIDCounter++
			}
			return nil
		})
	}()

	// Setup router
	router := gin.Default()

	// Setup CORS middleware
	router.Use(corsMiddleware())

	// Public routes
	router.POST("/api/auth/register", registerHandler)
	router.POST("/api/auth/login", loginHandler)
	router.GET("/api/models", getModelsHandler)
	router.Static("/uploads", "./uploads")
	// archive login (user token)
	router.POST("/api/archives/login", archiveLoginHandler)

	// archive file serving (secured)
	router.GET("/api/archives/:archiveName/files/:fileName", archiveAuthMiddleware(), archiveFileHandler)

	// Protected routes (admin)
	router.POST("/api/models/upload", authMiddleware(), uploadModelHandler)
	router.GET("/api/user/profile", authMiddleware(), getUserProfileHandler)
	router.DELETE("/api/models", authMiddleware(), deleteModelHandler)

	// Admin archive management
	router.POST("/api/archives", authMiddleware(), createArchiveHandler)
	router.GET("/api/archives", authMiddleware(), listArchivesHandler)
	router.DELETE("/api/archives", authMiddleware(), deleteArchiveHandler)

	fmt.Println("🚀 Server running on http://localhost:8080")
	router.Run(":8080")
}
