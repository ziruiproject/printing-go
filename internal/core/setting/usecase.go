package setting

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UsecaseDependency struct {
	DB                *gorm.DB
	SettingRepository Repository
}

type UsecaseImpl struct {
	UsecaseDependency
}

func NewUsecase(deps UsecaseDependency) Usecase {
	return &UsecaseImpl{
		deps,
	}
}

func (usecase *UsecaseImpl) Save(ctx context.Context, request Content) (Content, error) {
	tx := usecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	setting := request.ToEntity()
	err := usecase.SettingRepository.Save(tx, setting)
	if err = tx.Commit().Error; err != nil {
		log.Error().Err(err).Msgf("failed to commit transaction: %+v", err)
		return Content{}, errors.New("something went wrong")
	}

	return *ToContent(setting), nil
}

func (usecase *UsecaseImpl) Find(ctx context.Context) (*Content, error) {
	tx := usecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	setting, err := usecase.SettingRepository.Find(tx)
	if err != nil {
		log.Error().Err(err).Msg("failed to find setting")
		return nil, errors.New("setting not found")
	}

	if err = tx.Commit().Error; err != nil {
		log.Error().Err(err).Msgf("failed to commit transaction: %+v", err)
		return nil, errors.New("something went wrong")
	}
	return ToContent(&setting), nil
}
