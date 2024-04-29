package server

import (
	"database/sql"
	"log"

	"github.com/spf13/viper"
)

func InitDatabase(config *viper.Viper) *sql.DB {
	connectionString := config.GetString("database.connection_string")
	maxIdleConnection := config.GetInt("database.max_idle_connections")
	maxOpenConnection := config.GetInt("database.max_open_connections")
	connectionMaxLifetime := config.GetDuration("database.connection_max_lifetime")
	driverName := config.GetString("database.driver_name")

	if connectionString == "" {
		log.Fatalf("Database connection string is missing")
	}

	dbHandler, err := sql.Open(driverName, connectionString)
	if err != nil {
		log.Fatalf("Error while initializing database: %v", err)
	}
	dbHandler.SetMaxIdleConns(maxIdleConnection)
	dbHandler.SetMaxOpenConns(maxOpenConnection)
	dbHandler.SetConnMaxLifetime(connectionMaxLifetime)

	err = dbHandler.Ping()
	if err != nil {
		dbHandler.Close()
		log.Fatalf("Error while initializing database: %v", err)
	}

	return dbHandler
}
