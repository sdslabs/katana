package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sdslabs/katana/configs"
)

var db *sql.DB

func setup() {
	database, err := sql.Open("mysql", configs.MySQLConfig.Username+":"+configs.MySQLConfig.Password+"@tcp("+configs.ServicesConfig.ChallengeDeployer.Host+":"+configs.MySQLConfig.Port+")/mysql")
	if err != nil {
		panic(err.Error())
	}
	db = database
	fmt.Println("Connecting to MySQL")
	err = db.Ping()
	if err != nil {
		log.Println("MySQL connection was not established")
		log.Println("Error: ", err)
		time.Sleep(5 * time.Second)
		setup()
	} else {
		log.Println("MySQL Connection Established")
		setupGogs()
	}
}

func setupGogs() {
	CreateDatabase(gogsDatabase)
}

func Init() {
	go setup()
}
