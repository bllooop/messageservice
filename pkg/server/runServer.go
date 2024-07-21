package messageservice

import (
	"context"
	"messageservice"
	"messageservice/pkg/handlers"
	"messageservice/pkg/repository"

	//	messageservice "messageservice/pkg/server"
	"messageservice/pkg/service"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func Run() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).Level(zerolog.TraceLevel)
	if err := initConfig(); err != nil {
		logger.Error().Err(err).Msg("")
		logger.Fatal().Msg("There was an error with configs")
	}
	if err := godotenv.Load(); err != nil {
		logger.Error().Err(err).Msg("")
		logger.Fatal().Msg("There was an error with env")
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
		logger.Error().Err(err).Msg("")
		logger.Fatal().Msg("There was an error with database")
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handlers.NewHandler(services)
	serv := new(messageservice.Server)
	/*
		wg := &sync.WaitGroup{}
		wg.Add(1) */
	go func() {
		//	defer wg.Done()
		if err := handlers.CreateTableSQL(); err != nil {
			logger.Error().Err(err).Msg("")
			logger.Fatal().Msg("There was an error with creating a table in database")
		}
	}()
	//wg.Add(1)
	go func() {
		//	defer wg.Done()
		if err := serv.RunServer(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logger.Error().Err(err).Msg("")
			logger.Fatal().Msg("There was an error while running a server")
		}
	}()
	logger.Info().Msg("Server is running")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info().Msg("Server is shutting down")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := serv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("")
		logger.Fatal().Msg("There was an error while shutting down the server")
	}
	if err := db.Close(); err != nil {
		logger.Error().Err(err).Msg("")
		logger.Fatal().Msg("There was an error while closing db connection")
	}
	//wg.Wait()
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
