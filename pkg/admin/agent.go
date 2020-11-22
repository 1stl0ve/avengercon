package admin

import "fmt"

// sets an agent as the active agent based on a given agent id.
func (adm *Admin) use(args []string) {
	if len(args) < 1 {
		fmt.Println("not enough arguments for command: use")
		return
	}

	id := args[0]
	agent, ok := adm.agents[id]
	if !ok {
		fmt.Printf("%s is not a valid agent id\n", id)
		return
	}
	adm.currentAgent = agent

	fmt.Printf("\nUsing Agent: %s\n", agent.GetAgentID())
	fmt.Printf("OS: %s\n", agent.GetOS())
	fmt.Println("IP addresses:")
	for _, ip := range agent.GetIP() {
		fmt.Printf("\t%s\n", ip)
	}
	fmt.Println()
}

// lists all available agents
func (adm *Admin) listAgents() error {
	if err := adm.updateAgents(); err != nil {
		return err
	}

	if len(adm.agents) == 0 {
		fmt.Println("no agents available")
		return nil
	}

	fmt.Println("available agents:")
	fmt.Println()
	for _, agent := range adm.agents {
		fmt.Printf("%s - %s - %s\n", agent.AgentID, agent.OS, agent.IP)
	}
	fmt.Println()
	return nil
}
