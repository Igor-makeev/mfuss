package main

import (
	"context"
	"fmt"
	"log"
	"mfuss/configs"
	"mfuss/internal/handler"
	"mfuss/internal/repositories"
	"mfuss/internal/server"
	"mfuss/internal/service"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {

	logrus.Printf("Build version: %v", buildVersion)
	logrus.Printf("Build date: %v", buildDate)
	logrus.Printf("Build commit: %v", buildCommit)

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
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	srv := new(server.Server)

	serverErrChan := srv.Run(cfg, handler)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-signals:

		fmt.Println("main: got shutdown signal. Shutting down...")
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
