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

// Message rappresenta un messaggio o un inoltro
type Message struct {
	MessageId   int     `json:"messageId"`
	Body        Body    `json:"body"`
	Read        *string `json:"read,omitempty"`
	Time        string  `json:"time"`
	Sender      User    `json:"sender"`
	IsForwarded bool    `json:"isForwarded"`
}

func (db *appdbimpl) GetMessages(conversationId int, viewerId int) ([]Message, error) {
	rows, err := db.c.Query(`
		SELECT
			m.messageId,
			m.time,
			m.sender,
			uu.username,
			CASE WHEN m.originalMessage IS NOT NULL THEN 1 ELSE 0 END AS isForwarded
		FROM Message m
		JOIN UserUsername uu ON uu.userId = m.sender
		WHERE m.conversationId = ?
		  AND NOT EXISTS (
		      SELECT 1
		      FROM DeletedMessage d
		      WHERE d.messageId = m.messageId
		  )
		ORDER BY m.time DESC
	`, conversationId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	messages := make([]Message, 0)

	for rows.Next() {
		var messageId int
		var time string
		var senderId int
		var senderUsername string
		var isForwarded int

		err := rows.Scan(
			&messageId,
			&time,
			&senderId,
			&senderUsername,
			&isForwarded,
		)
		if err != nil {
			return nil, err
		}

		// Ottieni il contenuto originale (testo, foto, mittente) anche per inoltri multipli
		originalMsg, err := db.OriginalMessageInfo(messageId)
		if err != nil {
			return nil, err
		}

		// Popola i dati principali con le informazioni del messaggio visualizzato
		msg := originalMsg
		msg.MessageId = messageId // mantiene l'ID del messaggio corrente
		msg.Time = time
		msg.Sender = User{
			UserID:   fmt.Sprint(senderId),
			Username: senderUsername,
		}
		msg.IsForwarded = isForwarded == 1 // indica se il messaggio è un inoltro

		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// ---- Gestione Read / Unread ----
	for i := range messages {
		senderId, err := strconv.Atoi(messages[i].Sender.UserID)
		if err != nil {
			return nil, err
		}

		if senderId != viewerId {
			// Messaggi ricevuti → segna come letti
			_, err = db.c.Exec(`
				INSERT OR IGNORE INTO ReadMessage(userId, messageId)
				VALUES (?, ?)
			`, viewerId, messages[i].MessageId)
			if err != nil {
				return nil, err
			}
			messages[i].Read = nil
		} else {
			// Messaggi inviati → stato read/unread
			var exists int
			err = db.c.QueryRow(`
				SELECT 1
				FROM ReadMessage
				WHERE userId != ? AND messageId = ?
			`, viewerId, messages[i].MessageId).Scan(&exists)

			switch {
			case errors.Is(err, sql.ErrNoRows):
				status := "unread"
				messages[i].Read = &status
			case err != nil:
				return nil, err
			default:
				status := "read"
				messages[i].Read = &status
			}
		}
	}

	return messages, nil
}
