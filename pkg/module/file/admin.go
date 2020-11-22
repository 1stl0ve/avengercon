package file

import (
	"errors"

	"github.com/1stl0ve/avengercon/pkg/api"
	"github.com/1stl0ve/avengercon/pkg/module"
)

// Admin is a module.Admin implementation that allows a user of the admin tool
// to send a task for the Agent to upload or download a file to/from a remote
// system.
type Admin struct {
	Config config
}

// Name returns the name of the module.
func (m *Admin) Name() string {
	switch m.Config.Operation {
	case upload:
		return UploadModuleName
	case download:
		return DownloadModuleName
	default:
		return "unimplemented"
	}
}

// CreateTask sets up a download or upload configuration and creates a task to
// send to a remote Agent.
func (m *Admin) CreateTask(args []string) (*api.Task, error) {
	var err error
	task := &api.Task{
		Name:   m.Name(),
		Status: api.Status_ERROR,
	}

	if len(args) != 3 {
		return task, errors.New("invalid number of arguments")
	}

	op := operations[m.Name()]
	switch op {
	case upload:
		if err := m.setUploadConfig(args); err != nil {
			return task, err
		}
	case download:
		m.setDownloadConfig(args)
	default:
		return task, errors.New("invalid operation")
	}

	task.Data, err = module.EncodeConfig(m.Config)
	if err != nil {
		return task, err
	}

	task.Status = api.Status_OK
	return task, nil
}

// Do writes the contents from the response to the local path.
func (m *Admin) Do(resp *api.Response) error {
	switch operations[m.Name()] {
	case upload:
		m.upload(resp.GetAgentID())
	case download:
		if err := m.download(resp.GetAgentID(), resp.Data); err != nil {
			return err
		}
	}
	return nil
}
