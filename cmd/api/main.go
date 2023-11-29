package main

import (
	"log"
	"visa-hunter/internal/server"
)

func main() {
	log.Fatal(server.Start(":8080"))
}
