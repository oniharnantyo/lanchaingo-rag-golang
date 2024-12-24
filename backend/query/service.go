package query

import (
	"context"
	"github.com/google/uuid"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"
	mongomemory "github.com/tmc/langchaingo/memory/mongo"
	"github.com/tmc/langchaingo/schema"
)

type Service struct {
	llm            llms.Model
	retriever      schema.Retriever
	mongoDBUri     string
	mongoDBName    string
	collectionName string
}

func NewService(
	llm llms.Model,
	retriever schema.Retriever,
	mongoDBUri string,
	mongoDBName string,
	collectionName string,
) *Service {
	return &Service{llm, retriever, mongoDBUri, mongoDBName, collectionName}
}

func (s *Service) Query(ctx context.Context, sessionId string, question string) (*QueryResponse, error) {
	if sessionId == "" {
		sessionId = uuid.New().String()
	}

	mongoHistory, err := mongomemory.NewMongoDBChatMessageHistory(ctx,
		mongomemory.WithConnectionURL(s.mongoDBUri),
		mongomemory.WithDataBaseName(s.mongoDBName),
		mongomemory.WithSessionID(sessionId),
		mongomemory.WithCollectionName(s.collectionName),
	)
	if err != nil {
		return nil, err
	}

	conversationBuffer := memory.NewConversationBuffer(memory.WithChatHistory(mongoHistory))

	chain := chains.NewConversationalRetrievalQAFromLLM(s.llm, s.retriever, conversationBuffer)
	res, err := chains.Run(ctx, chain, question)
	if err != nil {
		return nil, err
	}

	return &QueryResponse{
		SessionId: sessionId,
		Answer:    res,
	}, nil
}
