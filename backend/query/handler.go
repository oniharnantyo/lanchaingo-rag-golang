package query

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Query(c echo.Context) error {
	ctx := c.Request().Context()

	var req QueryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	res, err := h.service.Query(ctx, req.SessionId, req.Query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": res})
}
