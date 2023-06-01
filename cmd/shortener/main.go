package main

import (
	"context"
	"fmt"
	"log"
	"mfuss/configs"
	shortener "mfuss/internal/grpc"
	"mfuss/internal/grpc/auth"
	"mfuss/internal/grpc/interceptors"
	"mfuss/internal/handler"
	"mfuss/internal/repositories"
	"mfuss/internal/server"
	"mfuss/internal/service"
	pb "mfuss/proto"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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
	authenticator := auth.New(cfg)
	go func() {
		ln, grpcErr := net.Listen("tcp", cfg.GRPCAdress)
		if grpcErr != nil {
			log.Printf("listen grpc port error: %s\n", grpcErr)
			return
		}

		i := interceptors.New(authenticator)
		g := grpc.NewServer(grpc.UnaryInterceptor(i.AuthUnaryInterceptor))
		pb.RegisterShortenerServer(g, shortener.NewGRPCServer(service))

		if grpcErr = g.Serve(ln); grpcErr != nil {
			log.Printf("gRPC server error: %s\n", grpcErr)
			return
		}
	}()
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
