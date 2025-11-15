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
	
	team := router.Group("/team")
	{
		team.POST("/add", h.addTeam)
		team.GET("/get", h.getTeam)
	}

	user := router.Group("/users")
	{
		user.POST("/setIsActive", h.setUserActive)
	}
	pr := router.Group("/pullRequest")
	{
		pr.POST("/create", h.createPullRequest)
		pr.POST("/merge", h.mergePullRequest)
		pr.POST("/reassign", h.reassignReviewer)
	}

	return router
}