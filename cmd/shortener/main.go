package main

import (
	"context"
	"mfuss/configs"
	"mfuss/internal/handler"
	"mfuss/internal/repositories"
	"mfuss/internal/server"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {

	cfg := configs.NewConfig()
	var urlstorage repositories.URLStorager

	if cfg.DBDSN == "" {
		dump, err := repositories.NewDump(cfg.FileStoragePath)
		if err != nil {
			logrus.Fatal(err)
		}
		urlstorage, err = repositories.NewMemoryStorage(cfg, dump)
		if err != nil {
			logrus.Fatal(err)
		}
	} else {
		var err error
		urlstorage, err = repositories.NewPostgresStorage(cfg)
		if err != nil {
			logrus.Fatal(err)
		}
	}

	rep, err := repositories.NewRepository(cfg, urlstorage)
	if err != nil {
		logrus.Fatal(err)
	}

	handler := handler.NewHandler(rep)

	srv := server.NewURLServer(handler)

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

	if rep.Close(); err != nil {
		logrus.Errorf("error occured on closing storage: %s", err.Error())
	}
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shuting down: %s", err.Error())
	}

}
