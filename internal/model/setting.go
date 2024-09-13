package model

import (
	"database/sql"
	"errors"
)

func (m *Model) GetAPIKey() (string, string, error) {
	var openAIKey, elevenLabsKey string
	err := m.conn.QueryRow("SELECT openai_key, elevenlabs_key FROM settings LIMIT 1").Scan(&openAIKey, &elevenLabsKey)

	return openAIKey, elevenLabsKey, err
}

func (m *Model) AreKeyExist() (bool, error) {
	openAIKey, elevenLabsKey, err := m.GetAPIKey()
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return openAIKey != "" || elevenLabsKey != "", nil
}

func (m *Model) UpdateAPIKeys(openAIKey, elevenLabsKey string) error {
	query := "UPDATE settings SET "
	args := []any{}

	// update openai_key if provided
	if openAIKey != "" {
		query += "openai_key = ?, "
		args = append(args, openAIKey)
	}

	// update elevenlabs_key if provided
	if elevenLabsKey != "" {
		query += "elevenlabs_key = ?, "
		args = append(args, elevenLabsKey)
	}

	// Remove trailing comma and space
	query = query[:len(query)-2]
	query += " WHERE id = 1"

	_, err := m.conn.Exec(query, args...)
	return err
}
