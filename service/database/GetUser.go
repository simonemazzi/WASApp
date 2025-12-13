package database

import (
	"database/sql"
	"errors"
)

type DBUser struct {
	UserId   string
	Username string
}

func (db *appdbimpl) Users() ([]DBUser, error) {

	rows, err := db.c.Query(`
        SELECT u.userId, uu.username
        FROM User u
        JOIN UserUsername uu
          ON uu.userId = u.userId
        WHERE uu.updateId = (
            SELECT MAX(updateId)
            FROM UserUsername
            WHERE userId = u.userId
        )
    `)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	var result []DBUser

	for rows.Next() {
		var u DBUser
		if err := rows.Scan(&u.UserId, &u.Username); err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (db *appdbimpl) UsersBySearch(token string) ([]DBUser, error) {
	rows, err := db.c.Query(`
        SELECT u.userId, uu.username
        FROM User u
        JOIN UserUsername uu
          ON uu.userId = u.userId
        WHERE uu.updateId = (
            SELECT MAX(updateId)
            FROM UserUsername
            WHERE userId = u.userId
        ) AND uu.username LIKE ?
    `, token)
	if err != nil {
		return []DBUser{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var result []DBUser

	for rows.Next() {
		var u DBUser
		if err := rows.Scan(&u.UserId, &u.Username); err != nil {
			return nil, err
		}
		result = append(result, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (db *appdbimpl) IDExists(userId int) bool {
	var tmp int
	err := db.c.QueryRow(`
		SELECT 1
		FROM User
		WHERE userId = ?
	`, userId).Scan(&tmp)

	return err == nil
}

func (db *appdbimpl) SearchUserByUsername(username string, time string) (int, error) {
	var id int
	err := db.c.QueryRow(`
SELECT uu.userId
FROM UserUsername uu
WHERE uu.username = ?
  AND uu.time <= ?
ORDER BY uu.time DESC
LIMIT 1`, username, time).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, errors.New("user not found")
	}
	return id, nil
}

func (db *appdbimpl) GetUsername(userId int, time string) (string, error) {
	var id string
	err := db.c.QueryRow(`
SELECT uu.username
FROM UserUsername uu
WHERE uu.userId= ?
  AND uu.time <= ?
ORDER BY uu.time DESC
LIMIT 1`, userId, time).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errors.New("user not found")
	}
	return id, nil
}
