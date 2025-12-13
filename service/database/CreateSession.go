package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) CreateSession(username string) (string, string, time.Time, error) {
	var userId int

	// Controlla se l'utente esiste
	query := `SELECT uu.userId
              FROM UserUsername uu
              WHERE uu.username = ?
              ORDER BY uu.updateId DESC
              LIMIT 1`

	err := db.c.QueryRow(query, username).Scan(&userId)
	if errors.Is(err, sql.ErrNoRows) {
		// Se non esiste, crea l'utente
		res, err := db.c.Exec("INSERT INTO User DEFAULT VALUES;")
		if err != nil {
			return "", "", time.Time{}, fmt.Errorf("failed to create user: %w", err)
		}

		newUserId, err := res.LastInsertId()
		if err != nil {
			return "", "", time.Time{}, fmt.Errorf("failed to get new userId: %w", err)
		}

		// Inserisci il nuovo username
		_, err = db.c.Exec(
			"INSERT INTO UserUsername(userId, username) VALUES (?, ?);",
			newUserId, username,
		)
		if err != nil {
			return "", "", time.Time{}, fmt.Errorf("failed to insert username: %w", err)
		}

		userId = int(newUserId)

		err = db.SetMyPhoto("uploads/default/default-avatar-profile-icon-social-600nw-1906669723.png", 600, 600, "image/png", userId)
		if err != nil {
			return "", "", time.Time{}, err
		}
	}

	// Genera token UUID
	token, err := uuid.NewV4()
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("failed to generate token: %w", err)
	}

	// Inserisci login con userId e token
	_, err = db.c.Exec("INSERT INTO Login(userId, loginId) VALUES (?, ?)", userId, token.String())
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("failed to insert login: %w", err)
	}
	// Recupera il timestamp del login appena inserito
	var loginTime time.Time
	err = db.c.QueryRow("SELECT time FROM Login WHERE loginId = ?", token.String()).Scan(&loginTime)
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("failed to query login time: %w", err)
	}

	// Ritorna userId, token e timestamp
	return fmt.Sprintf("%d", userId), token.String(), loginTime, nil
}
