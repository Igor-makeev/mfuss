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

	rep, err := repositories.NewRepository(cfg)
	if err != nil {
		logrus.Fatal(err)
	}
	service := service.NewService(rep)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	service.Queue.Run(ctx)

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
