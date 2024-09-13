package model

type StartChatResponse struct {
	ID       string `json:"id"`
	Secret   string `json:"secret"`
	Language string `json:"language"`

	Chat
}

type AnswerChatResponse struct {
	Language string `json:"language"`
	Prompt   Chat   `json:"prompt,omitempty"`
	Answer   Chat   `json:"answer,omitempty"`
}

type StatusResponse struct {
	Server    bool   `json:"server"`
	Key       bool   `json:"key"`
	API       *bool  `json:"api"`
	ApiStatus string `json:"apiStatus"`
}
