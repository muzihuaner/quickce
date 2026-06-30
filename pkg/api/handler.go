package api

import (
	"context"
	"net/http"
	"time"

	"quickce/pkg/probe"
	"quickce/pkg/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/sync/semaphore"
)

type Handler struct {
	sem   *semaphore.Weighted
	store *store.Store
}

type ProbeRequest struct {
	Target  string         `json:"target" binding:"required"`
	Type    probe.ProbeType `json:"type" binding:"required"`
	Port    int            `json:"port"`
	Timeout int            `json:"timeout"`
	Nodes   []string       `json:"nodes,omitempty"`
}

func NewHandler(maxConcurr int64, s *store.Store) *Handler {
	if s == nil {
		s = store.New()
	}
	return &Handler{
		sem:   semaphore.NewWeighted(maxConcurr),
		store: s,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/api/probe", h.HandleProbe)
}

func (h *Handler) HandleProbe(c *gin.Context) {
	var req ProbeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Nodes) > 0 {
		h.dispatchToNodes(c, &req)
		return
	}

	h.executeLocal(c, &req)
}

func (h *Handler) executeLocal(c *gin.Context, req *ProbeRequest) {
	if err := h.sem.Acquire(context.Background(), 1); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "too many requests"})
		return
	}
	defer h.sem.Release(1)

	timeout := req.Timeout
	if timeout <= 0 {
		timeout = 5
	}

	ctx, cancel := context.WithTimeout(context.Background(), probeTimeout(timeout))
	defer cancel()

	var result *probe.Result
	switch req.Type {
	case probe.HTTP:
		result = probe.ProbeHTTP(ctx, req.Target, timeout)
	case probe.Ping:
		result = probe.ProbePing(ctx, req.Target, timeout)
	case probe.TCP:
		if req.Port == 0 {
			req.Port = 80
		}
		result = probe.ProbeTCP(ctx, req.Target, req.Port, timeout)
	case probe.DNS:
		result = probe.ProbeDNS(ctx, req.Target, timeout)
	case probe.Traceroute:
		result = probe.ProbeTraceroute(ctx, req.Target, timeout)
	case probe.IP:
		result = probe.ProbeIP(ctx, req.Target, timeout)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported probe type: " + string(req.Type)})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) dispatchToNodes(c *gin.Context, req *ProbeRequest) {
	agents := h.store.ListAgents()
	if len(agents) == 0 {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "no agents available"})
		return
	}

	selected := filterAgents(agents, req.Nodes)
	if len(selected) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no matching agents found"})
		return
	}

	tasks := make([]gin.H, 0, len(selected))
	for _, a := range selected {
		id := uuid.New().String()
		task := &store.Task{
			ID:      id,
			AgentID: a.ID,
			Req: store.ProbeRequest{
				Target:  req.Target,
				Type:    req.Type,
				Port:    req.Port,
				Timeout: req.Timeout,
			},
			Status:    "pending",
			CreatedAt: time.Now(),
		}
		h.store.AddTask(task)
		tasks = append(tasks, gin.H{
			"task_id":  id,
			"agent_id": a.ID,
			"agent":    a.Name,
			"location": a.Location,
		})
	}

	c.JSON(http.StatusAccepted, gin.H{
		"dispatched": true,
		"tasks":      tasks,
		"message":    "tasks dispatched to agents, results will be reported asynchronously",
	})
}

func filterAgents(agents []*store.Agent, names []string) []*store.Agent {
	if len(names) == 0 {
		return agents
	}
	set := make(map[string]bool, len(names))
	for _, n := range names {
		set[n] = true
	}
	var result []*store.Agent
	for _, a := range agents {
		if set[a.Name] || set[a.ID] {
			result = append(result, a)
		}
	}
	return result
}

func probeTimeout(sec int) time.Duration {
	if sec <= 0 {
		return 5 * time.Second
	}
	return time.Duration(sec) * time.Second
}
