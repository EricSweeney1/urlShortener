package application

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Dashboard/urlShortener/config"
	"github.com/Dashboard/urlShortener/database"
	"github.com/Dashboard/urlShortener/internal/api"
	"github.com/Dashboard/urlShortener/internal/cache"
	"github.com/Dashboard/urlShortener/internal/service"
	"github.com/Dashboard/urlShortener/pkg/shortCode"
	"github.com/Dashboard/urlShortener/pkg/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application struct {
	e                  *echo.Echo
	db                 *sql.DB
	redisClient        *cache.RedisCache
	urlService         *service.URLService
	urlHandler         *api.URLHandler
	config             *config.Config
	shortCodeGenerator *shortCode.ShortCode
}

func (a *Application) Init(filePath string) error {
	config, err := config.LoadConfig(filePath)
	if err != nil {
		return fmt.Errorf("加载配置错误: %w", err)
	}
	a.config = config
	db, err := database.NewDB(config.Database)
	if err != nil {
		return err
	}
	a.db = db
	redisClient, err := cache.NewRedisCache(config.Redis)
	if err != nil {
		return err
	}
	a.redisClient = redisClient
	a.shortCodeGenerator = shortCode.NewShortCode(config.ShortCode.Length)
	BaseURL := config.App.BaseHost + config.App.BasePort
	a.urlService = service.NewURLService(db, a.shortCodeGenerator, config.App.DefaultDuration, redisClient, BaseURL)
	a.urlHandler = api.NewURLHandler(a.urlService)
	e := echo.New()
	e.Server.WriteTimeout = config.Server.WriteTimeout
	e.Server.ReadTimeout = config.Server.ReadTimeout
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.POST("/api/url", a.urlHandler.CreateURL)
	e.GET("/:code", a.urlHandler.RedirectURL)
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	e.Validator = validator.NewCustomValidator()
	a.e = e
	return nil
}
func (a *Application) Run() {
	go a.startServer()
	go a.cleanUp()
	a.Shutdown(context.Background())
}

func (a *Application) startServer() {
	if err := a.e.Start(a.config.Server.Address); err != nil {
		log.Panicln(err)
	}
}

func (a *Application) cleanUp() {
	ticker := time.NewTicker(a.config.App.CleanUpInterval)
	defer ticker.Stop()
	for range ticker.C {
		if err := a.urlService.DeleteURL(context.Background()); err != nil {
			log.Println(err)
		}
	}
}

func (a *Application) Shutdown(ctx context.Context) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	defer func() {
		if err := a.db.Close(); err != nil {
			log.Println(err)
		}
	}()
	defer func() {
		if err := a.redisClient.Close(); err != nil {
			log.Println(err)
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.e.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}
