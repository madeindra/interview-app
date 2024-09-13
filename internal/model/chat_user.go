package model

import (
	"github.com/google/uuid"
)

type ChatUser struct {
	ID       string `json:"id"`
	Secret   string `json:"secret"`
	Language string `json:"language"`
}

func (m *Model) CreateChatUser(secret, language string) (*ChatUser, error) {
	id := uuid.New().String()
	_, err := m.conn.Exec("INSERT INTO chat_users (id, secret, language) VALUES (?, ?, ?)", id, secret, language)
	if err != nil {
		return nil, err
	}

	return &ChatUser{ID: id, Secret: secret}, nil
}

func (m *Model) GetChatUser(id string) (*ChatUser, error) {
	var user ChatUser
	err := m.conn.QueryRow("SELECT id, secret, language FROM chat_users WHERE id = ?", id).Scan(&user.ID, &user.Secret, &user.Language)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
