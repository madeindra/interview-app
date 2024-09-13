package openai

type OpenAI struct {
	baseURL            string
	chatModel          string
	transcriptModel    string
	transcriptLanguage string
	ttsModel           string
	ttsVoice           string
}

const (
	baseURL            = "https://api.openai.com/v1"
	statusURL          = "https://status.openai.com/api/v2"
	chatModel          = "gpt-4o-mini-2024-07-18"
	transcriptModel    = "whisper-1"
	transcriptLanguage = "en"
	ttsModel           = "tts-1"
	ttsVoice           = "nova"
)

var supportedTranscriptLanguages = map[string]struct{}{
	"en": {},
}

func New() *OpenAI {
	return &OpenAI{
		baseURL:            baseURL,
		chatModel:          chatModel,
		transcriptModel:    transcriptModel,
		transcriptLanguage: transcriptLanguage,
		ttsModel:           ttsModel,
		ttsVoice:           ttsVoice,
	}
}
