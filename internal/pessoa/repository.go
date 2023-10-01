package pessoa

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	Conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) *Repository {
	return &Repository{Conn: conn}
}

func (r Repository) GetById(id string) (Schema, error) {
	selectStatement := `
		SELECT id, apelido, nome, nascimento, stack
		FROM pessoas
		WHERE id = $1;
	`

	var pessoa struct {
		ID         string
		Apelido    string
		Nome       string
		Nascimento string
		Stack      string
	}

	err := r.Conn.QueryRow(context.Background(), selectStatement, id).
		Scan(&pessoa.ID, &pessoa.Apelido, &pessoa.Nome, &pessoa.Nascimento, &pessoa.Stack)

	if err != nil {
		return Schema{}, err
	}

	return Schema{
		ID:         pessoa.ID,
		Apelido:    pessoa.Apelido,
		Nome:       pessoa.Nome,
		Nascimento: pessoa.Nascimento,
		Stack:      strings.Split(pessoa.Stack, ", "),
	}, nil
}

func (r Repository) GetByTermo(t string) ([]Schema, error) {
	selectStatement := `
		SELECT id, apelido, nome, nascimento, stack
		FROM pessoas
		WHERE trigram LIKE $1 LIMIT 50;
	`

	rows, err := r.Conn.Query(context.Background(), selectStatement, fmt.Sprintf("%%%v%%", strings.ToLower(t)))

	if err != nil {
		return []Schema{}, err
	}

	var schemas []Schema

	for rows.Next() {
		var schema struct {
			ID         string
			Apelido    string
			Nome       string
			Nascimento string
			Stack      string
		}
		err = rows.Scan(&schema.ID, &schema.Apelido, &schema.Nome, &schema.Nascimento, &schema.Stack)
		schemas = append(schemas, Schema{
			ID:         schema.ID,
			Apelido:    schema.Apelido,
			Nome:       schema.Nome,
			Nascimento: schema.Nascimento,
			Stack:      strings.Split(schema.Stack, ", "),
		})
	}

	return schemas, nil
}

func (r Repository) GetCount() (int, error) {
	selectStatement := `
		SELECT count(*)
		FROM pessoas;
	`

	var count int

	err := r.Conn.QueryRow(context.Background(), selectStatement).
		Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r Repository) InsertBatch(pessoas []Schema) error {
	_, err := r.Conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"pessoas"},
		[]string{"id", "apelido", "nome", "nascimento", "stack"},
		pgx.CopyFromSlice(len(pessoas), func(i int) ([]any, error) {
			p := pessoas[i]
			stack := strings.Join(p.Stack, ", ")
			return []any{p.ID, p.Apelido, p.Nome, p.Nascimento, stack}, nil
		}),
	)

	if err != nil {
		return err
	}

	return nil
}
