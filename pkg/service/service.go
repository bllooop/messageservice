package service

import (
	"messageservice"
	"messageservice/pkg/repository"
)

type Messages interface {
	Create(message messageservice.MessageItem) (int, error)
	CreateTable() error
}
type Service struct {
	Messages
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		Messages: NewMessageService(repository),
	}
}
