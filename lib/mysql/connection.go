package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:<yourMySQLdatabasepassword>@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	fmt.Println("Success!")
}
