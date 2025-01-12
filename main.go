package main

import (
	"log"
	"runner/config"
)

func main() {
	log.Println("Starting Runners App")
	log.Println("Initializing configuration")
	config := config.InitConfig("runners")
    log.Println(“Initializing database”)
    dbHandler := server.InitDatabase(config)
	log.Println("Initializing HTTP server")
	httpServer := server.InitHTTPServer(config, dbHandler)
	httpServer.Start()
}
