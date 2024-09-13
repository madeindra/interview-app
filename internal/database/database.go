package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

const (
	settingsSchema = `
	CREATE TABLE IF NOT EXISTS settings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		openai_key VARCHAR DEFAULT '',
		elevenlabs_key VARCHAR DEFAULT '',
	);`

	settingsData = "SELECT id, openai_key, elevenlabs_key FROM settings LIMIT 1;"

	settingInsert = "INSERT INTO settings (id) VALUES (1);"

	chatUsersSchema = `CREATE TABLE IF NOT EXISTS chat_users (
		id VARCHAR PRIMARY KEY,
		secret VARCHAR NOT NULL,
		language VARCHAR DEFAULT 'en'
	);`

	chatsSchema = `CREATE TABLE IF NOT EXISTS chats (
		id VARCHAR PRIMARY KEY,
		chat_user_id VARCHAR,
		role VARCHAR,
		text VARCHAR,
		audio VARCHAR,
		FOREIGN KEY(chat_user_id) REFERENCES chat_users(id)
	);`
)

func New() *sql.DB {
	db, err := sql.Open("sqlite", "file:app.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func migrate(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(settingsSchema)
	if err != nil {
		log.Fatal(err)
	}

	var id int
	var openaiKey, elevenlabsKey string
	err = tx.QueryRow(settingsData).Scan(&id, &openaiKey, &elevenlabsKey)
	if err != nil && err == sql.ErrNoRows {
		log.Default().Println("No settings found, creating default settings")

		if _, err := tx.Exec(settingInsert); err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec(chatUsersSchema)
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec(chatsSchema)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
