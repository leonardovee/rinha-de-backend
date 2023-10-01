package pessoa

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	repository *Repository
	cache      *redis.Client
	ch         chan Schema
}

func NewService(repository *Repository, cache *redis.Client, ch chan Schema) *Service {
	return &Service{repository: repository, cache: cache, ch: ch}
}

func (s *Service) BatchInsert() {
	sl := make([]Schema, 0)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case sch := <-s.ch:
			sl = append(sl, sch)
		case <-ticker.C:
			if len(sl) > 0 {
				err := s.repository.InsertBatch(sl)
				if err != nil {
					fmt.Println(err.Error())
				}
				sl = make([]Schema, 0)
			}
		}
	}
}

func (s *Service) InsertPessoa(ctx context.Context, cpr *CreateRequest) (Schema, error) {
	val, _ := s.cache.Get(ctx, cpr.Apelido).Result()
	if val != "" {
		return Schema{}, errors.New("duplicated entry")
	}
	schema := Schema{
		ID:         uuid.New().String(),
		Nome:       cpr.Nome,
		Apelido:    cpr.Apelido,
		Nascimento: cpr.Nascimento,
		Stack:      cpr.Stack,
	}
	go func() {
		s.ch <- schema
	}()
	cacheValue, _ := json.Marshal(schema)
	s.cache.Set(ctx, schema.ID, cacheValue, 0)
	s.cache.Set(ctx, cpr.Apelido, cpr.Apelido, 0)
	return schema, nil
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
