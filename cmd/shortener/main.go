package main

import (
	"context"
	"mfuss/configs"
	"mfuss/internal/handler"
	"mfuss/internal/repositories"
	"mfuss/internal/server"
	"mfuss/internal/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {

	cfg := configs.NewConfig()
	var urlstorage repositories.URLStorager
	var err error

	if cfg.DBDSN == "" {
		urlstorage, err = PrepareMemoryStorage(cfg)
		if err != nil {
			logrus.Fatal(err)
		}
	} else {
		conn, err := repositories.NewPostgresClient(cfg)
		if err != nil {
			logrus.Fatal(err)
		}
		urlstorage = repositories.NewPostgresStorage(cfg, conn)
	}

	rep, err := repositories.NewRepository(cfg, urlstorage)
	if err != nil {
		logrus.Fatal(err)
	}
	service := service.NewService(rep)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	service.Queue.Run(ctx, service.MarkAsDeleted)

	handler := handler.NewHandler(service)

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

	if service.Close(ctx); err != nil {
		logrus.Errorf("error occured on closing service: %s", err.Error())
	}
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shuting down: %s", err.Error())
	}

}

func PrepareMemoryStorage(cfg *configs.Config) (*repositories.MemoryStorage, error) {

	dump, err := repositories.NewDump(cfg.FileStoragePath)
	if err != nil {
		return nil, err
	}

	ms := repositories.NewMemoryStorage(cfg, dump)
	if err := ms.LoadFromDump(); err != nil {
		return nil, err
	}

	return ms, nil

}
