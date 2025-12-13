package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) UserConversation(userId int, conversationId int) (bool, error) {
	var tmp int
	err := db.c.QueryRow(`
        SELECT 1
        FROM Conversation c
        WHERE c.conversationId = ?
          AND (c.component_A = ? OR c.component_B = ?)
    `, conversationId, userId, userId).Scan(&tmp)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
