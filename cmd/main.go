package main

import (
	"log"
	testapi "test-api"
	"test-api/pkg/handler"
	"test-api/pkg/repository"
	"test-api/pkg/service"

	"github.com/spf13/viper"
)

func main() {

	if err := initConfig(); err != nil {
		log.Fatalf("error while reading config, %s", err.Error())
	}

	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(testapi.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error runnig server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
