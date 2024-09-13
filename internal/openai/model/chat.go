package model

type ChatMessage struct {
	Content string `json:"content"`
	Role    Role   `json:"role"`
}

type ChatRequest struct {
	Messages []ChatMessage `json:"messages"`
	Model    string        `json:"model"`
}

type ChatResponse struct {
	Choices []Choice `json:"choices"`
}

type TTSRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
	Voice string `json:"voice"`
}

type TranscriptResponse struct {
	Text string `json:"text"`
}
