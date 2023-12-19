package main

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
)

func main() {
	// server
	e := echo.New()
	e.HideBanner = true
	e.Logger.SetLevel(log.INFO)
	e.Logger.SetHeader("${time_rfc3339} [${level}]")

	// cache
	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	// metrics
	elasticsearchClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://elasticsearch:9200"},
	})
	if err != nil {
		e.Logger.Fatalf("error creating the elasticsearch client: %v", err)
	}

	// routes
	e.GET("/:username", func(ctx echo.Context) error {
		metrics := NewMetrics(elasticsearchClient)
		handler := NewHandler(ctx, redisClient, metrics)
		return handler.User()
	})

	// start
	e.Logger.Fatal(e.Start(":8080"))
}
