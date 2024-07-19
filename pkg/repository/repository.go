package repository

import (
	"messageservice"

	"github.com/jmoiron/sqlx"
)

type Messages interface {
	Create(message messageservice.MessageItem) (int, error)
	CreateTable() error
}

type Repository struct {
	Messages
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Messages: NewMessagePostgres(db),
	}
}
