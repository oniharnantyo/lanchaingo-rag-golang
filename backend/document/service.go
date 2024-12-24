package document

import (
	"bytes"
	"context"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"strings"
	"unicode/utf8"
)

type Service struct {
	splitter textsplitter.TextSplitter
	store    vectorstores.VectorStore
}

func NewService(splitter textsplitter.TextSplitter, store vectorstores.VectorStore) *Service {
	return &Service{splitter, store}
}

func (s *Service) AddDocument(ctx context.Context, pdf *documentloaders.PDF) error {
	documents, err := pdf.LoadAndSplit(ctx, s.splitter)
	if err != nil {
		return err
	}

	for i, document := range documents {
		documents[i].PageContent = sanitizeUTF8([]byte(strings.ReplaceAll(document.PageContent, "\x00", "")))
	}

	_, err = s.store.AddDocuments(ctx, documents)
	if err != nil {
		return err
	}

	return nil
}

func sanitizeUTF8(input []byte) string {
	if !utf8.Valid(input) {
		decoder := charmap.Windows1252.NewDecoder()
		result, _, err := transform.Bytes(decoder, input)
		if err != nil {
			return string(bytes.ReplaceAll(input, []byte{0}, []byte("")))
		}
		return string(result)
	}
	return string(input)
}
