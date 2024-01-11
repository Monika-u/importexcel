package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DbClient *sql.DB

func DbConnect() {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "sidha", "pass", "127.0.0.1", 3306, "employeedb")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("impossible to create the connection: %s", err)
	}

	// mysql := mysql.Open("dev:password@123A@tcp(127.0.0.1:3306)/demo")
	// db, err := gorm.Open(mysql, &gorm.Config{})

	if err != nil {
		fmt.Println(err)
		//fmt.Println("Error while connecting to DB")
	}
	DbClient = db
}
