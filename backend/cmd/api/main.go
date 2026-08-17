package main

import (
    "database/sql"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
  db, err := database.Connect()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatal("failed to run migrations:", err)
    }

    r := gin.Default()

    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status":  "ok",
            "message": "Smart House Hunt API is running",
        })
    })

    log.Println("Database: SQLite 3")
    log.Println("Migrations: completed")
    log.Println("Server starting on :8080")

    if err := r.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}
