package pessoa

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Setup(conn *pgxpool.Pool, cache *redis.Client) *Handler {
	repository := NewRepository(conn)
	service := NewService(repository, cache)
	return NewHandler(service)
}
