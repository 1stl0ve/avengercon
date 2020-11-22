package http

import (
	"errors"
	"net/http"
)

// ModuleName is the universal name for the HTTP server module.
const ModuleName = "http"

// the valid HTTP server operations
type operation int

const (
	start  operation = iota // start a new HTTP server
	stop                    // shutdown an existing HTTP server
	status                  // check if a HTTP server is running
)

var operations map[string]operation = map[string]operation{
	"start":  start,
	"stop":   stop,
	"status": status,
}

// struct members must be exported for serialization to work properly
type config struct {
	Operation operation // the operation the agent is supposed to execute
	Name      string
	Addr      string
	http.Handler
}

func startConfig(args []string) (config, error) {
	config := config{
		Name:      ModuleName,
		Operation: start,
	}

	if len(args) < 3 {
		return config, errors.New("not enough args for http start")
	}
	config.Addr = args[2]
	return config, nil
}

func newConfig(op operation, args []string) (config, error) {
	config := config{
		Name:      ModuleName,
		Operation: op,
	}
	switch op {
	case start:
		return startConfig(args)
	case stop, status:
		return config, nil
	default:
		return config, errors.New("invalid operation")
	}
}
