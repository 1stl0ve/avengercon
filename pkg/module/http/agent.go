package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/1stl0ve/avengercon/pkg/api"
	"github.com/1stl0ve/avengercon/pkg/module"
)

// Agent is a module.Agent implementation that allows a user of the admin tool
// to start, stop, or get the status of an HTTP server running on the remote
// system.
type Agent struct {
	Server  *http.Server
	message string
}

// Do executes the appropriate HTTP module operation.
func (m *Agent) Do(task *api.Task) error {
	var c config
	if err := module.DecodeConfig(task.Data, &c); err != nil {
		return err
	}

	switch c.Operation {
	case start:
		m.start(c)
	case stop:
		if err := m.stop(); err != nil {
			return err
		}
	case status:
		m.message = m.status()
	}
	return nil
}

// CreateResponse creates a task response with the appropriate status and server
// message.
func (m *Agent) CreateResponse(id string) (*api.Response, error) {
	resp := &api.Response{
		AgentID: id,
		Status:  api.Status_OK,
		Data:    []byte(m.message),
	}
	return resp, nil
}

// starts a new HTTP server on the remote system.
func (m *Agent) start(c config) {
	if m.Server != nil {
		m.message = "http server already running"
		return
	}

	m.Server = &http.Server{
		Addr: c.Addr,
	}

	m.message = "server started"
	go func() {
		log.Printf("starting server on %s\n", m.Server.Addr)
		if err := m.Server.ListenAndServe(); err != nil {
			m.message = fmt.Sprintf("error: %s\n", err.Error())
			log.Println(err)
			m.Server = nil
		}
	}()
	// give the HTTP server time to start
	time.Sleep(3 * time.Second)
}

// shutdown the HTTP server if it is already running.
func (m *Agent) stop() error {
	if m.Server == nil {
		m.message = "no http server running"
		return nil
	}

	if err := m.Server.Shutdown(context.Background()); err != nil {
		return err
	}
	m.Server = nil
	m.message = "server stopped"
	return nil
}

// returns a status message depending on if the HTTP server is running or not.
func (m *Agent) status() string {
	if m.Server == nil {
		return "no server running"
	}
	return "server already running"
}
