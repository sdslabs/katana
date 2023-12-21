package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/utils"
)

var db *sql.DB

func setup() error {
	for i := 0; i < 10; i++ {
		log.Printf("Trying to connect to MySQL, attempt %d", i+1)
		katanaLB, err := utils.GetKatanaLoadbalancer()
		if err != nil {
			return err
		}
		database, err := sql.Open("mysql", g.MySQLConfig.Username+":"+g.MySQLConfig.Password+"@tcp("+katanaLB+":3306)/mysql")
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

			if err := CreateDatabase(gogsDatabase); err != nil {
				if !strings.Contains(err.Error(), "database exists") {
					return fmt.Errorf("cannot create gogs database: %w", err)
				} else {
					if err := setupGogs(); err != nil {
						return fmt.Errorf("cannot setup gogs: %w", err)
					}else{
						log.Println("Gogs MySQL database setup successfully")
					}
				}
			} else {
				log.Println("Gogs database created successfully in MySQL")
			}
			
			return nil
		}
	}
	return fmt.Errorf("cannot connect to mysql")
}

func setupGogs() error {
	if err := CreateGogsAdmin(g.AdminConfig.Username, g.AdminConfig.Password); err != nil {
		return fmt.Errorf("cannot create gogs admin: %w", err)
	}
	if err := CreateAccessToken(g.AdminConfig.Username, g.AdminConfig.Password); err != nil {
		return fmt.Errorf("cannot create access token: %w", err)
	}
	return nil
}

func Init() error {
	if err := setup(); err != nil {
		return fmt.Errorf("cannot setup: %w", err)
	}
	return nil
}
