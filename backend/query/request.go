package query

type QueryRequest struct {
	SessionId string `json:"sessionId"`
	Query     string `json:"query"`
}
