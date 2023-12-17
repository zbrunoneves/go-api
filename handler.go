package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	redisClient *redis.Client
}

func NewHandler(redisClient *redis.Client) *Handler {
	return &Handler{
		redisClient: redisClient,
	}
}

func (h *Handler) Name(c echo.Context) error {
	name := c.Param("name")

	if userData, err := h.redisClient.Get(c.Request().Context(), name).Result(); err == nil {
		c.Logger().Info(fmt.Sprintf("Cache HIT key=%s", name))

		var user User
		err = json.Unmarshal([]byte(userData), &user)
		if err != nil {
			c.Logger().Error(err)
			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, user)
	}

	c.Logger().Info(fmt.Sprintf("Cache MISS key=%s", name))

	user := User{
		Name:    name,
		SavedAt: time.Now().Format("15:04:05"),
	}

	userData, err := json.Marshal(&user)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	err = h.redisClient.Set(c.Request().Context(), user.Name, userData, 2*time.Minute).Err()
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, user)
}
