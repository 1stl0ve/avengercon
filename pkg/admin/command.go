package admin

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/1stl0ve/avengercon/pkg/api"
	"github.com/1stl0ve/avengercon/pkg/module/kill"
)

// the command type is used to define the 'built in' commands (i.e. exit, use,
// back, etc) that are used for controlling the admin client itself, not for
// tasking an Agent.
type command string

const (
	exitCommand   = "exit"
	agentsCommand = "agents"
	useCommand    = "use"
	backCommand   = "back"
)

var builtinCommands []command = []command{
	helpCommand,
	exitCommand,
	agentsCommand,
	useCommand,
	backCommand,
}

// checks if a user has input a builtin command.
func (cmd command) isBuiltIn() bool {
	for _, known := range builtinCommands {
		if cmd == known {
			return true
		}
	}
	return false
}

func (adm *Admin) handleInput(args []string) error {
	if len(args) == 0 {
		return errors.New("no command")
	}

	// the 'command' is going to be the first token
	cmd := command(args[0])

	// if the command is builtin, handle it appropriately
	if cmd.isBuiltIn() {
		return adm.handleBuiltIn(cmd, args[1:])
	}

	// if its not builtin, assume that its a module command
	if adm.currentAgent != nil {
		return adm.handleModule(args)
	}
	return nil
}

func (adm *Admin) handleBuiltIn(cmd command, args []string) error {
	switch cmd {
	case helpCommand:
		adm.help()
	case exitCommand:
		os.Exit(0) // does not return
	case useCommand:
		adm.use(args)
	case agentsCommand:
		return adm.listAgents()
	case backCommand:
		adm.back()
	default:
		return errors.New("invalid command")
	}
	return nil
}

func (adm *Admin) handleModule(args []string) error {
	// load the appropriate module
	mod := adm.modules[args[0]]
	if mod == nil {
		return errors.New("invalid module")
	}

	// create a task for the selected module
	task, err := mod.CreateTask(args)
	if err != nil {
		return err
	}
	task.AgentID = adm.currentAgent.GetAgentID()

	// send the task to the agent and wait for the response
	resp, err := adm.TaskAgent(context.Background(), task)
	if err != nil {
		return err
	}

	if resp.GetStatus() == api.Status_ERROR {
		fmt.Println(string(resp.GetData()))
	}

	// do something based on the response
	if err := mod.Do(resp); err != nil {
		return err
	}

	// the 'kill' module is unique; if its called, back out of the agent's
	// context and update the list of registered agents.
	if mod.Name() == kill.ModuleName {
		adm.back()
		return adm.updateAgents()
	}
	return nil
}

func (adm *Admin) back() {
	adm.currentAgent = nil
}
