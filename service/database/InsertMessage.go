package database

func (db *appdbimpl) InsertMessage(conversationId int, userId int, text string, photoId *int) error {

	if photoId != nil {
		_, err := db.c.Exec(`INSERT INTO Message(conversationId,sender,text,photoId) VALUES (?,?,?,?)`, conversationId, userId, text, *photoId)
		if err != nil {
			return err
		}
	} else {
		_, err := db.c.Exec(`INSERT INTO Message(conversationId,sender,text) VALUES (?,?,?)`, conversationId, userId, text)
		if err != nil {
			return err
		}
	}
	return nil

}

func (db *appdbimpl) InsertGroupMessage(groupId int, userId int, text string, photoId *int) error {

	if photoId != nil {
		_, err := db.c.Exec(`INSERT INTO Message(groupId,sender,text,photoId) VALUES (?,?,?,?)`, groupId, userId, text, *photoId)
		if err != nil {
			return err
		}
	} else {
		_, err := db.c.Exec(`INSERT INTO Message(groupId,sender,text) VALUES (?,?,?)`, groupId, userId, text)
		if err != nil {
			return err
		}
	}
	return nil

}
