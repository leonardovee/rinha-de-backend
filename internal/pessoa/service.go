package pessoa

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	repository *Repository
	cache      *redis.Client
}

func NewService(repository *Repository, cache *redis.Client) *Service {
	return &Service{repository: repository, cache: cache}
}

func (s *Service) InsertPessoa(ctx context.Context, cpr *CreateRequest) (Schema, error) {
	val, _ := s.cache.Get(ctx, cpr.Apelido).Result()
	if val != "" {
		return Schema{}, errors.New("duplicated entry")
	}
	schema, err := s.repository.Insert(Schema{
		ID:         uuid.New().String(),
		Nome:       cpr.Nome,
		Apelido:    cpr.Apelido,
		Nascimento: cpr.Nascimento,
		Stack:      cpr.Stack,
	})
	cacheValue, _ := json.Marshal(schema)
	s.cache.Set(ctx, schema.ID, cacheValue, 0)
	s.cache.Set(ctx, cpr.Apelido, cpr.Apelido, 0)
	return schema, err
}

func (s *Service) GetPessoaById(ctx context.Context, id string) (Schema, error) {
	val, _ := s.cache.Get(ctx, id).Result()
	if val != "" {
		var schema Schema
		err := json.Unmarshal([]byte(val), &schema)
		if err == nil {
			return schema, nil
		}
	}
	return s.repository.GetById(id)
}

func (s *Service) GetPessoasByTermo(ctx context.Context, t string) ([]Schema, error) {
	val, _ := s.cache.Get(ctx, t).Result()
	if val != "" {
		var schema []Schema
		err := json.Unmarshal([]byte(val), &schema)
		if err == nil {
			return schema, nil
		}
	}
	schema, err := s.repository.GetByTermo(t)
	cacheValue, _ := json.Marshal(schema)
	s.cache.Set(ctx, t, cacheValue, 0)
	return schema, err
}

func (s *Service) GetPessoaCount() (int, error) {
	return s.repository.GetCount()
}
