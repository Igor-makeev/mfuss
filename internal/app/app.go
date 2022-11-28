package app

import (
	"context"
	"mfuss/internal/handler"
	"mfuss/internal/server"
	"net/http"
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

	mux := http.NewServeMux()
	handler := handler.NewHandler()

	mux.HandleFunc("/", handler.RootHandler)

	srv := new(server.URLserver)
	go func() {

		if err := srv.Run(mux); err != nil {
			logrus.Fatalf("failed to listen and serve: %+v", err.Error())
		}
	}()

	logrus.Print("shortener started...")

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logrus.Print("shortener shuting down.")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Error("error occured on server shuting down: %s", err.Error())
	}
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return srv.Shutdown(ctx)
}
