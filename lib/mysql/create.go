package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func CreateDatabase(db *sql.DB, database string) error {
	_, err := db.Exec("CREATE DATABASE " + database)
	if err != nil {
		return err
	}
	return nil
}

func CreateGogsUser(db *sql.DB, username, password string) error {
	_, err := db.Exec("CREATE USER '" + username + "'@'localhost' IDENTIFIED BY '" + password + "'")
	if err != nil {
		return err
	}
	return nil
}

func CreateGogsUsers(db *sql.DB, users map[string]string) error {
	for username, password := range users {
		err := CreateGogsUser(db, username, password)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateGogsWebhook(db *sql.DB, database, webhook string) error {
	_, err := db.Exec("GRANT ALL PRIVILEGES ON " + database + ".* TO '" + webhook + "'@'localhost'")
	if err != nil {
		return err
	}
	return nil
}
