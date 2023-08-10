package mysql

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sdslabs/katana/configs"
	g "github.com/sdslabs/katana/configs"
)

var db *sql.DB

func setup(LoadbalancerIP string) {
	database, err := sql.Open("mysql", configs.MySQLConfig.Username+":"+configs.MySQLConfig.Password+"@tcp("+LoadbalancerIP+":3306)/mysql")
	if err != nil {
		panic(err.Error())
	}
	db = database
	log.Println("Connecting to MySQL")
	err = db.Ping()
	if err != nil {
		log.Println("MySQL connection was not established")
		log.Println("Error: ", err)
		time.Sleep(time.Duration(g.KatanaConfig.TimeOut) * time.Second)
		setup(LoadbalancerIP)
	} else {
		log.Println("MySQL Connection Established")
		setupGogs()
	}
}

func setupGogs() {
	CreateDatabase(gogsDatabase)
	CreateGogsAdmin(configs.AdminConfig.Username, configs.AdminConfig.Password)
	CreateAccessToken(configs.AdminConfig.Username, configs.AdminConfig.Password)
}

func Init(LoadbalancerIP string) {
	go setup(LoadbalancerIP)
}
