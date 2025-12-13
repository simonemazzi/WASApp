package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

type User struct {
	Username string `json:"username"`
	UserID   string `json:"userId"`
}

type Photo struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Mime   string `json:"mime"`
}

type Body struct {
	Text  string `json:"text"`
	Photo *Photo `json:"photo,omitempty"`
}

type Message struct {
	MessageId int     `json:"messageId"`
	Body      Body    `json:"body"`
	Read      *string `json:"read,omitempty"`
	Time      string  `json:"time"`
	Sender    User    `json:"sender"`
}

func (db *appdbimpl) GetMessages(conversationId int, viewerId int) ([]Message, error) {
	rows, err := db.c.Query(`
		SELECT
			m.messageId,
			m.text,
			p.url,
			p.mime,
			p.width,
			p.height,
			m.time,
			uu.username AS senderUsername,
			m.sender
		FROM Message m
		LEFT JOIN Photo p ON m.photoId = p.photoId
		JOIN UserUsername uu ON m.sender = uu.userId
		WHERE m.conversationId = ? AND NOT EXISTS (
      SELECT 1
      FROM DeletedMessage d
      WHERE d.messageId = m.messageId
  )
		ORDER BY m.time DESC;
	`, conversationId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	messages := make([]Message, 0)

	for rows.Next() {
		var msg Message
		var text sql.NullString
		var photoURL, photoMime sql.NullString
		var photoWidth, photoHeight sql.NullInt64
		var senderId int

		err := rows.Scan(
			&msg.MessageId,
			&text,
			&photoURL,
			&photoMime,
			&photoWidth,
			&photoHeight,
			&msg.Time,
			&msg.Sender.Username,
			&senderId,
		)
		if err != nil {
			return nil, err
		}

		if text.Valid {
			msg.Body.Text = text.String
		}

		if photoURL.Valid {
			msg.Body.Photo = &Photo{
				Url:    photoURL.String,
				Mime:   photoMime.String,
				Width:  int(photoWidth.Int64),
				Height: int(photoHeight.Int64),
			}
		}

		msg.Sender.UserID = fmt.Sprint(senderId)
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// ---- Gestione Read / Unread e ReadMessage ----
	for i := range messages {
		senderId, err := strconv.Atoi(messages[i].Sender.UserID)
		if err != nil {
			return nil, err
		}

		if senderId != viewerId {
			// Messaggi ricevuti: segna come letti nel DB ma non includere il campo Read nel JSON
			_, err = db.c.Exec(`
				INSERT OR IGNORE INTO ReadMessage(userId, messageId)
				VALUES (?, ?)
			`, viewerId, messages[i].MessageId)
			if err != nil {
				return nil, err
			}
			messages[i].Read = nil
		} else {
			// Messaggi inviati: calcola lo stato Read / Unread
			var exists int
			err = db.c.QueryRow(`
				SELECT 1
				FROM ReadMessage
				WHERE userId != ? AND messageId = ?
			`, viewerId, messages[i].MessageId).Scan(&exists)

			if errors.Is(err, sql.ErrNoRows) {
				status := "unread"
				messages[i].Read = &status
			} else if err != nil {
				return nil, err
			} else {
				status := "read"
				messages[i].Read = &status
			}
		}
	}

	return messages, nil
}
