package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"crud/user"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
	})

	router.GET("/users", user.GetUsers)
	router.GET("/users/:id", user.GetUser)
	router.POST("/users", user.CreateUser)
	router.PUT("/users", user.UpdateUser)
	router.DELETE("/users/:id", user.DeleteUser)

	router.Run(":8080")
}
