package exec

import (
	"errors"
	"fmt"
	"strings"

	"github.com/1stl0ve/avengercon/pkg/api"
	"github.com/1stl0ve/avengercon/pkg/module"
)

// Admin is a module.Admin implementation that allows a user of the admin tool
// to send a command for the Agent to execute on the remote system.
type Admin struct {
	cmd string
}

// Name returns the name of the module.
func (m *Admin) Name() string {
	return ModuleName
}

// CreateTask sets up a command execution configuration and creates a task to
// send to a remote Agent.
func (m *Admin) CreateTask(args []string) (*api.Task, error) {
	var err error
	task := &api.Task{
		Name:   m.Name(),
		Status: api.Status_ERROR,
	}

	if len(args) <= 1 {
		return task, errors.New("command required")
	}

	config := config{
		Name: m.Name(),
		Cmd:  strings.TrimSuffix(strings.Join(args[1:], " "), "\n"),
	}
	m.cmd = config.Cmd

	task.Data, err = module.EncodeConfig(config)
	if err != nil {
		return nil, err
	}

	task.Status = api.Status_OK
	return task, nil
}

// Do prints out the results of the command that would've been written to stdout
// or stderr on the remote system.
func (m *Admin) Do(resp *api.Response) error {
	fmt.Printf("\n[%s:%s:%s:%v]\n\n%s\n",
		resp.GetAgentID(),
		m.Name(),
		m.cmd,
		resp.GetStatus(),
		string(resp.Data))
	return nil
}
