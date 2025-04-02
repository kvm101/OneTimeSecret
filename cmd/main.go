package main

import (
	"log"
	"one_time_secret/config"
	"one_time_secret/routes"
)

func main() {
	err := config.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	r := routes.SetupRouter()
	r.Run(":8999")
}
