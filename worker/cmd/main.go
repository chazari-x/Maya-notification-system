package main

import (
	"log"

	"worker/internal/app/worker"
)

func main() {
	log.Print("worker closed err: ", worker.StartWorker())
}
