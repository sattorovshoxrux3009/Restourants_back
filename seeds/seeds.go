package seeds

import (
	"context"
	"log"

	"github.com/sattorovshoxrux3009/Restourants_back/storage"
	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
	"golang.org/x/crypto/bcrypt"
)

// Seed interface defines methods for all seeds
type Seed interface {
	Run() error
}

// Seeds contains all seeds that should be run
type Seeds struct {
	storage storage.StorageI
}

// NewSeeds creates a new Seeds instance
func NewSeeds(storage storage.StorageI) *Seeds {
	return &Seeds{
		storage: storage,
	}
}

// RunAll runs all seeds
func (s *Seeds) RunAll() {
	if err := s.seedSuperAdmin(); err != nil {
		log.Printf("Error seeding super admin: %v\n", err)
	}
	// Add other seeds here in the future
}

// seedSuperAdmin creates a default super admin if it doesn't exist
func (s *Seeds) seedSuperAdmin() error {
	// Check if super admin already exists
	existingAdmin, err := s.storage.SuperAdmin().GetByUsername(context.Background(), "superadmin")
	if err == nil && existingAdmin != nil {
		log.Println("Super admin already exists, skipping creation")
		return nil
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("superadmin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create super admin
	_, err = s.storage.SuperAdmin().Create(context.Background(), &repo.SuperAdmin{
		Username:  "superadmin",
		FirstName: "Super",
		LastName:  "Admin",
		Password:  string(hashedPassword),
	})

	if err != nil {
		return err
	}

	log.Println("Super admin created successfully")
	log.Println("Username: superadmin")
	log.Println("Password: superadmin123")

	return nil
}
