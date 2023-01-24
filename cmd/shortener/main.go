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
	queu := repositories.NewQueue()
	rep, err := repositories.NewRepository(cfg, urlstorage, queu)
	if err != nil {
		logrus.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go queu.Listen(ctx, rep.MarkAsDeleted, rep.Queue.UpdateInterval)
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
	cancel()
	logrus.Print("shortener shuting down.")

	if rep.Close(); err != nil {
		logrus.Errorf("error occured on closing storage: %s", err.Error())
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
