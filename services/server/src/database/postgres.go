package database

import (
	"fmt"

	"github.com/bananocoin/boompow-next/services/server/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

func NewConnection(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return db, err
	}
	return db, nil
}

func Migrate(db *gorm.DB) error {
	createTypes(db)
	return db.AutoMigrate(&models.User{}, &models.WorkRequest{})
}

// Create types in postgres
func createTypes(db *gorm.DB) error {
	result := db.Exec(fmt.Sprintf("SELECT 1 FROM pg_type WHERE typname = '%s';", models.PG_USER_TYPE_NAME))

	switch {
	case result.RowsAffected == 0:
		if err := db.Exec(fmt.Sprintf("CREATE TYPE %s AS ENUM ('%s', '%s');", models.PG_USER_TYPE_NAME, models.PROVIDER, models.REQUESTER)).Error; err != nil {
			fmt.Printf("Error creating %s ENUM", models.PG_USER_TYPE_NAME)
			return err
		}

		return nil
	case result.Error != nil:
		return result.Error

	default:
		return nil
	}
}
