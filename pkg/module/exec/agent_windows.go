package exec

import (
	"strings"
)

func getCommandTokens(cmd string) []string {
	// for Windows, pass commands to powershell
	tokens := []string{"powershell", "-c"}
	return append(tokens, strings.Split(cmd, " ")...)
}
