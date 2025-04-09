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
	err = r.RunTLS(":443", "/etc/ssl/certs/selfsigned.crt", "/etc/ssl/private/selfsigned.key")
	if err != nil {
		log.Fatal(err)
	}
}
