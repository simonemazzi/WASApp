package database

import "database/sql"

type Avatar struct {
	Url    string `json:"url"`
	Mime   string `json:"mime"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Conversation struct {
	ConversationID int    `json:"conversation_id"`
	Name           string `json:"name"`
	Avatar         Avatar `json:"avatar"`
}

func (db *appdbimpl) GetConversations(userId int) ([]Conversation, error) {
	rows, err := db.c.Query(`SELECT
    c.conversationId,
    uu.username,
    p.url,
    p.mime,
    p.width,
    p.height
FROM Conversation c
JOIN UserUsername uu ON uu.userId = CASE
    WHEN c.component_A = ? THEN c.component_B
    ELSE c.component_A
END
AND uu.updateId = (
    SELECT MAX(updateId)
    FROM UserUsername
    WHERE userId = uu.userId
)
LEFT JOIN UsPhoto up ON up.userId = uu.userId
AND up.updateId = (
    SELECT up2.updateId
    FROM UsPhoto up2
    WHERE up2.userId = uu.userId
    ORDER BY up2.updateId DESC
    LIMIT 1
)
LEFT JOIN Photo p ON p.photoId = up.photoId
WHERE c.component_A = ? OR c.component_B = ?;`, userId, userId, userId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	conversations := make([]Conversation, 0)

	for rows.Next() {
		var conv Conversation

		err := rows.Scan(
			&conv.ConversationID,
			&conv.Name,
			&conv.Avatar.Url,
			&conv.Avatar.Mime,
			&conv.Avatar.Width,
			&conv.Avatar.Height,
		)
		if err != nil {
			return nil, err
		}

		conversations = append(conversations, conv)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}
