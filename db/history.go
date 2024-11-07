package db

import (
	"chat-service/models"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func SaveMessage(db *sql.DB, username uint64, room uint64, message string) error {
	query := "INSERT INTO messages (username, room, message, timestamp) VALUES (?, ?, ?, NOW())"
	_, err := db.Exec(query, username, room, message)
	if err != nil {
		log.Println("error saving msg to DB: ", err)
		return err
	}
	fmt.Println("savemsg was acivated")
	return nil
}

func LoadMessages(db *sql.DB, room uint64, offset int, limit int) ([]models.Message, error) {
	query := "SELECT username, room, message, timestamp FROM messages WHERE room = ? ORDER BY timestamp ASC LIMIT ? OFFSET ?"
	rows, err := db.Query(query, room, limit, offset)
	if err != nil {
		log.Println("error load history: ", err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		var timestamp string
		err := rows.Scan(&msg.Username, &room, &msg.Message, &timestamp)
		if err != nil {
			log.Println("error scan message db.LoadMessages: ", err)
			return nil, err
		}

		msg.Timestamp, err = time.Parse("2006-01-02 15:04:05", timestamp)
		if err != nil {
			log.Println("err pars timestamp: ", err)
			return nil, err
		}

		messages = append(messages, msg)
	}
	log.Println("loaded msg1: ")
	return messages, nil
}
