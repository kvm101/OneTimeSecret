package main

import (
	"log"
	"one_time_secret/internal/model"
	"one_time_secret/routes"
)

func main() {
	err := model.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	r := routes.SetupRouter()
	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
