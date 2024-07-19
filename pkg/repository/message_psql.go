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

	createMessageQuery := fmt.Sprintf("INSERT INTO %s (text, time) VALUES ($1,$2)", messageTable)
	_, err = tr.Exec(createMessageQuery, message.Text, message.Time)
	if err != nil {
		tr.Rollback()
		return 0, err
	}
	return itemId, tr.Commit()

}

func (r *MessagePostgres) CreateTable() error {
	// Wait for a short time to simulate background processing
	time.Sleep(2 * time.Second)

	// SQL statement to create a table
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS messageTable (
        id SERIAL PRIMARY KEY,
        text VARCHAR(250) NOT NULL,
        time VARCHAR(100) UNIQUE NOT NULL
    );`

	// Execute the SQL statement
	_, err := r.db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	return nil
}
