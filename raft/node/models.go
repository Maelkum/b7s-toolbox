package node

type SolveRequest struct {
	Expression string             `json:"expression,omitempty"`
	Parameters map[string]float64 `json:"parameters,omitempty"`
}

type SolveResponse struct {
	Expression string             `json:"expression,omitempty"`
	Parameters map[string]float64 `json:"parameters,omitempty"`
	Result     float64            `json:"result,omitempty"`
	Error      string             `json:"error,omitempty"`
}
