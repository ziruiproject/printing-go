package database

import (
	"errors"
	"print-agent/internal/core/setting"
	"print-agent/pkg/helper"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type SettingRepository struct {
}

func NewSettingRepository() setting.Repository {
	return &SettingRepository{}
}

func (repository *SettingRepository) Save(db *gorm.DB, settings *setting.Setting) error {
	settings.ID = 1
	var existing setting.Setting
	err := db.First(&existing, 1).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := db.Create(settings).Error; err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}
	updated := helper.Differ(existing, settings).(setting.Setting)

	result := db.Save(&updated)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to save setting")
		return result.Error
	}

	return nil
}

func (repository *SettingRepository) Find(db *gorm.DB) (setting.Setting, error) {
	var setting setting.Setting
	result := db.First(&setting, 1)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to find setting")
		return setting, result.Error
	}

	return setting, nil
}
