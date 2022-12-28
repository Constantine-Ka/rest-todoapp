package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"todo"
	"todo/pkg/handler"
	"todo/pkg/repository"
	"todo/pkg/service"
)

// @title Todo App RestAPI
// @version 0.1
// @description RestAPI Server for TodoList Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @contact.name   Constantine Ka

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		log.Fatalf("error initialization configs:%s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables:%s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host1"),
		Port:     viper.GetString("db.port1"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db:%s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(todo.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server:%s", err.Error())
		}
	}()
	logrus.Print("TodoApp Started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Print("TodoApp Shutting Down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down:%s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db-connection shutting down:%s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
