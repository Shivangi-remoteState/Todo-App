package main

import (
	"log"
	"mytodoApp/database"
	"mytodoApp/server"
)

//import "github.com/gin-gonic/gin"

func main() {
	log.Println("Server is starting")

	err := database.ConnectAndMigrate("localhost", "5432", "todoApp", "local", "local", database.SSLMODEDisables)
	if err != nil {
		log.Fatal(err)
	}
	server.Start()
}
