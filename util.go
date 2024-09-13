package main

import (
	"regexp"
	"strings"

	"github.com/madeindra/interview-app/internal/model"
	oaiModel "github.com/madeindra/interview-app/internal/openai/model"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

func sanitizeString(text string) string {
	reStrong := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	text = reStrong.ReplaceAllString(text, "$1")

	reItalic := regexp.MustCompile(`\*([^*]+)\*`)
	text = reItalic.ReplaceAllString(text, "$1")

	reLink := regexp.MustCompile(`\[(.*?)\]\(.*?\)`)
	text = reLink.ReplaceAllString(text, "$1")

	reBullet := regexp.MustCompile(`\n- `)
	text = reBullet.ReplaceAllString(text, ", ")

	text = strings.Replace(text, "\n", " ", -1)

	return text
}

func generateRandom() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 10

	random := make([]byte, length)
	for i := range random {
		random[i] = charset[rand.Intn(len(charset))]
	}

	return string(random)
}

func createHash(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func compareHash(plain, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}

func entryToChatMessage(chats []model.Entry) []oaiModel.ChatMessage {
	var chatHistory []oaiModel.ChatMessage
	for _, chat := range chats {
		chatHistory = append(chatHistory, oaiModel.ChatMessage{
			Role:    oaiModel.Role(chat.Role),
			Content: chat.Text,
		})
	}

	return chatHistory
}

func pointer[T any](v T) *T {
	return &v
}
