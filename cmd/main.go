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

	config.CleanInappropriateDB()

	r := routes.SetupRouter()
	err = r.RunTLS(":443", "/etc/ssl/certs/selfsigned.crt", "/etc/ssl/private/selfsigned.key")
	if err != nil {
		log.Fatal(err)
	}
}
