package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Utility script untuk seed database dengan test data
// Sekarang menggunakan in-memory storage di main.go
func seedDatabase() {
	fmt.Println("Using in-memory storage. Test data initialized in initData() function.")

	// Hash passwords untuk reference
	adminPass, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	userPass, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	fmt.Printf("Admin hash: %s\n", string(adminPass))
	fmt.Printf("User hash: %s\n", string(userPass))
	fmt.Println("\nTest Credentials:")
	fmt.Println("  Admin:  admin@test.com / admin123")
	fmt.Println("  User:   user@test.com / password123")
}
