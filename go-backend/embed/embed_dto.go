package embedservice

type embedTextReq struct {
	Text  string `json:"text"`
	Model string `json:"model"`
}

type embedTextResp struct {
	Embedding []float32 `json:"embedding"`
}
