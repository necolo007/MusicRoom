package dbs

import (
	"github.com/necolo007/MusicRoom/internal/app/user/user_entity"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&user_entity.User{},
	)
	return err
}
