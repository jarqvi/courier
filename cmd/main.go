package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jarqvi/courier/internal/db"
	"github.com/jarqvi/courier/internal/log"
	"github.com/jarqvi/courier/internal/smtp"
)

func main() {
	err := log.NewZapLogger()
	if err != nil {
		panic(err)
	}

	log.Logger.Info("logger initialized")

	err = db.Connect()
	if err != nil {
		panic(err)
	}

	smtp.Init()

	if smtp.ServerError != nil {
		panic(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	smtp.Shutdown()
	db.Instance.Disconnect()
	log.Logger.Sync()
}
