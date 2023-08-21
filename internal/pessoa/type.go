package pessoa

type CreateRequest struct {
	Apelido    string   `json:"apelido" validate:"max=32"`
	Nome       string   `json:"nome" validate:"max=100"`
	Nascimento string   `json:"nascimento" validate:"max=10"`
	Stack      []string `json:"stack,omitempty" validate:"max=32,dive,max=32"`
}

type Schema struct {
	ID         string   `json:"id"`
	Apelido    string   `json:"apelido"`
	Nome       string   `json:"nome"`
	Nascimento string   `json:"nascimento"`
	Stack      []string `json:"stack,omitempty" validate:"max=5,dive,max=32"`
}
