package main

import (
	"database/sql"
	"fmt"
)

type user struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Eusr     string `json:"eUsr"`
}

func (u *user) getUser(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT username, eUsr FROM users WHERE id=%d", u.ID)
	return db.QueryRow(statement).Scan(&u.Username, &u.Eusr)
}

func (u *user) updateUser(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE users SET username='%s', eUsr=%s WHERE id=%d", u.Username, u.Eusr, u.ID)
	_, err := db.Exec(statement)
	return err
}

func (u *user) deleteUser(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM users WHERE id=%d", u.ID)
	_, err := db.Exec(statement)
	return err
}

func (u *user) createUser(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO users (username, eUsr) VALUES ('%s', %s)", u.Username, u.Eusr)
	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

func getUsers(db *sql.DB, start, count int) ([]user, error) {
	statement := fmt.Sprintf("SELECT id, username, eUsr FROM users LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []user{}

	for rows.Next() {
		var u user
		if err := rows.Scan(&u.ID, &u.Username, &u.Eusr); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
