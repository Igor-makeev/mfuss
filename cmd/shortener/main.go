package main

import (
	"context"
	"log"
	"mfuss/configs"
	"mfuss/internal/handler"
	"mfuss/internal/repositories"
	"mfuss/internal/server"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

func main() {

	var cfg configs.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal("failed to parse config environment variables")
	}

	storage := repositories.NewMemoryStorage()

	handler := handler.NewHandler(storage, cfg)

	srv := server.NewURLServer(handler.Router, cfg)

	go func() {

		if err := srv.ListenAndServe(); err != nil {
			logrus.Fatalf("failed to listen and serve: %+v", err.Error())
		}
	}()

	logrus.Print("shortener started...")

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logrus.Print("shortener shuting down.")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shuting down: %s", err.Error())
	}

}
