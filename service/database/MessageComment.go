package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) MessageComment(messageId int, commentId int) (bool, error) {
	var exists int
	err := db.c.QueryRow(`SELECT 1 FROM Comment c WHERE c.messageId = ? AND c.commentId = ?`, messageId, commentId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (db *appdbimpl) CommentUser(commentId int, userId int) (bool, error) {
	var exists int
	err := db.c.QueryRow(`SELECT 1 FROM Comment c WHERE c.commentId = ? AND c.userId = ?`, commentId, userId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
