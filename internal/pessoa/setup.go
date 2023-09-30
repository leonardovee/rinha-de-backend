package pessoa

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Setup(conn *pgxpool.Pool, cache *redis.Client, ch chan Schema) *Handler {
	repository := NewRepository(conn)
	service := NewService(repository, cache, ch)
	for i := 0; i < 8; i++ {
		go service.BatchInsert()
	}
	return NewHandler(service)
}
