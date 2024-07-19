package main

import (
	"context"
	"messageservice"
	"messageservice/pkg/handlers"
	"messageservice/pkg/repository"
	"messageservice/pkg/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).Level(zerolog.TraceLevel)

	if err := initConfig(); err != nil {
		logger.Fatal().Msg("There was an error with configs")
		logger.Error().Err(err).Msg("")
	}
	if err := godotenv.Load(); err != nil {
		logger.Fatal().Msg("There was an error with env")
		logger.Error().Err(err).Msg("")
	}
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBname:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logger.Fatal().Msg("There was an error with database")
		logger.Error().Err(err).Msg("")
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handlers.NewHandler(services)
	serv := new(messageservice.Server)
	go func() {
		if err := handlers.CreateTableSQL(); err != nil {
			logger.Fatal().Msg("There was an error with creating a table in database")
			logger.Error().Err(err).Msg("")
		}
	}()

	go func() {
		if err := serv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logger.Fatal().Msg("There was an error while running a server")
			logger.Error().Err(err).Msg("")
		}
	}()
	logger.Info().Msg("Server is running")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logger.Info().Msg("Server is shutting down")
	if err := serv.Shutdown(context.Background()); err != nil {
		logger.Fatal().Msg("There was an error while shutting down the server")
		logger.Error().Err(err).Msg("")
	}
	if err := db.Close(); err != nil {
		logger.Fatal().Msg("There was an error while closing db connection")
		logger.Error().Err(err).Msg("")
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
