package main

import (
	"context"
	"fmt"
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

	srv := new(server.Server)

	serverErrChan := srv.Run(cfg, handler)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-signals:

		fmt.Println("main: got terminate signal. Shutting down...")
		if service.Close(ctx); err != nil {
			logrus.Errorf("error occured on closing service: %s", err.Error())
		}
		if err := srv.Shutdown(); err != nil {
			fmt.Printf("main: received an error while shutting down the server: %v", err)
		}

	case <-serverErrChan:
		if service.Close(ctx); err != nil {
			logrus.Errorf("error occured on closing service: %s", err.Error())
		}
		fmt.Println("main: got server err signal. Shutting down...")

	}
}
