package api

import (
	"context"
	"github.com/Dashboard/urlShortener/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

type URLService interface {
	CreateURL(ctx context.Context, req model.CreateURLRequest) (*model.CreateURLResponse, error)
	GetURL(ctx context.Context, shortCode string) (string, error)
}

type URLHandler struct {
	URLService URLService
}

func NewURLHandler(s URLService) *URLHandler {
	return &URLHandler{URLService: s}
}

// CreateURL POST /api/url original_url, custom_code, duration, -> shortUrl,expired_time
func (h *URLHandler) CreateURL(c echo.Context) error {
	//提取数据
	var req model.CreateURLRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	//验证数据格式
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	//调用业务函数
	resp, err := h.URLService.CreateURL(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	//返回响应
	return c.JSON(http.StatusCreated, resp)
}

// RedirectURL GET /:code redirect
func (h *URLHandler) RedirectURL(c echo.Context) error {
	// 把code取出来
	shortCode := c.Param("code")
	// shortcode -> url调用业务函数
	originalURL, err := h.URLService.GetURL(c.Request().Context(), shortCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.Redirect(http.StatusPermanentRedirect, originalURL)
}
