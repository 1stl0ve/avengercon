package kill

import (
	"fmt"

	"github.com/1stl0ve/avengercon/pkg/api"
)

// Admin tasks an Agent to cleanup and terminate.
type Admin struct{}

// Name returns the name of the module.
func (m *Admin) Name() string {
	return ModuleName
}

// CreateTask returns a new Kill module task. A kill command does not require
// any arguments, so any values passed in 'args' will be ignored.
func (m *Admin) CreateTask(args []string) (*api.Task, error) {
	return &api.Task{Name: "kill"}, nil
}

// Do prints out the results of the command that would've been written to stdout
// or stderr on the remote system.
func (m *Admin) Do(resp *api.Response) error {
	fmt.Printf("\n[%s:%s:%v]\n",
		resp.GetAgentID(),
		m.Name(),
		resp.GetStatus())
	return nil
}
