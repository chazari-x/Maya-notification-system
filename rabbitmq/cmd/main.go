package main

import (
	"log"

	"rabbitmq/internal/app/server"
)

func main() {
	log.Print(server.StartServer())
}
