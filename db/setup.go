package postgresdb

import (
	"errors"

	"github.com/M1123Ananda/tododo/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Setup(dsn string) (*gorm.DB, error) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return DB, nil
}

func InitTables() error {
	if DB == nil {
		return errors.New("DB is not initialized")
	} else {
		err := DB.AutoMigrate(&model.User{})
		if err != nil {
			return errors.New("failed to init table: user")
		}
	}
	return nil
}