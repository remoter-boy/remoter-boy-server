package common

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/lib/pq" // PostgreSQL 드라이버
)

func InitDatabase() *sql.DB {
	config := GetConfig()

	dbPort, err := strconv.Atoi(config.DatabasePort)

	if err != nil {
		return nil
	}

	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.DatabaseHost, dbPort, config.DatabaseUser, config.DatabasePass, config.DatabaseName)

	db, err := sql.Open("postgres", dbInfo)

	if err != nil {
		log.Println(err)
		return nil
	}

	return db
}
