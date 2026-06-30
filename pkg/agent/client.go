package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const pollInterval = 2 * time.Second

type AgentClient struct {
	ServerURL string
	ID        string
	Name      string
	Location  string
	Client    *http.Client
}

type registerRequest struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	ID       string `json:"id,omitempty"`
}

type registerResponse struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type taskResponse struct {
	ID        string `json:"id"`
	Target    string `json:"target"`
	Type      string `json:"type"`
	Port      int    `json:"port"`
	Timeout   int    `json:"timeout"`
}

type resultRequest struct {
	TaskID  string      `json:"task_id"`
	Success bool        `json:"success"`
	Latency int64       `json:"latency_ms"`
	Error   string      `json:"error,omitempty"`
	Detail  interface{} `json:"detail,omitempty"`
}

func New(server, name, location, agentID string) *AgentClient {
	return &AgentClient{
		ServerURL: server,
		ID:        agentID,
		Name:      name,
		Location:  location,
		Client:    &http.Client{Timeout: 30 * time.Second},
	}
}

func (a *AgentClient) Register(ctx context.Context) error {
	body := registerRequest{Name: a.Name, Location: a.Location, ID: a.ID}
	data, _ := json.Marshal(body)
	resp, err := a.Client.Post(a.ServerURL+"/api/agents/register", "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("register failed: %w", err)
	}
	defer resp.Body.Close()

	var res registerResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("decode register response: %w", err)
	}
	a.ID = res.ID
	log.Printf("Registered as agent %s (ID: %s)", a.Name, a.ID)
	return nil
}

func (a *AgentClient) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			a.pollAndExecute(ctx)
			time.Sleep(pollInterval)
		}
	}
}

func (a *AgentClient) pollAndExecute(ctx context.Context) {
	if a.ID == "" {
		return
	}

	task, err := a.pollTask(ctx)
	if err != nil {
		log.Printf("Poll task error: %v", err)
		return
	}
	if task == nil {
		return
	}

	log.Printf("Executing task %s: %s %s", task.ID, task.Type, task.Target)
	result := a.executeTask(ctx, task)
	a.reportResult(ctx, task.ID, result)
}

func (a *AgentClient) pollTask(ctx context.Context) (*taskResponse, error) {
	url := fmt.Sprintf("%s/api/agents/%s/tasks", a.ServerURL, a.ID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	var task taskResponse
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("decode task: %w", err)
	}
	if task.ID == "" {
		return nil, nil
	}
	return &task, nil
}

func (a *AgentClient) executeTask(ctx context.Context, task *taskResponse) *probeResult {
	return executeProbe(ctx, task.Target, task.Type, task.Port, task.Timeout)
}

func (a *AgentClient) reportResult(ctx context.Context, taskID string, r *probeResult) {
	body := resultRequest{
		TaskID:  taskID,
		Success: r.Success,
		Latency: r.LatencyMs,
		Error:   r.Error,
		Detail:  r.Detail,
	}
	data, _ := json.Marshal(body)
	url := fmt.Sprintf("%s/api/agents/%s/results", a.ServerURL, a.ID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		log.Printf("Report result error: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.Client.Do(req)
	if err != nil {
		log.Printf("Report result error: %v", err)
		return
	}
	resp.Body.Close()
	log.Printf("Task %s completed, success=%v latency=%dms", taskID, r.Success, r.LatencyMs)
}
