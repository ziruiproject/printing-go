package database

import (
	"print-agent/config"
	"print-agent/internal/core/setting"

	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func NewSqliteConn() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.AppConfig.DBName))
	if err != nil {
		log.Panic().
			Err(err).
			Str("database", config.AppConfig.DBName).
			Msg("unable to connect to database")

		return nil, err
	}

	err = db.AutoMigrate(setting.Setting{})

	return db, nil
}
