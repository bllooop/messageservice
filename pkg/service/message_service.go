package service

import (
	"messageservice"
	"messageservice/pkg/repository"
)

type MessageService struct {
	repository repository.Messages
}

func NewMessageService(repository repository.Messages) *MessageService {
	return &MessageService{repository: repository}
}

func (s *MessageService) Create(message messageservice.MessageItem) (int, error) {
	return s.repository.Create(message)
}

func (s *MessageService) CreateTable() error {
	return s.repository.CreateTable()
}
