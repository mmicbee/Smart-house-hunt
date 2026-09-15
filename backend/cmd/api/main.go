package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"smart-house-hunt/config"
	"smart-house-hunt/database"
	"smart-house-hunt/routes"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatal("failed to run migrations:", err)
	}

	router := gin.Default()
	routes.Setup(router, db, cfg)

	log.Println("Database: SQLite 3")
	log.Println("Migrations: completed")
	log.Println("Server starting on", cfg.Port)

	if err := router.Run(cfg.Port); err != nil {
		log.Fatal(err)
	}
}