package db

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/kpango/glg"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func ConnectDB() *gorm.DB {
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")

	connStr := fmt.Sprintf("host=%s user=%s password=%s sslmode=disable", dbHost, dbUser, dbPass)

	DB, err := gorm.Open("postgres", connStr)

	if err != nil {
		glg.Errorf("Failed to connect to database: %v", err)
	}

	var count []uint32
	DB.Raw("SELECT count(*) FROM pg_database WHERE datname  = ?", dbName).Scan(&count)
	if count[0] == 0 {
		sql := fmt.Sprintf("CREATE DATABASE %s", dbName)
		result := DB.Exec(sql)

		if result.Error != nil {
			glg.Errorf("Failed create database: %v", result.Error)
		}

	}
	conn_db_url := fmt.Sprintf("%s dbname=%s", connStr, dbName)

	DB, err = gorm.Open("postgres", conn_db_url)
	if err != nil {
		glg.Errorf("Failed to connect to database: %v", err)
	}

	DB.AutoMigrate(&Person{})
	return DB
}
