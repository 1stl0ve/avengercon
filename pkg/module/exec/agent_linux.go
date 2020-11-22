package exec

import (
	"strings"
)

func getCommandTokens(cmd string) []string {
	return strings.Split(cmd, " ")
}
