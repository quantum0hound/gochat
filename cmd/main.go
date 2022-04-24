package main

import (
	"github.com/joho/godotenv"
	"github.com/quantum0hound/gochat/configs"
	"github.com/quantum0hound/gochat/internal/handler"
	"github.com/quantum0hound/gochat/internal/handler/server/http"
	"github.com/quantum0hound/gochat/internal/handler/server/ws"
	"github.com/quantum0hound/gochat/internal/repository"
	"github.com/quantum0hound/gochat/internal/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {
	//logrus.SetReportCaller(true)
	//logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.SetLevel(logrus.DebugLevel)
	if err := configs.InitConfig(); err != nil {
		logrus.Fatalf("Error occured, during reading the config file: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error occured during reading .env file :%s", err.Error())
	}

	db, err := repository.NewPostgresDb(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetUint("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   viper.GetString("db.dbname"),
		SslMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("Error occured during connection to DB :%s", err.Error())
	}

	wsServer := ws.NewWebSocketServer()
	repo := repository.NewRepository(db)
	srv := service.NewService(repo)
	hnd := handler.NewHandler(srv, wsServer)

	go wsServer.Run()
	server := new(http.Server)
	if err := server.Run(viper.GetUint("port"), hnd.InitRoutes()); err != nil {
		logrus.Fatalf("Error occured during HTTP server execution :%s", err.Error())
	}
}
