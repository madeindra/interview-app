package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/madeindra/interview-app/internal/openai/model"
)

func (ai *OpenAI) IsKeyValid(apiKey string) (bool, error) {
	url, err := url.JoinPath(ai.baseURL, "/models")
	if err != nil {
		return false, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	return true, nil
}

func (c *OpenAI) Status() (model.Status, error) {
	url, err := url.JoinPath(statusURL, "/components.json")
	if err != nil {
		return model.STATUS_UNKNOWN, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return model.STATUS_UNKNOWN, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.STATUS_UNKNOWN, err
	}

	if resp.StatusCode != http.StatusOK {
		return model.STATUS_UNKNOWN, nil
	}

	var statusResp model.ComponentStatusResponse
	err = unmarshalJSONResponse(resp, &statusResp)
	if err != nil {
		return model.STATUS_UNKNOWN, err
	}

	for _, component := range statusResp.Components {
		if component.Name == "API" {
			switch component.Status {
			case "operational":
				return model.STATUS_OPERATIONAL, nil
			case "degraded_performance":
				return model.STATUS_DEGRADED_PERFORMANCE, nil
			case "partial_outage":
				return model.STATUS_PARTIAL_OUTAGE, nil
			case "major_outage":
				return model.STATUS_MAJOR_OUTAGE, nil
			}
		}
	}

	return model.STATUS_UNKNOWN, nil
}

func (ai *OpenAI) Chat(apiKey string, messages []model.ChatMessage) (model.ChatResponse, error) {
	url, err := url.JoinPath(ai.baseURL, "/chat/completions")
	if err != nil {
		log.Default().Println("error joining url path", err)

		return model.ChatResponse{}, err
	}

	chatReq := model.ChatRequest{
		Model:    chatModel,
		Messages: messages,
	}

	body, err := json.Marshal(chatReq)
	if err != nil {
		log.Default().Println("error marshalling chat request", err)

		return model.ChatResponse{}, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Default().Println("error creating http request", err)

		return model.ChatResponse{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Default().Println("error sending http request", err)

		return model.ChatResponse{}, err
	}

	var chatResp model.ChatResponse
	err = unmarshalJSONResponse(resp, &chatResp)
	if err != nil {
		log.Default().Println("error unmarshalling chat response", err)

		return model.ChatResponse{}, err
	}

	return chatResp, nil
}

func (ai *OpenAI) Transcribe(apiKey string, file io.Reader, filename string) (model.TranscriptResponse, error) {
	if file == nil {
		log.Default().Println("audio is nil")

		return model.TranscriptResponse{}, fmt.Errorf("audio is nil")
	}

	url, err := url.JoinPath(baseURL, "/audio/transcriptions")
	if err != nil {
		log.Default().Println("error joining url path", err)

		return model.TranscriptResponse{}, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		log.Default().Println("error creating form file", err)

		return model.TranscriptResponse{}, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		log.Default().Println("error copying file to form file", err)

		return model.TranscriptResponse{}, err
	}

	err = writer.WriteField("model", transcriptModel)
	if err != nil {
		log.Default().Println("error writing model field", err)

		return model.TranscriptResponse{}, err
	}

	err = writer.WriteField("language", transcriptLanguage)
	if err != nil {
		log.Default().Println("error writing language field", err)

		return model.TranscriptResponse{}, err
	}

	err = writer.Close()
	if err != nil {
		log.Default().Println("error closing writer", err)

		return model.TranscriptResponse{}, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, body)
	if err != nil {
		log.Default().Println("error creating http request", err)

		return model.TranscriptResponse{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Default().Println("error sending http request", err)

		return model.TranscriptResponse{}, err
	}

	if resp == nil || resp.Body == nil {
		log.Default().Println("response is nil")

		return model.TranscriptResponse{}, fmt.Errorf("response is nil")
	}

	var transcriptResp model.TranscriptResponse
	err = unmarshalJSONResponse(resp, &transcriptResp)
	if err != nil {
		log.Default().Println("error unmarshalling transcript response", err)

		return model.TranscriptResponse{}, err
	}

	return transcriptResp, nil
}

func (ai *OpenAI) Speechify(apiKey string, text string) (io.ReadCloser, error) {
	url, err := url.JoinPath(ai.baseURL, "/audio/speech")
	if err != nil {
		log.Default().Println("error joining url path", err)

		return nil, err
	}

	ttsReq := model.TTSRequest{
		Model: ttsModel,
		Voice: ttsVoice,
		Input: text,
	}

	body, err := json.Marshal(ttsReq)
	if err != nil {
		log.Default().Println("error marshalling tts request", err)

		return nil, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Default().Println("error creating http request", err)

		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Default().Println("error sending http request", err)

		return nil, err
	}
	respBody, err := getResponseBody(resp)
	if err != nil {
		log.Default().Println("error getting response body", err)

		return nil, err
	}

	return respBody, nil
}
