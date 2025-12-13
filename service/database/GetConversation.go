package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetConversationById(conversationId int) (Conversation, error) {
	rows, err := db.c.Query(`SELECT
    c.conversationId,
    uu.username,
    p.url,
    p.mime,
    p.width,
    p.height
FROM Conversation c
JOIN UserUsername uu ON uu.userId = c.component_B
AND uu.updateId = (
    SELECT MAX(updateId)
    FROM UserUsername
    WHERE userId = c.component_B
)
LEFT JOIN UsPhoto up ON up.userId = c.component_B
AND up.updateId = (
    SELECT MAX(updateId)
    FROM UsPhoto
    WHERE userId = c.component_B
)
LEFT JOIN Photo p ON p.photoId = up.photoId
WHERE c.conversationId = ?`, conversationId)
	if err != nil {
		return Conversation{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	if !rows.Next() {
		return Conversation{}, errors.New("no conversation found")
	}

	var conv Conversation
	err = rows.Scan(
		&conv.ConversationID,
		&conv.Name,
		&conv.Avatar.Url,
		&conv.Avatar.Mime,
		&conv.Avatar.Width,
		&conv.Avatar.Height,
	)
	if err != nil {
		return Conversation{}, err
	}

	if err = rows.Err(); err != nil {
		return Conversation{}, err
	}

	return conv, nil

}
