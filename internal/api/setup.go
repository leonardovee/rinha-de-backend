package api

import (
	"github.com/labstack/echo/v4"
	"leonardovee.com/rinha-de-backend/internal/pessoa"
)

type Handlers struct {
	Pessoa *pessoa.Handler
}

func Setup(e *echo.Echo, h *Handlers) {
	e.GET("/pessoas", h.Pessoa.Get)
	e.GET("/pessoas/:id", h.Pessoa.GetById)
	e.POST("/pessoas", h.Pessoa.Post)
	e.GET("/contagem-pessoas", h.Pessoa.GetTotal)
}
