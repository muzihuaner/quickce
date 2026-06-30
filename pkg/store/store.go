package store

import (
	"sync"
	"time"

	"quickce/pkg/probe"
)

type Agent struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	IP        string    `json:"ip"`
	LastSeen  time.Time `json:"last_seen"`
	CreatedAt time.Time `json:"created_at"`
}

type Task struct {
	ID        string         `json:"id"`
	AgentID   string         `json:"agent_id"`
	Req       ProbeRequest   `json:"request"`
	Result    *probe.Result  `json:"result,omitempty"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	DoneAt    *time.Time     `json:"done_at,omitempty"`
}

type ProbeRequest struct {
	Target  string         `json:"target"`
	Type    probe.ProbeType `json:"type"`
	Port    int            `json:"port"`
	Timeout int            `json:"timeout"`
}

type Store struct {
	mu      sync.RWMutex
	agents  map[string]*Agent
	tasks   map[string]*Task
}

func New() *Store {
	return &Store{
		agents: make(map[string]*Agent),
		tasks:  make(map[string]*Task),
	}
}

func (s *Store) AddAgent(a *Agent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.agents[a.ID] = a
}

func (s *Store) GetAgent(id string) *Agent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.agents[id]
}

func (s *Store) ListAgents() []*Agent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*Agent, 0, len(s.agents))
	for _, a := range s.agents {
		list = append(list, a)
	}
	return list
}

func (s *Store) AddTask(t *Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[t.ID] = t
}

func (s *Store) GetTask(id string) *Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tasks[id]
}

func (s *Store) PollTask(agentID string) *Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range s.tasks {
		if t.AgentID == agentID && t.Status == "pending" {
			t.Status = "running"
			return t
		}
	}
	return nil
}

func (s *Store) CompleteTask(id string, r *probe.Result) *Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := s.tasks[id]
	if t == nil {
		return nil
	}
	now := time.Now()
	t.Status = "done"
	t.Result = r
	t.DoneAt = &now
	return t
}

func (s *Store) ListTasks() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		list = append(list, t)
	}
	return list
}
