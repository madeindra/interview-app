package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/madeindra/interview-app/internal/language"
	"github.com/madeindra/interview-app/internal/model"
	oaiModel "github.com/madeindra/interview-app/internal/openai/model"
)

func (a *App) AreKeyExist() (bool, error) {
	return a.model.AreKeyExist()
}

func (a *App) UpdateAPIKeys(oaiKey, elKey string) error {
	return a.model.UpdateAPIKeys(oaiKey, elKey)
}

func (a *App) Status() (model.StatusResponse, error) {
	oaiKey, _, err := a.model.GetAPIKey()
	if err != nil {
		return model.StatusResponse{}, fmt.Errorf("failed to get api key: %v", err)
	}

	isKeyValid, err := a.oaiAPI.IsKeyValid(oaiKey)
	if err != nil {
		return model.StatusResponse{}, fmt.Errorf("failed to check api key: %v", err)
	}

	status, err := a.oaiAPI.Status()
	if err != nil {
		return model.StatusResponse{}, fmt.Errorf("failed to get api status: %v", err)
	}

	var apiState *bool

	switch status {
	case oaiModel.STATUS_OPERATIONAL:
		apiState = pointer(true)
	case oaiModel.STATUS_DEGRADED_PERFORMANCE, oaiModel.STATUS_PARTIAL_OUTAGE, oaiModel.STATUS_MAJOR_OUTAGE:
		apiState = pointer(false)
	case oaiModel.STATUS_UNKNOWN:
		apiState = nil
	}

	response := model.StatusResponse{
		Server:    true,           // always true when the server is running
		Key:       isKeyValid,     // true if the API key is valid, false otherwise
		API:       apiState,       // nil if status unknown, true if operational, false otherwise
		ApiStatus: string(status), // always return the status string
	}

	return response, nil
}

func (a *App) StartChat(role string, skills []string, lang string) (model.StartChatResponse, error) {
	oaiKey, elKey, err := a.model.GetAPIKey()
	if err != nil {
		return model.StartChatResponse{}, fmt.Errorf("failed to get api key: %v", err)
	}

	chatLanguage := a.oaiAPI.GetDefaultTranscriptLanguage()
	if lang != "" {
		chatLanguage = language.GetLanguage(lang)
	}

	systempPrompt, err := a.oaiAPI.GetSystemPrompt(role, skills, chatLanguage)
	if err != nil {
		return model.StartChatResponse{}, fmt.Errorf("failed to get system prompt: %v", err)
	}

	initialText, err := a.oaiAPI.GetInitialChat(role, chatLanguage)
	if err != nil {
		return model.StartChatResponse{}, fmt.Errorf("failed to get initial text: %v", err)
	}

	var initialAudio io.Reader
	if a.oaiAPI.IsSpeechAvailable(chatLanguage) {
		speech, err := a.oaiAPI.Speechify(oaiKey, sanitizeString(initialText))
		if err != nil {
			return model.StartChatResponse{}, fmt.Errorf("failed to create initial audio: %v", err)
		}

		initialAudio = speech
	} else {
		speech, err := a.elAPI.Speechify(elKey, sanitizeString(initialText))
		if err != nil {
			return model.StartChatResponse{}, fmt.Errorf("failed to create initial audio: %v", err)
		}

		initialAudio = speech
	}

	var audioBase64 string
	if initialAudio != nil {
		audioByte, err := io.ReadAll(initialAudio)
		if err != nil {
			return model.StartChatResponse{}, fmt.Errorf("failed to read speech: %v", err)
		}

		audioBase64 = base64.StdEncoding.EncodeToString(audioByte)
	}

	plainSecret := generateRandom()
	hashed, err := createHash(plainSecret)
	if err != nil {
		return model.StartChatResponse{}, fmt.Errorf("failed to create hash: %v", err)
	}

	newUser, err := a.model.CreateChatUser(hashed, chatLanguage)
	if err != nil {
		return model.StartChatResponse{}, fmt.Errorf("failed to create new chat: %v", err)
	}

	if _, err := a.model.CreateChat(newUser.ID, string(oaiModel.ROLE_SYSTEM), systempPrompt, audioBase64); err != nil {
		return model.StartChatResponse{}, fmt.Errorf("failed to create chat: %v", err)
	}

	if _, err := a.model.CreateChat(newUser.ID, string(oaiModel.ROLE_ASSISTANT), initialText, audioBase64); err != nil {
		return model.StartChatResponse{}, fmt.Errorf("failed to create chat: %v", err)
	}

	initialChat := model.StartChatResponse{
		ID:       newUser.ID,
		Secret:   plainSecret,
		Language: lang,
		Chat: model.Chat{
			Text:  initialText,
			Audio: audioBase64,
		},
	}

	return initialChat, nil
}

