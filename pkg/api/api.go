//go:generate protoc -I . agent.proto --go-grpc_out=. --go_out=.

// Package api provides the communications protocol for the C2 framework.
// The protocol is implemented using gRPC.
package api

// EmptyMessage is a pointer to an Empty message. It can be returned whenever
// a "NULL" response is expected since NULL does not exist in gRPC.
var EmptyMessage *Empty = &Empty{}

// SetTaskStatus updates the status field of a given Task.
func SetTaskStatus(task *Task, s Status) {
	task.Status = s
}

// SetResponseStatus updates the status field of a given Response.
func SetResponseStatus(resp *Response, s Status) {
	resp.Status = s
}

// DefaultTask returns a Task with an OK status.
func DefaultTask() *Task {
	return &Task{
		Status: Status_OK,
	}
}
