package main

import (
	"log"
	"os"

	postgresdb "github.com/M1123Ananda/tododo/db"
	"github.com/M1123Ananda/tododo/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found (or failed to load), relying on system env vars")
	}

	db, err := postgresdb.Setup(os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("Fail to connect to DB: %v", err)
	}

	err = postgresdb.InitTables()
	if err != nil {
		log.Fatalf("Fail to init tables in DB: %v", err)
	}

	log.Println("Connected to Database: ", db.Name())

	router := gin.Default()
	router.POST("/register", service.RegisterUser)
	router.POST("/login", service.LoginUser)
	router.POST("/todos", service.CreateToDoItem)
	router.PUT("/:id", service.UpdateToDoItem)
	
	router.Run("localhost:6969")
}
