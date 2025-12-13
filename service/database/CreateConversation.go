package database

func (db *appdbimpl) CreateConversation(userId int, username string, time string) (Conversation, error) {
	componentB, err := db.SearchUserByUsername(username, time)
	if err != nil {
		return Conversation{}, err
	}

	res, err := db.c.Exec(` INSERT INTO Conversation (component_A,component_B) VALUES (?,?)
`, userId, componentB)
	if err != nil {
		if sqliteErr, ok := err.(interface{ Error() string }); ok && sqliteErr.Error() == "Conversation already created" {
			var conv Conversation
			row := db.c.QueryRow(`
                SELECT conversationId
                FROM Conversation
                WHERE (component_A = ? AND component_B = ?)
                   OR (component_A = ? AND component_B = ?)
            `, userId, componentB, componentB, userId)

			var id int
			if err := row.Scan(&id); err != nil {
				return Conversation{}, err
			}
			conv.ConversationID = id
			conv.Name = username
			conv.Avatar, err = db.GetPhoto(componentB, time)
			if err != nil {
				return Conversation{}, err
			}
			return conv, nil
		}

		return Conversation{}, err
	}
	var conversation Conversation
	var id int64
	id, err = res.LastInsertId()
	conversation.ConversationID = int(id)
	if err != nil {
		return Conversation{}, err
	}
	conversation.Name = username
	conversation.Avatar, err = db.GetPhoto(componentB, time)
	if err != nil {
		return Conversation{}, err
	}

	return conversation, nil
}
