package server

import (
	"context"

	"github.com/1stl0ve/avengercon/pkg/api"
)

// an implementation of an api.AdminServer
type admin struct {
	api.UnimplementedAdminServer
	work   map[string]chan *api.Task
	output map[string]chan *api.Response
}

func newAdmin(work map[string]chan *api.Task, output map[string]chan *api.Response) *admin {
	return &admin{
		work:   work,
		output: output,
	}
}

// TaskAgent sends a new task to a specified Agent
func (a *admin) TaskAgent(ctx context.Context, task *api.Task) (*api.Response, error) {
	id := task.GetAgentID()

	go func() {
		a.work[id] <- task
	}()

	resp := <-a.output[id]
	return resp, nil
}

// GetAgents returns a list of Agents registered with the server
func (a *admin) GetAgents(ctx context.Context, _ *api.Empty) (*api.RegisteredAgents, error) {
	var agents *api.RegisteredAgents = new(api.RegisteredAgents)
	for _, a := range registeredAgents {
		agents.Agents = append(agents.Agents, a)
	}
	return agents, nil
}
