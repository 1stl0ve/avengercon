package kill

import (
	"context"
	"fmt"
	"os"

	"github.com/1stl0ve/avengercon/pkg/api"
)

// Agent is a module.Agent implementation that ends the Agent process.
type Agent struct {
	api.AgentClient
}

// SetAgentClient sets the AgentClient
func (m *Agent) SetAgentClient(ac api.AgentClient) {
	m.AgentClient = ac
}

// Do unregisters the Agent and calls exit().
func (m *Agent) Do(task *api.Task) error {
	ctx := context.Background()
	resp := &api.Response{
		AgentID: task.GetAgentID(),
		Status:  api.Status_OK,
	}
	if _, err := m.SendResponse(ctx, resp); err != nil {
		return err
	}

	// let the AgentServer know that this agent is done
	id := &api.AgentID{
		AgentID: task.GetAgentID(),
	}
	if _, err := m.UnregisterAgent(ctx, id); err != nil {
		return err
	}
	fmt.Println("killed")
	os.Exit(0)
	return nil // unreachable
}

// CreateResponse returns a response with an OK status.
func (m *Agent) CreateResponse(id string) (*api.Response, error) {
	return &api.Response{
		AgentID: id,
		Status:  api.Status_OK,
	}, nil
}
