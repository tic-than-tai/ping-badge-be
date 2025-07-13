package database

import (
	"ping-badge-be/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Initialize(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(
		&model.User{},
		&model.Organization{},
		&model.OrganizationAdmin{},
		&model.Badge{},
		&model.IssuedBadge{},
		&model.UserBadgeProgress{},
		&model.Activity{},
		&model.ActivityParticipation{},
		&model.BadgeView{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
