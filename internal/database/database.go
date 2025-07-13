package database

import (
	"ping-badge-be/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Initialize(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.Organization{},
		&models.OrganizationAdmin{},
		&models.Badge{},
		&models.IssuedBadge{},
		&models.UserBadgeProgress{},
		&models.Activity{},
		&models.ActivityParticipation{},
		&models.BadgeView{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
