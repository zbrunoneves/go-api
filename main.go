package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
)

type User struct {
	Name    string `json:"name"`
	SavedAt string `json:"saved_at"`
}

func main() {
	e := echo.New()
	e.HideBanner = true
	e.Logger.SetLevel(log.INFO)
	e.Logger.SetHeader("${time_rfc3339} [${level}]")

	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	handler := NewHandler(redisClient)
	e.GET("/:name", handler.Name)

	e.Logger.Fatal(e.Start(":8080"))
}
