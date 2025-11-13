package app

import (
	"context"
	"os"
	"os/signal"
	"review-rotator/internal/handler"
	"review-rotator/internal/repository"
	"review-rotator/internal/service"
	"review-rotator/pkg"
	"syscall"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)


type App struct {
	server *pkg.Server
}

func NewApp() *App {
	return & App{
		server: new(pkg.Server),
	}
}

func (a *App) Run() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := pkg.InitConfig(); err != nil {
		logrus.Fatal(err.Error())
	}
	if err := pkg.LoadEnv(); err != nil {
		logrus.Fatalf("error of load env %s", err.Error())
	}
	db, err := pkg.NewPostgresDB(pkg.GetDBconfig())
	if err != nil {
		logrus.Fatal(err.Error())
	}


	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := handler.NewHandler(service)

	go func () {
		if err := a.server.Run(pkg.GetPort(), handler.InitRouts()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()
	logrus.Println("Service started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<- quit

	logrus.Println("server shutting down")
	if err := a.server.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}