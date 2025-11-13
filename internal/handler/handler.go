package handler

import (
	"review-rotator/internal/service"

	"github.com/gin-gonic/gin"
)



type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRouts() *gin.Engine {
	router := gin.New()
	
	return router
}