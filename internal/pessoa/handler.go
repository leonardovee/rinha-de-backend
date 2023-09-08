package pessoa

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Get(c echo.Context) error {
	t := c.QueryParam("t")
	if t == "" {
		return c.JSON(http.StatusBadRequest, nil)
	}
	pessoas, err := h.service.GetPessoasByTermo(c.Request().Context(), t)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return c.JSON(http.StatusNotFound, nil)
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	if len(pessoas) == 0 {
		return c.JSON(http.StatusOK, "[]")
	}
	return c.JSON(http.StatusOK, pessoas)
}

func (h *Handler) GetTotal(c echo.Context) error {
	total, err := h.service.GetPessoaCount()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, total)
}

func (h *Handler) GetById(c echo.Context) error {
	id := c.Param("id")
	pessoa, err := h.service.GetPessoaById(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return c.JSON(http.StatusNotFound, nil)
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, pessoa)
}

func (h *Handler) Post(c echo.Context) error {
	cpr := new(CreateRequest)
	if err := c.Bind(cpr); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "")
	}
	pessoa, err := h.service.InsertPessoa(c.Request().Context(), cpr)
	if err != nil {
		if err.Error() == "duplicated entry" {
			return c.JSON(http.StatusUnprocessableEntity, nil)
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	locationHeader := fmt.Sprintf("/pessoas/%s", pessoa.ID)
	c.Response().Header().Set("Location", locationHeader)
	return c.JSON(http.StatusCreated, "")
}
