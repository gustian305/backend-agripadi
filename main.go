package main

import (
	seed "backend/seeding"
	"backend/config"
	"backend/internal/routes"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// ======================
	// LOAD CONFIG
	// ======================
	config.LoadConfig()

	// ======================
	// CONNECT DATABASE
	// ======================
	config.ConnectPostgres()

	err := config.AutoMigrate(config.DB)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database migration completed")

	// ======================
	// SEED DATA
	// ======================
	err = seed.SeedPestisida()
	if err != nil {
		log.Println("Seed warning:", err)
	} else {
		log.Println("Pestisida seeding completed")
	}

	// ======================
	// ROUTER
	// ======================
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))

	api := router.Group("/api")

	// ======================
	// STATIC FILES
	// ======================
	api.Static("/uploads", "./uploads")

	// ======================
	// ROUTES
	// ======================
	routes.RegisterRoutes(api)
	routes.LoginRoutes(api)

	// protected
	routes.ConversationRoutes(api)
	routes.ChatRoutes(api)
	routes.ProfileRoutes(api)

	// ======================
	// PORT
	// ======================
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port:", port)

	log.Fatal(router.Run(":" + port))
}