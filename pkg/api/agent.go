package api

import (
	"net/http"
	"time"

	"quickce/pkg/probe"
	"quickce/pkg/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AgentHandler struct {
	store *store.Store
}

func NewAgentHandler(s *store.Store) *AgentHandler {
	return &AgentHandler{store: s}
}

func (h *AgentHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/api/agents/register", h.HandleRegister)
	r.GET("/api/agents", h.HandleListAgents)
	r.GET("/api/agents/:id/tasks", h.HandlePollTask)
	r.POST("/api/agents/:id/results", h.HandleReportResult)
	r.GET("/api/tasks/:id", h.HandleGetTask)
}

type registerReq struct {
	Name     string `json:"name" binding:"required"`
	Location string `json:"location" binding:"required"`
	ID       string `json:"id,omitempty"`
}

func (h *AgentHandler) HandleRegister(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := req.ID
	if id == "" || h.store.GetAgent(id) == nil {
		id = uuid.New().String()
	}

	agent := &store.Agent{
		ID:        id,
		Name:      req.Name,
		Location:  req.Location,
		IP:        c.ClientIP(),
		LastSeen:  time.Now(),
		CreatedAt: time.Now(),
	}
	h.store.AddAgent(agent)

	c.JSON(http.StatusOK, gin.H{
		"id":    id,
		"token": id,
	})
}

func (h *AgentHandler) HandleListAgents(c *gin.Context) {
	agents := h.store.ListAgents()
	c.JSON(http.StatusOK, agents)
}

func (h *AgentHandler) HandlePollTask(c *gin.Context) {
	agentID := c.Param("id")
	if h.store.GetAgent(agentID) == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	task := h.store.PollTask(agentID)
	if task == nil {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      task.ID,
		"target":  task.Req.Target,
		"type":    task.Req.Type,
		"port":    task.Req.Port,
		"timeout": task.Req.Timeout,
	})
}

type resultReq struct {
	TaskID  string      `json:"task_id" binding:"required"`
	Success bool        `json:"success"`
	Latency int64       `json:"latency_ms"`
	Error   string      `json:"error,omitempty"`
	Detail  interface{} `json:"detail,omitempty"`
}

func (h *AgentHandler) HandleReportResult(c *gin.Context) {
	agentID := c.Param("id")
	if h.store.GetAgent(agentID) == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	var req resultReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	r := &probe.Result{
		Success:   req.Success,
		LatencyMs: req.Latency,
		Error:     req.Error,
		Detail:    req.Detail,
	}
	task := h.store.CompleteTask(req.TaskID, r)
	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *AgentHandler) HandleGetTask(c *gin.Context) {
	taskID := c.Param("id")
	task := h.store.GetTask(taskID)
	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         task.ID,
		"agent_id":   task.AgentID,
		"status":     task.Status,
		"target":     task.Req.Target,
		"type":       task.Req.Type,
		"port":       task.Req.Port,
		"timeout":    task.Req.Timeout,
		"created_at": task.CreatedAt,
		"done_at":    task.DoneAt,
		"result":     task.Result,
	})
}
