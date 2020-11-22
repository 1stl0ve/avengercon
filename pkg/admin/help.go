package admin

import "fmt"

var helpCommand command = "help"

func (adm *Admin) help() {
	fmt.Println("available commands:")
	for _, cmd := range builtinCommands {
		fmt.Println(cmd)
	}

	if adm.currentAgent != nil {
		fmt.Printf("\nagent commands:\n")
		for _, cmd := range adm.modules {
			fmt.Println(cmd.Name())
		}
	}
}
