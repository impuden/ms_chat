package db

import "database/sql"

func CheckCreateRoom(db *sql.DB, user1, user2, itemID uint64) (uint64, error) {
	var id uint64

	query := `SELECT id FROM rooms
	WHERE (user1 = ? AND user2 = ? AND itemID = ?)
	OR (user1 = ? AND user2 = ? AND itemID = ?)
	LIMIT 1
	`

	err := db.QueryRow(query, user1, user2, itemID, user2, user1, itemID).Scan(&id)
	if err == sql.ErrNoRows {
		insertQuery := `
	INSERT INTO rooms (user1, user2, itemID)
	VALUES (?, ?, ?)
	`

		result, err := db.Exec(insertQuery, user1, user2, itemID)
		if err != nil {
			return 0, err
		}

		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}

		return uint64(lastInsertID), nil
	} else if err != nil {
		return 0, err
	}

	return uint64(id), nil
}
