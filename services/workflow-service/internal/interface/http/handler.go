package http

import (
	"net/http"
	"workflow-service/internal/application/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	createUsecase  *usecase.CreateWorkflowUsecase
	approveUsecase *usecase.ApproveWorkflowUsecase
}

func NewHandler(create *usecase.CreateWorkflowUsecase, approve *usecase.ApproveWorkflowUsecase) *Handler {
	return &Handler{createUsecase: create, approveUsecase: approve}
}

func (h *Handler) Create(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.createUsecase.Create(req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func (h *Handler) Approve(c *gin.Context) {
	id := c.Param("id")
	if err := h.approveUsecase.Approve(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "approved"})
}
