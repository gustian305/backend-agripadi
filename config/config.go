package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string
	JWT string
	GROQAPIKEY string

	DBURL string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
	DBSSL  bool
}

var Cfg Config

func LoadConfig(envFilePath ...string) {

	env := ".env"

	if len(envFilePath) > 0 {
		env = envFilePath[0]
	}

	err := godotenv.Load(env)
	if err != nil {
		log.Println("No .env file found, using system environment")
	}

	dbSSL := false
	if os.Getenv("DB_SSLMODE") == "enable" {
		dbSSL = true
	}

	Cfg = Config{
		AppEnv:  os.Getenv("APP_ENV"),
		AppPort: os.Getenv("APP_PORT"),
		JWT:     os.Getenv("JWT"),
		GROQAPIKEY: os.Getenv("GROQ_API_KEY"),

		DBURL: os.Getenv("DATABASE_URL"),
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),
		DBSSL:  dbSSL,
	}

	log.Println("Config Loaded")
}