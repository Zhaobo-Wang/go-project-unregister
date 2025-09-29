package main

import (
	"log"
	"os"

	"github.com/Zhaobo-Wang/go-projects/database"
	"github.com/Zhaobo-Wang/go-projects/routes"
	"github.com/joho/godotenv"
)

func main() {
	f, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.SetOutput(f)

	err = godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	database.ConnectDatabase()

	r := routes.SetupRouter()

	r.Run(":8081")
}
