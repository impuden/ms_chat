package models

import "time"

type Message struct {
	Username  uint64    `json:"username"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type CoreData struct {
	User1  uint64 `json:"user_id1"`
	User2  uint64 `json:"user_id2"`
	ItemID uint64 `json:"item_id"`
}
