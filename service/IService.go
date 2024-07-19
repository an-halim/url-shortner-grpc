package service

import (
	"context"

	"github.com/an-halim/url-shortner-grpc/entity"
)

type IService interface {
	Short(ctx context.Context, shortUrl entity.Url) (entity.Url, error)
	GetByShort(ctx context.Context, shortUrl string) (entity.Url, error)
}
