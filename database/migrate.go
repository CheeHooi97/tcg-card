package database

import (
	"pkm/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	models := []any{
		&model.Admin{},
		&model.BGS{},
		&model.BGSUrl{},
		&model.BGSLogging{},
		&model.Card{},
		&model.CardPrice{},
		&model.CGC{},
		&model.CGCUrl{},
		&model.CGCLogging{},
		&model.PSA{},
		&model.PSAUrl{},
		&model.PSALogging{},
		&model.Set{},
		&model.TAG{},
		&model.TAGUrl{},
		&model.TAGLogging{},
		&model.Token{},
		&model.User{},
		&model.UserCard{},
		&model.UserCardSearchLog{},
		&model.UserDevice{},
	}
	err := db.AutoMigrate(models...)
	if err != nil {
		return err
	}
	return nil
}
