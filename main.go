package main

import (
	"log"
	"runners/config"
	"runners/server"

	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting Runners App")

	log.Println("Initializing configuration")
	config := config.InitConfig("config_runners.toml")

	log.Println("Initializing database")
	dbHandler := server.InitDatabase(config)

	log.Println("Initializing HTTP server")
	httpServer := server.InitHttpServer(config, dbHandler)
	httpServer.Start()
}
