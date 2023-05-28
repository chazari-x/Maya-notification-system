package main

import (
	"log"

	"Maya-notification-system/rabbitmq/internal/app/server"
)

func main() {
	log.Print(server.StartServer())
}
