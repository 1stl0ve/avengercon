// Package admin provides the functionality for a user to interact with a remote
// agent. The package utilizes implementations of the module.Agent interface
// to define what tasks an Admin client is capable of sending to an Agent.
package admin

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/1stl0ve/avengercon/pkg/api"
	"github.com/1stl0ve/avengercon/pkg/module"
	"github.com/1stl0ve/avengercon/pkg/module/kill"
	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Admin is a client for interacting with and tasking Agents.
//
// Admin is a wrapper around an api.AdminClient, which provides the methods for
// sending gRPC messages to the server.
type Admin struct {
	// an AdminClient is embedded within the Admin for gRPC functionality.
	api.AdminClient

	// agents is a map of the agents that are currently registered with the server.
	agents map[string]*api.Registration

	// currentAgent is the agent selected by the 'use' command.
	currentAgent *api.Registration

	// modules is a list of modules that this Admin was configured with.
	modules map[string]module.Admin
}

// New creates a new Admin and sets the default values
func New(addr string) (*Admin, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return &Admin{}, err
	}

	adm := &Admin{
		AdminClient: api.NewAdminClient(conn),
		agents:      make(map[string]*api.Registration),
		modules:     make(map[string]module.Admin),
	}

	// the 'kill' module is included by default on all Admin clients.
	adm.AddModule("kill", &kill.Admin{})
	return adm, nil
}

// AddModule registers a module.Admin with the client.
func (adm *Admin) AddModule(name string, m module.Admin) error {
	adm.modules[name] = m
	return nil
}

func (adm *Admin) prompt() string {
	if adm.currentAgent == nil {
		return "> "
	}
	// the uuid is not the most useful indicator... but a prompt that begins
	// with '(UUID)' indicates that the client is in an agent's context.
	return fmt.Sprintf("(%s)> ", adm.currentAgent.GetAgentID())
}

// Run starts the admin client and begins a user input loop.
func (adm *Admin) Run() {
	// periodically update the Agents list
	go func() {
		if err := adm.updateAgents(); err != nil {
			fmt.Println(err)
		}
		time.Sleep(10 * time.Second)
	}()

	// enter command loop
	for {
		// get the user's input (augmented with suggestions)
		text := prompt.Input(adm.prompt(), adm.suggestions)

		// tokenize and handle the user's input
		tokens := strings.Split(text, " ")
		if err := adm.handleInput(tokens); err != nil {
			log.Error(err)
		}
	}
}

func (adm *Admin) updateAgents() error {
	// get an updated list of registered agents from the server
	agents, err := adm.GetAgents(context.Background(), api.EmptyMessage)
	if err != nil {
		return err
	}

	// update the admin's agent tracker
	adm.agents = make(map[string]*api.Registration)
	for _, a := range agents.Agents {
		adm.agents[a.GetAgentID()] = a
	}
	return nil
}
