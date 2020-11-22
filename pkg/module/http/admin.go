package http

import (
	"errors"
	"fmt"

	"github.com/1stl0ve/avengercon/pkg/api"
	"github.com/1stl0ve/avengercon/pkg/module"
)

// Admin is a module.Client implementation that allows a user of the
// admin tool to start and stop an HTTP server on the remote system
type Admin struct{}

// Name returns the name of the module.
func (m *Admin) Name() string {
	return ModuleName
}

// CreateTask returns a api.Task_EXEC task populated with the commands provided
// by the user of the admin tool.
func (m *Admin) CreateTask(args []string) (*api.Task, error) {
	task := &api.Task{Name: ModuleName}

	if len(args) < 2 {
		return task, errors.New("http requires an operation")
	}

	c, err := newConfig(operations[args[1]], args)
	if err != nil {
		return task, err
	}

	if task.Data, err = module.EncodeConfig(c); err != nil {
		return task, err
	}
	return task, nil
}

// Do prints out the results of the command that would've been written to stdout
// or stderr on the remote system.
func (m *Admin) Do(resp *api.Response) error {
	fmt.Printf("\n[%s:http]\n%s\n",
		resp.GetAgentID(),
		string(resp.Data))
	return nil
}
