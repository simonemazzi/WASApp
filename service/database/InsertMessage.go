package database

import "database/sql"

func (db *appdbimpl) InsertMessage(conversationId int, userId int, text string, photoId *int) (Message, error) {

	var res sql.Result
	var err error

	if photoId != nil {
		res, err = db.c.Exec(
			`INSERT INTO Message (conversationId, sender, text, photoId)
			 VALUES (?, ?, ?, ?)`,
			conversationId, userId, text, *photoId,
		)
	} else {
		res, err = db.c.Exec(
			`INSERT INTO Message (conversationId, sender, text)
			 VALUES (?, ?, ?)`,
			conversationId, userId, text,
		)
	}

	if err != nil {
		return Message{}, err
	}

	messageId, err := res.LastInsertId()
	if err != nil {
		return Message{}, err
	}

	// ritorno il mess appena creato
	var msg Message

	var photoUrl sql.NullString
	var photoWidth, photoHeight sql.NullInt64
	var photoMime sql.NullString

	err = db.c.QueryRow(`
		SELECT
			m.messageId,
			m.text,

			p.URL,
			p.width,
			p.height,
			p.mime,

			m.time,

			u.userId,
			uu.username,

			m.originalMessage IS NOT NULL

		FROM Message m

		JOIN User u
			ON u.userId = m.sender

		JOIN UserUsername uu
			ON uu.userId = u.userId
		   AND uu.updateId = (
				SELECT MAX(updateId)
				FROM UserUsername
				WHERE userId = u.userId
		   )

		LEFT JOIN Photo p
			ON p.photoId = m.photoId

		WHERE m.messageId = ?
	`, messageId).Scan(
		&msg.MessageId,
		&msg.Body.Text,

		&photoUrl,
		&photoWidth,
		&photoHeight,
		&photoMime,

		&msg.Time,

		&msg.Sender.UserID,
		&msg.Sender.Username,

		&msg.IsForwarded,
	)

	if err != nil {
		return Message{}, err
	}

	// se c'è la foto la metto
	if photoUrl.Valid {
		msg.Body.Photo = &Photo{
			Url:    photoUrl.String,
			Width:  int(photoWidth.Int64),
			Height: int(photoHeight.Int64),
			Mime:   photoMime.String,
		}
	}

	return msg, nil
}
func (db *appdbimpl) InsertGroupMessage(
	groupId int,
	userId int,
	text string,
	photoId *int,
) (Message, error) {

	var res sql.Result
	var err error

	if photoId != nil {
		res, err = db.c.Exec(
			`INSERT INTO Message (groupId, sender, text, photoId)
			 VALUES (?, ?, ?, ?)`,
			groupId, userId, text, *photoId,
		)
	} else {
		res, err = db.c.Exec(
			`INSERT INTO Message (groupId, sender, text)
			 VALUES (?, ?, ?)`,
			groupId, userId, text,
		)
	}

	if err != nil {
		return Message{}, err
	}

	messageId, err := res.LastInsertId()
	if err != nil {
		return Message{}, err
	}

	var msg Message

	var photoUrl sql.NullString
	var photoWidth, photoHeight sql.NullInt64
	var photoMime sql.NullString

	err = db.c.QueryRow(`
		SELECT
			m.messageId,
			m.text,

			p.URL,
			p.width,
			p.height,
			p.mime,

			m.time,

			u.userId,
			uu.username,

			m.originalMessage IS NOT NULL

		FROM Message m

		JOIN User u
			ON u.userId = m.sender

		JOIN UserUsername uu
			ON uu.userId = u.userId
		   AND uu.updateId = (
				SELECT MAX(updateId)
				FROM UserUsername
				WHERE userId = u.userId
		   )

		LEFT JOIN Photo p
			ON p.photoId = m.photoId

		WHERE m.messageId = ?
	`, messageId).Scan(
		&msg.MessageId,
		&msg.Body.Text,

		&photoUrl,
		&photoWidth,
		&photoHeight,
		&photoMime,

		&msg.Time,

		&msg.Sender.UserID,
		&msg.Sender.Username,

		&msg.IsForwarded,
	)

	if err != nil {
		return Message{}, err
	}

	if photoUrl.Valid {
		msg.Body.Photo = &Photo{
			Url:    photoUrl.String,
			Width:  int(photoWidth.Int64),
			Height: int(photoHeight.Int64),
			Mime:   photoMime.String,
		}
	}

	return msg, nil
}
