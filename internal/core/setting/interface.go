package setting

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Save(db *gorm.DB, setting *Setting) error
	Find(db *gorm.DB) (Setting, error)
}

type Usecase interface {
	Save(ctx context.Context, request Content) (Content, error)
	Find(ctx context.Context) (*Content, error)
}
