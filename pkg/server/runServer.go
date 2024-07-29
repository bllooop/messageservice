package messageservice

import (

	//	messageservice "messageservice/pkg/server"

	"context"
	"log"
	"messageservice"
	"messageservice/pkg/handlers"
	cons "messageservice/pkg/messaging"
	"messageservice/pkg/repository"
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
	logger := zerolog.New(os.Stdout).Level(zerolog.TraceLevel)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	if err := initConfig(); err != nil {
		logger.Error().Err(err).Msg("")
		logger.Fatal().Msg("There was an error with configs")
	}
	if err := godotenv.Load(); err != nil {
		logger.Error().Err(err).Msg("")
		logger.Fatal().Msg("There was an error with env")
	}
	dbpool, err := repository.NewPostgresDB(repository.Config{
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
	repos := repository.NewRepository(dbpool)
	services := service.NewService(repos)
	handlers := handlers.NewHandler(services)
	kc := &cons.KafkaConsumer{}
	if err := kc.ConsumeKafka(); err != nil {
		log.Fatalf("Failed to start consumer: %s", err)
	}

	defer kc.Stop()

	srv := new(messageservice.Server)
	go func() {
		if err := handlers.CreateTableSQL(); err != nil {
			logger.Error().Err(err).Msg("")
			logger.Fatal().Msg("There was an error with creating a table in database")
		}
	}()
	go func() {
		if err := srv.RunServer("8000", handlers.InitRoutes()); err != nil {
			errorLog.Fatal(err)
		}
	}()
	logger.Info().Msg("Server is running")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	go func() {
		<-quit
		kc.Stop()
		os.Exit(0)
	}()
	logger.Info().Msg("Server is shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("")
		logger.Fatal().Msg("There was an error while shutting down the server")
	}
	if err := dbpool.Close(); err != nil {
		logger.Error().Err(err).Msg("")
		logger.Fatal().Msg("There was an error while closing db connection")
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
