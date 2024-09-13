package elevenlabs

import (
	"io"

	"github.com/madeindra/interview-app/internal/elevenlabs/model"
)

type Client interface {
	TextToSpeech(string) (io.ReadCloser, error)
}

type ElevenLab struct {
	baseURL  string
	ttsModel string
	ttsVoice string
}

const (
	baseURL  = "https://api.elevenlabs.io/v1"
	ttsModel = "eleven_multilingual_v2"
	ttsVoice = "cgSgspJ2msm6clMCkdW9"
)

var defaultVoiceSetting = model.VoiceSetting{
	Stability:       0.5,
	SimilarityBoost: 0.75,
}

func New() *ElevenLab {
	return &ElevenLab{
		baseURL:  baseURL,
		ttsModel: ttsModel,
		ttsVoice: ttsVoice,
	}
}
