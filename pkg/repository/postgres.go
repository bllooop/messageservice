package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	messageTable = "messageTable"
	//host = "db"   //comment when starting on local without docker-compose
	//port = "5432" //comment when starting on local without docker-compose
//	host     = "localhost" //uncomment when starting on local without docker-compose
//	port     = "5433"      //uncomment when starting on local without docker-compose
//	user     = "postgres"
//dbname   = "postgres"
//sslmode  = "disable"
//password = "54321"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBname   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBname, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
