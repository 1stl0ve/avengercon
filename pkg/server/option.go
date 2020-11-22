package server

// Option is a method for using functional options for setting up a Server.
type Option func(*Server)

// AdminAddr sets the ip address and port for the admin component of the server.
func AdminAddr(addr string) Option {
	return func(s *Server) {
		s.adminAddr = addr
	}
}

// AgentAddr sets the ip address and port for the agent component of the server.
func AgentAddr(addr string) Option {
	return func(s *Server) {
		s.agentAddr = addr
	}
}
