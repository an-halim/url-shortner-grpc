package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/an-halim/url-shortner-grpc/entity"
	"github.com/an-halim/url-shortner-grpc/repository"
	"github.com/redis/go-redis/v9"
	"github.com/teris-io/shortid"
)

type UrlService struct {
	UrlRepository *repository.UrlRepository
	redis         *redis.Client
}

var redisKey = "key:id"

func NewUrlService(urlRepository *repository.UrlRepository, redis *redis.Client) *UrlService {
	return &UrlService{UrlRepository: urlRepository, redis: redis}
}

func (s *UrlService) Short(ctx context.Context, shortUrl entity.Url) (entity.Url, error) {

	sid, err := shortid.New(10, shortid.DefaultABC, 1231241)

	if err != nil {
		log.Println("Gagal membuat short id")
		return entity.Url{}, fmt.Errorf("Gagal membuat short id")
	}
	urlEntity := entity.Url{
		Original: shortUrl.Original,
		ShortUrl: sid.MustGenerate(),
	}

	created, err := s.UrlRepository.Short(ctx, urlEntity)
	if err != nil {
		log.Println("Gagal Membuat Short url")
		return entity.Url{}, fmt.Errorf("Gagal membuat short url")
	}

	// set to redis
	redisKey := fmt.Sprintf(redisKey, created.ShortUrl)

	data, err := json.Marshal(created)

	if err != nil {
		log.Println("Error marshal data redis")
	}

	if err := s.redis.Set(ctx, redisKey, data, 60*time.Second).Err(); err != nil {
		log.Println("Error set data redis")
	}

	return created, nil
}

func (s *UrlService) GetByShort(ctx context.Context, shortUrl string) (entity.Url, error) {
	var url entity.Url

	redisKey := fmt.Sprintf(redisKey, shortUrl)
	// get from redis
	val, err := s.redis.Get(ctx, redisKey).Result()
	if err == nil {
		log.Println("Data ditemukan di redis cache")
		err = json.Unmarshal([]byte(val), &url)
		if err != nil {
			log.Println("Error unmarshal data redis")
		}
	}

	if err != nil {
		log.Println("Data tidak ditemukan di redis cache")
		url, err = s.UrlRepository.GetByShort(ctx, shortUrl)
		if err != nil {
			return entity.Url{}, err
		}

		return url, nil
	}

	return url, nil
}
