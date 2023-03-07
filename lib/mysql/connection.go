package mysql

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func setup() {
	database, err := sql.Open("mysql", "root:<yourMySQLdatabasepassword>@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err.Error())
	}
	db = database
	err = db.Ping()
	if err != nil {
		log.Println("MySQL connection was not established")
		time.Sleep(5 * time.Second)
		setup()
	} else {
		log.Println("MySQL Connection Established")
		setupGogs()
	}
}

func setupGogs() {
	CreateDatabase(db, gogsDatabase)
}

func Init() {
	go setup()
}
