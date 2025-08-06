// Filename: api/handlers.go
package handler

import (
	"github.com/mrjvadi/BackendPanelVpn/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}
