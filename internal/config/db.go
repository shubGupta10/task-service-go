package config

import (
	"fmt"
	"log"
	"os"
	"testapp/internal/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	var db *gorm.DB
	var err error

	// Try connecting 10 times with 2-second intervals
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		db.AutoMigrate(&models.User{}, &models.Todo{})
		if err == nil {
			break
		}

		log.Println("Failed to connect to the database. Retrying in 2 seconds...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Could not connect to the database after retries: ", err)
	}

	DB = db
	fmt.Println("Database connected successfully")
}
