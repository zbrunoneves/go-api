package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	ctx     echo.Context
	cache   *redis.Client
	metrics *Metrics
}

func NewHandler(ctx echo.Context, redisClient *redis.Client, metrics *Metrics) *Handler {
	return &Handler{
		ctx:     ctx,
		cache:   redisClient,
		metrics: metrics,
	}
}

func (h *Handler) User() error {
	name := h.ctx.Param("username")

	h.metrics.Add("path", "/:username")
	h.metrics.Add("method", h.ctx.Request().Method)

	if userData, err := h.cache.Get(h.ctx.Request().Context(), name).Result(); err == nil {
		h.metrics.Add("redis_cache", "hit")

		var user User
		err = json.Unmarshal([]byte(userData), &user)
		if err != nil {
			h.metrics.SendError(err)
			return echo.ErrInternalServerError
		}
		h.metrics.AddAll(user.ToMetrics())

		h.metrics.Send()
		return h.ctx.JSON(http.StatusOK, user)
	}
	h.metrics.Add("redis_cache", "miss")

	user := User{
		Name:    name,
		SavedAt: time.Now().Format(time.TimeOnly),
	}
	h.metrics.AddAll(user.ToMetrics())

	userData, err := json.Marshal(&user)
	if err != nil {
		h.metrics.SendError(err)
		return echo.ErrInternalServerError
	}

	err = h.cache.Set(h.ctx.Request().Context(), user.Name, userData, 2*time.Minute).Err()
	if err != nil {
		h.metrics.SendError(err)
		return echo.ErrInternalServerError
	}

	h.metrics.Send()
	return h.ctx.JSON(http.StatusOK, user)
}
