package repository

import (
	"fmt"
	"messageservice"
	"time"

	"github.com/jmoiron/sqlx"
)

type MessagePostgres struct {
	db *sqlx.DB
}

func NewMessagePostgres(db *sqlx.DB) *MessagePostgres {
	return &MessagePostgres{db: db}
}

func (r *MessagePostgres) Create(message messageservice.MessageItem) (int, error) {
	tr, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var itemId int

	createMessageQuery := fmt.Sprintf("INSERT INTO %s (text, time) VALUES ($1,$2) RETURNING id", messageTable)
	row := tr.QueryRow(createMessageQuery, message.Text, message.Time)
	if err := row.Scan(&itemId); err != nil {
		tr.Rollback()
		return 0, err
	}
	return itemId, tr.Commit()

}

func (r *MessagePostgres) CreateTable() error {
	time.Sleep(2 * time.Second)

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS messageTable (
        id SERIAL PRIMARY KEY,
        text VARCHAR(250) NOT NULL,
        time TIMESTAMP NOT NULL
    );`

	_, err := r.db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	return nil
}
