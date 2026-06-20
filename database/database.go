package database

import (
	"fmt"
	"strings"

	"hotelmanagementsystem.com/v2/config"
	"hotelmanagementsystem.com/v2/models"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(driver string, dsn string, sqlitePath string) (*gorm.DB, error) {
	switch strings.ToLower(strings.TrimSpace(driver)) {
	case "", "sqlite":
		return gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
	case "postgres", "postgresql":
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}
}

func Seed(db *gorm.DB, cfg config.Config) error {
	var userCount int64
	if err := db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		return err
	}

	if userCount == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte(cfg.SeedAdminPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		admin := models.User{
			Name:         cfg.SeedAdminName,
			Email:        cfg.SeedAdminEmail,
			PasswordHash: string(hash),
			Role:         cfg.SeedAdminRole,
		}

		if err := db.Create(&admin).Error; err != nil {
			return err
		}
	}

	return nil
}
