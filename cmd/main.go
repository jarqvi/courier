package main

import (
	"github.com/jarqvi/courier/internal/db"
	"github.com/jarqvi/courier/internal/dns"
	"github.com/jarqvi/courier/internal/log"
)

func main() {
	logger, err := log.NewZapLogger()
	if err != nil {
		panic(err)
	}

	logger.Sync()

	logger.Info("logger initialized")

	database, err := db.Connect()
	if err != nil {
		panic(err)
	}

	defer func() {
		err := database.Disconnect()
		if err != nil {
			logger.Error("failed to disconnect from database: ", err)
		}
	}()

	logger.Info("database initialized")

	_, err = dns.Init()
	if err != nil {
		panic(err)
	}

	logger.Info("dns client initialized")
}
