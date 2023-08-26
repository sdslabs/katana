package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sdslabs/katana/configs"
	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/utils"
)

var db *sql.DB

func setup() error {
	for i := 0; i < 10; i++ {
		log.Printf("Trying to connect to MySQL, attempt %d", i+1)
		database, err := sql.Open("mysql", configs.MySQLConfig.Username+":"+configs.MySQLConfig.Password+"@tcp("+utils.GetKatanaLoadbalancer()+":3306)/mysql")
		if err != nil {
			return fmt.Errorf("cannot connect to mysql: %w", err)
		}
		db = database
		log.Println("Connecting to MySQL")
		err = db.Ping()
		if err != nil {
			log.Println("MySQL connection was not established")
			log.Println("Error: ", err)
			time.Sleep(time.Duration(g.KatanaConfig.TimeOut) * time.Second)
		} else {
			log.Println("MySQL Connection Established")
			if err := setupGogs(); err != nil {
				return fmt.Errorf("cannot setup gogs: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("cannot connect to mysql")
}

func setupGogs() error {
	if err := CreateDatabase(gogsDatabase); err != nil {
		fmt.Errorf("cannot create database: %w", err)
	}
	if err := CreateGogsAdmin(configs.AdminConfig.Username, configs.AdminConfig.Password); err != nil {
		fmt.Errorf("cannot create gogs admin: %w", err)
	}
	if err := CreateAccessToken(configs.AdminConfig.Username, configs.AdminConfig.Password); err != nil {
		fmt.Errorf("cannot create access token: %w", err)
	}
	return nil
}

func Init() error {
	if err := setup(); err != nil {
		return fmt.Errorf("cannot setup: %w", err)
	}
	return nil
}
