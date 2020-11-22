package main

import (
	"flag"
	"fmt"

	"github.com/1stl0ve/avengercon/pkg/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	adminPort := flag.Int("admin", 9090, "the port to listen for incoming clients")
	agentPort := flag.Int("agent", 4444, "the port to listen for incoming agents")
	serverAddr := flag.String("addr", "127.0.0.1", "the ip address to bind to")
	debug := flag.Bool("debug", false, "run the server in debug mode")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	// default admin listener is localhost:9090
	adminAddr := fmt.Sprintf("%s:%d", *serverAddr, *adminPort)

	// default agent listener is localhost:4444
	agentAddr := fmt.Sprintf("%s:%d", *serverAddr, *agentPort)

	// create a new instance of the server
	srv := server.New(
		server.AdminAddr(adminAddr), // sets the admin listener address
		server.AgentAddr(agentAddr)) // sets the agent listener address

	// start the server
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
