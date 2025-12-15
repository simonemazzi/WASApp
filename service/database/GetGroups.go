package database

//TODO:METTERE I LOCATION AI POST

type Group struct {
	GroupId      int    `json:"group_id"`
	Name         string `json:"name"`
	Photo        Avatar `json:"photo"`
	Participants []User `json:"participants"`
}

func (db *appdbimpl) GetGroups(userId int) ([]Group, error) {
	groups := make([]Group, 0)

	// 1) Prendo i gruppi in cui l'utente è attualmente presente
	rows, err := db.c.Query(`
		SELECT DISTINCT g.groupId, g.name
		FROM Group_ g
		JOIN Components c ON c.groupId = g.groupId
		WHERE c.userId = ? AND c.timeLeft IS NULL
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var g Group
		err := rows.Scan(&g.GroupId, &g.Name)
		if err != nil {
			return nil, err
		}

		// 2) Foto del gruppo (ultima)
		var avatar Avatar
		err = db.c.QueryRow(`
			SELECT p.url, p.mime, p.width, p.height
			FROM GroupPhoto gp
			JOIN Photo p ON p.photoId = gp.photoId
			WHERE gp.groupId = ?
			ORDER BY gp.updateId DESC
			LIMIT 1
		`, g.GroupId).Scan(
			&avatar.Url,
			&avatar.Mime,
			&avatar.Width,
			&avatar.Height,
		)
		if err == nil {
			g.Photo = avatar
		}

		// 3) Partecipanti del gruppo
		urows, err := db.c.Query(`
			SELECT u.userId, uu.username
			FROM Components c
			JOIN User u ON u.userId = c.userId
			JOIN UserUsername uu ON uu.userId = u.userId
			WHERE c.groupId = ?
			  AND c.timeLeft IS NULL
			  AND uu.updateId = (
			      SELECT MAX(updateId)
			      FROM UserUsername
			      WHERE userId = u.userId
			  )
		`, g.GroupId)
		if err != nil {
			return nil, err
		}

		var users []User
		for urows.Next() {
			var u User
			err := urows.Scan(&u.UserID, &u.Username)
			if err != nil {
				urows.Close()
				return nil, err
			}
			users = append(users, u)
		}
		urows.Close()
		if err := urows.Err(); err != nil {
			return nil, err
		}

		g.Participants = users
		groups = append(groups, g)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}
