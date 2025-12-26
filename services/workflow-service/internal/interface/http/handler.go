package http

import (
	"net/http"
	"workflow-service/internal/application/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	createUsecase  *usecase.CreateWorkflowUsecase
	approveUsecase *usecase.ApproveWorkflowUsecase
	getUsecase     *usecase.GetWorkflowUsecase
}

func NewHandler(create *usecase.CreateWorkflowUsecase, approve *usecase.ApproveWorkflowUsecase, getUsecase *usecase.GetWorkflowUsecase) *Handler {
	return &Handler{createUsecase: create, approveUsecase: approve, getUsecase: getUsecase}
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

func (h *Handler) GetWorkflow(c *gin.Context) {
	id := c.Param("id")
	// cache
	workflow, err := h.getUsecase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "workflow not found"})
		return
	}
	c.JSON(http.StatusOK, workflow)
}
