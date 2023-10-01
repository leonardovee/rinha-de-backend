package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"leonardovee.com/rinha-de-backend/internal/api"
	"leonardovee.com/rinha-de-backend/internal/pessoa"
)

func main() {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	cache := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})

	ch := make(chan pessoa.Schema)
	p := pessoa.Setup(pool, cache, ch)

	e := echo.New()
	e.Use(echoprometheus.NewMiddleware("rinha_de_backend"))
	//e.Use(middleware.Logger())
	e.GET("/metrics", echoprometheus.NewHandler())

	api.Setup(e, &api.Handlers{Pessoa: p})

	e.Logger.Fatal(e.Start(":80"))
}
