package main

import (
	"log"
	"project/user-management/internal/repository"
	"project/user-management/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	db, err := repository.InitDB()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	r := gin.Default()

	routes.InitRoutes(r, db)

	r.Run(":8080")
}
