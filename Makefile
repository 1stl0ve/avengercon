default: api

api: pkg/api/agent.proto
	go generate ./...