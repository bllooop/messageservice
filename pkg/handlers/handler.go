package handlers

import (
	"messageservice/pkg/service"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
func (h *Handler) InitRoutes() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("/", h.Home)
	r.HandleFunc("/create", h.createMessage)
	r.Handle("/metrics", promhttp.Handler())
	return r
}
func (h *Handler) CreateTableSQL() error {
	return h.services.Messages.CreateTable()
}
