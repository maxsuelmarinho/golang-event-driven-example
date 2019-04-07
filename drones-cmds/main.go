package main

import (
	"fmt"
	"os"

	service "github.com/maxsuelmarinho/golang-event-driven-example/drones-cmds/service"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	server := service.NewServer()
	server.Run(fmt.Sprintf(":%s", port))
}
