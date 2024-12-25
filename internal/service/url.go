package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Dashboard/urlShortener/internal/model"
	"github.com/Dashboard/urlShortener/internal/repo"
	"time"
)

type ShortCodeGenerator interface {
	GenerateShortCode() string
}

type Cacher interface {
	SetURL(ctx context.Context, url repo.Url) error
	GetURL(ctx context.Context, shortCode string) (*repo.Url, error)
}

type URLService struct {
	querier         repo.Querier
	shortener       ShortCodeGenerator
	defaultDuration time.Duration
	cache           Cacher
	baseURL         string
}

func NewURLService(db *sql.DB, shortCodeGenerator ShortCodeGenerator, duration time.Duration, cache Cacher, baseURL string) *URLService {
	return &URLService{
		querier:         repo.New(db),
		shortener:       shortCodeGenerator,
		defaultDuration: duration,
		cache:           cache,
		baseURL:         baseURL,
	}
}

func (s *URLService) CreateURL(ctx context.Context, req model.CreateURLRequest) (*model.CreateURLResponse, error) {
	var shortCode string
	var isCustom bool
	var expiredTime time.Time
	if req.CustomCode != "" {
		isAvailable, err := s.querier.IsShortCodeAvailable(ctx, req.CustomCode)
		if err != nil {
			return nil, err
		}
		if !isAvailable {
			return nil, errors.New("别名已存在")
		}
		//可得到
		shortCode = req.CustomCode
		isCustom = true
	} else {
		//别名不存在
		code, err := s.createShortCode(ctx, 0)
		if err != nil {
			return nil, err
		}
		shortCode = code
	}
	shanghaiLoc, _ := time.LoadLocation("Asia/Shanghai")
	rd := &req.Duration
	if rd == nil {
		expiredTime = time.Now().In(shanghaiLoc).Add(s.defaultDuration)
	} else {
		expiredTime = time.Now().In(shanghaiLoc).Add(time.Hour * time.Duration(*rd))
	}
	//插入数据库
	url, err := s.querier.CreateURL(ctx, repo.CreateURLParams{
		OriginalUrl: req.OriginalURL,
		ShortCode:   shortCode,
		IsCustom:    isCustom,
		ExpiredTime: expiredTime,
	})
	if err != nil {
		return nil, err
	}

	//插入缓存
	if err := s.cache.SetURL(ctx, url); err != nil {
		return nil, err
	}
	return &model.CreateURLResponse{
		ShortURL:    s.baseURL + "/" + url.ShortCode,
		ExpireTime:  expiredTime.Format("2006-01-02 15:04:05"),
		OriginalURL: req.OriginalURL,
	}, nil

}

func (s *URLService) GetURL(ctx context.Context, shortCode string) (string, error) {
	// 先访问cache
	url, err := s.cache.GetURL(ctx, shortCode)
	if err != nil {
		return "", err
	}
	if url != nil {
		return url.OriginalUrl, nil
	}

	//访问数据库
	url2, err := s.querier.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}
	//存入缓存
	if err := s.cache.SetURL(ctx, url2); err != nil {
		return "", err
	}
	return url2.OriginalUrl, nil
}

func (s *URLService) createShortCode(ctx context.Context, n int) (string, error) {
	if n > 5 {
		return "", errors.New("重试过多")
	}
	shortCode := s.shortener.GenerateShortCode()
	isAvailable, err := s.querier.IsShortCodeAvailable(ctx, shortCode)
	if err != nil {
		return "", err
	}
	if isAvailable {
		return shortCode, nil
	}

	return s.createShortCode(ctx, n+1)
}

func (s *URLService) DeleteURL(ctx context.Context) error {
	return s.querier.DeleteURLExpired(ctx)
}
