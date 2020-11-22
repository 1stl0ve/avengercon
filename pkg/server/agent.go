package server

import (
	"context"
	"errors"

	"github.com/1stl0ve/avengercon/pkg/api"
	log "github.com/sirupsen/logrus"
)

var registeredAgents = map[string]*api.Registration{}

// an implementation of an api.AgentServer
type agent struct {
	api.UnimplementedAgentServer
	work   map[string]chan *api.Task
	output map[string]chan *api.Response
}

func newAgent(work map[string]chan *api.Task, output map[string]chan *api.Response) *agent {
	return &agent{
		work:   work,
		output: output,
	}
}

// GetTask handles requests from an Agent for new Tasks. The specific
// Agent is defined by an AgentInfo argument. The server will return a Task
// from that Agent's work channel. If there are not any Tasks in the channel,
// the server will respond with a Task that does not have the Agent member
// set. This will inform the Agent that there are no tasks to execute.
func (a *agent) GetTask(ctx context.Context, msg *api.AgentID) (*api.Task, error) {
	var task *api.Task = new(api.Task)
	select {
	case task, ok := <-a.work[msg.GetAgentID()]:
		if ok {
			return task, nil
		}
		return task, errors.New("channel closed")
	default:
		return task, nil
	}
}

// SendResponse retrieves the result of a Task from an Agent and adds it to
// that Agent's output channel. Upon success, the server returns an Empty message
// to the Agent.
func (a *agent) SendResponse(ctx context.Context, resp *api.Response) (*api.Empty, error) {
	a.output[resp.GetAgentID()] <- resp

	// if there was an error handling the command, the error message will be set
	// in the Command's Out field. The server should log any errors that the
	// agent encounters. This function will only return an error if there is an
	// issue handling/forwarding the response to a client.
	if resp.GetStatus() == api.Status_ERROR {
		log.Error(string(resp.GetData()))
	}
	return api.EmptyMessage, nil
}

// RegisterAgent adds a new Agent to the registeredAgents map. It creates both
// the work and output channels for the new Agent. On success, the server returns
// an Empty message to the Agent.
func (a *agent) RegisterAgent(ctx context.Context, reg *api.Registration) (*api.Empty, error) {
	id := reg.GetAgentID()

	a.work[id] = make(chan *api.Task)
	a.output[id] = make(chan *api.Response)
	registeredAgents[id] = reg
	log.WithFields(log.Fields{
		"agent":    reg.GetAgentID(),
		"ip":       reg.GetIP(),
		"hostname": reg.GetHostname(),
		"os":       reg.GetOS(),
	}).Info("registered new agent")
	return api.EmptyMessage, nil
}

// UnregisterAgent removes an Agent from the registeredAgents map. It deletes all
// of the data structures used for tracking the Agent and sends an Emtpy message
// back to the Agent.
func (a *agent) UnregisterAgent(ctx context.Context, msg *api.AgentID) (*api.Empty, error) {
	id := msg.GetAgentID()
	delete(a.work, id)
	delete(a.output, id)
	delete(registeredAgents, id)
	log.WithFields(log.Fields{
		"agent": id,
	}).Info("unregistered agent")
	return api.EmptyMessage, nil
}
