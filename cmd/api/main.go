package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"leonardovee.com/rinha-de-backend/internal/api"
	"leonardovee.com/rinha-de-backend/internal/pessoa"
	"os"
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

	p := pessoa.Setup(pool, cache)

	e := echo.New()
	e.Use(middleware.Logger())

	api.Setup(e, &api.Handlers{Pessoa: p})

	e.Logger.Fatal(e.Start(":80"))
}
