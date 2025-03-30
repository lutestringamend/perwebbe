package config

import (
	"log"

	"github.com/lutestringamend/perwebbe/internal/model"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	log.Println("running db migration")

	err := db.AutoMigrate(
		&model.BlogPost{},
		&model.Tag{},
		&model.PortfolioProject{},
		&model.ContactSubmission{},
	)
	if err != nil {
		return err
	}

	log.Println("db migrations completed successfully")
	return nil
}