func (a *App) AnswerChat(userID, userSecret string, audioData []byte) (model.AnswerChatResponse, error) {
	apiKey, elKey, err := a.model.GetAPIKey()
	if err != nil {
		return model.AnswerChatResponse{}, fmt.Errorf("failed to get api key: %v", err)
	}

	user, err := a.model.GetChatUser(userID)
	if err != nil {
		return model.AnswerChatResponse{}, fmt.Errorf("failed to get chat: %v", err)
	}

	if err := compareHash(userSecret, user.Secret); err != nil {
		return model.AnswerChatResponse{}, fmt.Errorf("invalid user secret")
	}

	entry, err := a.model.GetChatsByChatUserID(userID)
	if err != nil {
		return model.AnswerChatResponse{}, fmt.Errorf("failed to get chat: %v", err)
	}

	audioReader := bytes.NewReader(audioData)
	transcript, err := a.oaiAPI.Transcribe(apiKey, audioReader, "audio.wav")
	if err != nil {
		return model.AnswerChatResponse{}, fmt.Errorf("failed to transcribe audio: %v", err)
	}

	if transcript.Text == "" {
		return model.AnswerChatResponse{}, fmt.Errorf("cannot complete audio transcription: no transcript")
	}

	if _, err := a.model.CreateChat(userID, string(oaiModel.ROLE_USER), transcript.Text, base64.StdEncoding.EncodeToString(audioData)); err != nil {
		return model.AnswerChatResponse{}, fmt.Errorf("failed to create chat: %v", err)
	}

	chatHistory := append(entry, model.Entry{
		ChatUserID: userID,
		Role:       string(oaiModel.ROLE_USER),
		Text:       transcript.Text,
	})

	chatMessages := entryToChatMessage(chatHistory)

	chatCompletion, err := a.oaiAPI.Chat(apiKey, chatMessages)
	if err != nil {
		return model.AnswerChatResponse{}, fmt.Errorf("failed to get chat completion: %v", err)
	}

	if len(chatCompletion.Choices) == 0 {
		return model.AnswerChatResponse{}, fmt.Errorf("cannot complete chat completion: no chat completion")
	}

	speechText := chatCompletion.Choices[0].Message.Content

	var speech io.Reader
	if a.oaiAPI.IsSpeechAvailable(user.Language) {
		speech, err = a.oaiAPI.Speechify(apiKey, sanitizeString(chatCompletion.Choices[0].Message.Content))
		if err != nil {
			return model.AnswerChatResponse{}, fmt.Errorf("failed to create speech: %v", err)
		}
	} else {
		speech, err = a.elAPI.Speechify(elKey, sanitizeString(chatCompletion.Choices[0].Message.Content))
		if err != nil {
			return model.AnswerChatResponse{}, fmt.Errorf("failed to create speech: %v", err)
		}
	}

	var speechBase64 string
	if speech != nil {
		speechByte, err := io.ReadAll(speech)
		if err != nil {
			return model.AnswerChatResponse{}, fmt.Errorf("failed to read speech: %v", err)
		}

		speechBase64 = base64.StdEncoding.EncodeToString(speechByte)
	}

	if _, err := a.model.CreateChat(userID, string(oaiModel.ROLE_ASSISTANT), speechText, speechBase64); err != nil {
		return model.AnswerChatResponse{}, fmt.Errorf("failed to create chat: %v", err)
	}

	response := model.AnswerChatResponse{
		Language: language.GetCode(user.Language),
		Prompt: model.Chat{
			Text: transcript.Text,
		},
		Answer: model.Chat{
			Text:  speechText,
			Audio: speechBase64,
		},
	}

	return response, nil
}

func (a *App) EndChat(userID, userSecret string) (model.AnswerChatResponse, error) {
	apiKey, elKey, err := a.model.GetAPIKey()
	if err != nil {
		return model.AnswerChatResponse{}, fmt.Errorf("failed to get api key: %v", err)
	}

	user, err := a.model.GetChatUser(userID)
	if err != nil {
		return model.AnswerChatResponse{}, fmt.Errorf("failed to get chat: %v", err)
	}

	if err := compareHash(userSecret, user.Secret); err != nil {
		return model.AnswerChatResponse{}, fmt.Errorf("invalid user secret")
	}

	entry, err := a.model.GetChatsByChatUserID(userID)
	if err != nil {
		return model.AnswerChatResponse{}, fmt.Errorf("failed to get chat: %v", err)
	}

	chatHistory := append(entry, model.Entry{
		ChatUserID: userID,
		Role:       string(oaiModel.ROLE_USER),
		Text:       "That is the end of the mock interview, thank you, please provide your feedbacks on my strength and which area to improve, and whether you are confident that I fits the role.",
	})

	chatMessages := entryToChatMessage(chatHistory)

	chatCompletion, err := a.oaiAPI.Chat(apiKey, chatMessages)
	if err != nil {
		return model.AnswerChatResponse{}, fmt.Errorf("failed to get chat completion: %v", err)
	}

	if len(chatCompletion.Choices) == 0 {
		return model.AnswerChatResponse{}, fmt.Errorf("cannot complete chat completion: no chat completion")
	}

	speechText := chatCompletion.Choices[0].Message.Content

	var speech io.Reader
	if a.oaiAPI.IsSpeechAvailable(user.Language) {
		speech, err = a.oaiAPI.Speechify(apiKey, sanitizeString(chatCompletion.Choices[0].Message.Content))
		if err != nil {
			return model.AnswerChatResponse{}, fmt.Errorf("failed to create speech: %v", err)
		}
	} else {
		speech, err = a.elAPI.Speechify(elKey, sanitizeString(chatCompletion.Choices[0].Message.Content))
		if err != nil {
			return model.AnswerChatResponse{}, fmt.Errorf("failed to create speech: %v", err)
		}
	}

	var speechBase64 string
	if speech != nil {
		speechByte, err := io.ReadAll(speech)
		if err != nil {
			return model.AnswerChatResponse{}, fmt.Errorf("failed to read speech: %v", err)
		}

		speechBase64 = base64.StdEncoding.EncodeToString(speechByte)
	}

	if _, err := a.model.CreateChat(userID, string(oaiModel.ROLE_ASSISTANT), speechText, speechBase64); err != nil {
		return model.AnswerChatResponse{}, fmt.Errorf("failed to create chat: %v", err)
	}

	response := model.AnswerChatResponse{
		Language: language.GetCode(user.Language),
		Answer: model.Chat{
			Text:  speechText,
			Audio: speechBase64,
		},
	}

	return response, nil
}
