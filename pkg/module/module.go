package module

import (
	"bytes"
	"encoding/gob"

	"github.com/1stl0ve/avengercon/pkg/api"
)

// Admin is an interface for creating and sending a task to an agent.
type Admin interface {
	Name() string
	CreateTask([]string) (*api.Task, error)
	Do(*api.Response) error
}

// Agent is an interface for handling tasks sent by clients.
type Agent interface {
	//api.AgentClient
	CreateResponse(string) (*api.Response, error)
	Do(*api.Task) error
	//SetAgent(api.AgentClient)
}

// EncodeConfig serializes a module configuration using the gob package.
func EncodeConfig(c interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(c); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// DecodeConfig unmarshalls an array of bytes into a module configuration.
func DecodeConfig(data []byte, c interface{}) error {
	var b bytes.Buffer
	_, err := b.Write(data)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(&b)
	if err := dec.Decode(c); err != nil {
		return err
	}
	return nil
}
