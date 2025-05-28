package dbs

import (
	"github.com/necolo007/MusicRoom/internal/app/music/music_entity"
	"github.com/necolo007/MusicRoom/internal/app/user/user_entity"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&user_entity.User{},
		&music_entity.Music{},
		&music_entity.MusicDTO{},
	)
	return err
}
