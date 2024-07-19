package handlers

import (
	"messageservice/pkg/service"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
func (h *Handler) InitRoutes() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("/create", h.createMessage)
	return r
}
func (h *Handler) CreateTableSQL() error {
	return h.services.Messages.CreateTable()
}
