package main

import (
	"backend/config"
	"backend/internal/routes"
	"backend/seeding"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {

	config.LoadConfig()
	config.ConnectPostgres()

	err := config.AutoMigrate(config.DB)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Database migration completed")

	// seed pestisida
	err = seed.SeedPestisida()
	if err != nil {
		log.Fatal("Seeding failed:", err)
	}
	log.Println("Pestisida seeding completed")

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api")

	// serve static uploads
	api.Static("/uploads", "./uploads")

	// Public Routes
	routes.RegisterRoutes(api)
	routes.LoginRoutes(api)

	// Protected Routes
	routes.ConversationRoutes(api)
	routes.ChatRoutes(api)
	routes.ProfileRoutes(api)

	router.Run(":8080")
}
