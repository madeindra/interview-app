package elevenlabs

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/madeindra/interview-app/internal/elevenlabs/model"
)

func (c *ElevenLab) Speechify(apiKey, input string) (io.ReadCloser, error) {
	url, err := url.JoinPath(c.baseURL, "text-to-speech", c.ttsVoice)
	if err != nil {
		return nil, err
	}

	ttsReq := model.TTSRequest{
		Text:         input,
		ModelID:      c.ttsModel,
		VoiceSetting: defaultVoiceSetting,
	}

	body, err := json.Marshal(ttsReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("xi-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	respBody, err := getResponseBody(resp)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
