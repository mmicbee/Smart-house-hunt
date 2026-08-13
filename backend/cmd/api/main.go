package main

import (
        "log"
        "net/http"

        "github.com/gin-gonic/gin"
)

func main() {
        r := gin.Default()

        // Health check
        r.GET("/health", func(c *gin.Context) {
                c.JSON(http.StatusOK, gin.H{
                        "status":  "ok",
                        "message": "Smart House Hunt API is running",
                })
        })

        // TODO: Add routes for auth, properties, etc.

        log.Println("Server starting on :8080")
        if err := r.Run(":8080"); err != nil {
                log.Fatal(err)
        }
}
