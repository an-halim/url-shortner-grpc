package repository

import (
	"context"

	"github.com/an-halim/url-shortner-grpc/entity"
	"gorm.io/gorm"
)

type UrlRepository struct {
	db *gorm.DB
}

func NewUrlRepository(db *gorm.DB) *UrlRepository {
	return &UrlRepository{db: db}
}

func (r *UrlRepository) Short(ctx context.Context, shortUrl entity.Url) (entity.Url, error) {
	err := r.db.Create(&shortUrl).Error
	if err != nil {
		return entity.Url{}, err
	}
	return shortUrl, nil
}

func (r *UrlRepository) GetByShort(ctx context.Context, shortUrl string) (entity.Url, error) {
	var entity entity.Url
	err := r.db.Where("short_url = ?", shortUrl).First(&entity).Error

	if err != nil {
		return entity, err
	}

	return entity, nil
}
