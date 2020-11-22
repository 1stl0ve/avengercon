package server

import (
	"net"

	"github.com/1stl0ve/avengercon/pkg/api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Server is a structure that manages transfering tasks between admin clients
// and agents.
type Server struct {
	adminAddr string
	*admin
	agentAddr string
	*agent
}

// New returns a new instance of Server and creates both the adminServer and
// the agentServer that run inside of it.
func New(options ...Option) *Server {
	work := make(map[string]chan *api.Task)
	output := make(map[string]chan *api.Response)

	// default server configuration values
	server := &Server{
		adminAddr: "127.0.0.1:9090",
		admin:     newAdmin(work, output),
		agentAddr: "127.0.0.1:4444",
		agent:     newAgent(work, output),
	}

	for _, option := range options {
		option(server)
	}
	return server
}

// Run creates a listeners for each of the gRPC servers and runs both of them.
func (s *Server) Run() error {
	go s.startAgent()
	return s.startAdmin()
}

func (s *Server) startAgent() error {
	log.WithFields(log.Fields{
		"addr": s.agentAddr,
	}).Info("starting agent listener")
	l, err := net.Listen("tcp", s.agentAddr)
	if err != nil {
		return err
	}
	gs := grpc.NewServer([]grpc.ServerOption{}...)
	api.RegisterAgentServer(gs, s.agent)
	return gs.Serve(l)
}

func (s *Server) startAdmin() error {
	log.WithFields(log.Fields{
		"addr": s.adminAddr,
	}).Info("starting admin listener")
	l, err := net.Listen("tcp", s.adminAddr)
	if err != nil {
		return err
	}
	gs := grpc.NewServer([]grpc.ServerOption{}...)
	api.RegisterAdminServer(gs, s.admin)
	return gs.Serve(l)
}
