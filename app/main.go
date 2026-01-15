package main

import (
	"auth-service/src/db"
	"auth-service/src/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db.Init()

	r := gin.Default()

	r.GET("/health", handlers.Health)
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.GET("/validate", handlers.JWTAuthMiddleware(), handlers.Validate)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(r.Run(":" + port))
}
