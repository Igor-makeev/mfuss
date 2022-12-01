package app

import (
	"context"
	"mfuss/internal/handler"
	"mfuss/internal/repositories"
	"mfuss/internal/server"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

type App struct {
}

func NewApp() *App {
	app := &App{}
	return app
}

func (app *App) Run() error {

	storage := repositories.NewMemoryStorage()
	repository := repositories.NewRepository(storage)
	handler := handler.NewHandler(repository)

	srv := new(server.URLserver)
	go func() {

		if err := srv.Run(handler.InitRoutes()); err != nil {
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
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return srv.Shutdown(ctx)
}
