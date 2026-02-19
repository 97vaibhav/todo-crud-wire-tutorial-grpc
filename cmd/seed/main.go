// cmd/seed seeds the database with the two default groups and an initial admin user.
// Run this ONCE after `make migrate-up`. It is idempotent — safe to run multiple times.
package main

import (
	"fmt"
	"log"

	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/config"
	infradb "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/infrastructure/db"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Stable UUIDs so the seed is idempotent: running it again just hits the unique constraint.
const (
	adminGroupID = "00000000-0000-0000-0000-000000000001"
	guestGroupID = "00000000-0000-0000-0000-000000000002"

	defaultAdminEmail    = "admin@example.com"
	defaultAdminPassword = "admin123"
)

// Inline structs are fine here — the seed command is not part of the app,
// so we don't need the domain/repository abstraction.
type groupSeed struct {
	ID   string `gorm:"column:id"`
	Name string `gorm:"column:name"`
	Type string `gorm:"column:type"`
}

func (groupSeed) TableName() string { return "groups" }

type userSeed struct {
	ID           string `gorm:"column:id"`
	Name         string `gorm:"column:name"`
	Email        string `gorm:"column:email"`
	PasswordHash string `gorm:"column:password_hash"`
	GroupID      string `gorm:"column:group_id"`
}

func (userSeed) TableName() string { return "users" }

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := infradb.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	seedGroups(db)
	seedAdminUser(db)

	fmt.Println("\nSeed complete.")
	fmt.Printf("  Admin login → email: %s  password: %s\n", defaultAdminEmail, defaultAdminPassword)
	fmt.Println("  Change the password after first login!")
}

func seedGroups(db *gorm.DB) {
	groups := []groupSeed{
		{ID: adminGroupID, Name: "Admins", Type: "ADMIN"},
		{ID: guestGroupID, Name: "Guests", Type: "GUEST"},
	}
	for _, g := range groups {
		// ON CONFLICT DO NOTHING makes this idempotent.
		db.Exec(
			`INSERT INTO groups (id, name, type) VALUES (?, ?, ?) ON CONFLICT (id) DO NOTHING`,
			g.ID, g.Name, g.Type,
		)
		fmt.Printf("  group seeded: %s (%s)\n", g.Name, g.Type)
	}
}

func seedAdminUser(db *gorm.DB) {
	hash, err := bcrypt.GenerateFromPassword([]byte(defaultAdminPassword), 12)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	db.Exec(
		`INSERT INTO users (id, name, email, password_hash, group_id)
		 VALUES (gen_random_uuid(), ?, ?, ?, ?)
		 ON CONFLICT (email) DO NOTHING`,
		"Admin", defaultAdminEmail, string(hash), adminGroupID,
	)
	fmt.Printf("  user seeded: %s\n", defaultAdminEmail)
}
