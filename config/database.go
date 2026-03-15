package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgres() {

	// ======================
	// DETECT RAILWAY ENV
	// ======================

	host := os.Getenv("PGHOST")
	port := os.Getenv("PGPORT")
	user := os.Getenv("PGUSER")
	password := os.Getenv("PGPASSWORD")
	dbname := os.Getenv("PGDATABASE")

	sslMode := "require"

	// ======================
	// FALLBACK TO LOCAL ENV
	// ======================

	if host == "" {

		host = Cfg.DBHost
		port = Cfg.DBPort
		user = Cfg.DBUser
		password = Cfg.DBPass
		dbname = Cfg.DBName

		if Cfg.DBSSL {
			sslMode = "require"
		} else {
			sslMode = "disable"
		}
	}

	// ======================
	// BUILD DSN
	// ======================

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		host,
		user,
		password,
		dbname,
		port,
		sslMode,
	)

	// ======================
	// CONNECT DATABASE
	// ======================

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("sql db instance error: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("database unreachable: %v", err)
	}

	// ======================
	// CONNECTION POOL
	// ======================

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	log.Println("PostgreSQL connected successfully")
}
