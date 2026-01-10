package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("Failed to connect to database!")
	}

	// ========== Connection Pool Settings ==========
	// Important for production performance
	sqlDB, err := DB.DB()
	if err != nil {
		panic("Failed to get database connection pool!")
	}

	// SetMaxIdleConns: จำนวน idle connections ที่เก็บไว้ใน pool
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns: จำนวน connections สูงสุดที่เปิดพร้อมกันได้
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime: อายุสูงสุดของ connection (ป้องกัน stale connections)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// SetConnMaxIdleTime: เวลาสูงสุดที่ connection จะ idle ก่อนถูกปิด
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	// fmt.Println("Database connection successfully opened")
}
