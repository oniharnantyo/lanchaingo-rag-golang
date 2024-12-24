package document

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/tmc/langchaingo/documentloaders"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service}
}

func (h *Handler) AddDocument(c echo.Context) error {
	ctx := c.Request().Context()

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Error retrieving file")
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(src)
	if err != nil {
		return err
	}

	readerAt := bytes.NewReader(buf.Bytes())

	pdf := documentloaders.NewPDF(readerAt, file.Size)
	if err = h.service.AddDocument(ctx, &pdf); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, "Document added")
}
