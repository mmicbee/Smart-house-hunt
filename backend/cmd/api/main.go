package main

import (
    "database/sql"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, err := sql.Open("sqlite3", "./smart-house-hunt.db")
    if err != nil {
        log.Fatal("failed to open database:", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatal("failed to connect to database:", err)
    }

    r := gin.Default()

    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status":  "ok",
            "message": "Smart House Hunt API is running",
        })
    })

    log.Println("Database: SQLite 3")
    log.Println("Server starting on :8080")

    if err := r.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}
