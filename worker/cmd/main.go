package main

import (
	"log"

	"Maya-notification-system/worker/internal/app/worker"
)

func main() {
	log.Print("worker closed err: ", worker.StartWorker())
}
