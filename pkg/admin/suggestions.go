package admin

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

var noSuggestions = []prompt.Suggest{}

// suggestions provides auto-complete suggestions
func (adm *Admin) suggestions(d prompt.Document) []prompt.Suggest {
	// populate command suggestions
	suggestions := []prompt.Suggest{}
	for _, cmd := range builtinCommands {
		suggestions = append(
			suggestions,
			prompt.Suggest{Text: string(cmd)})
	}
	if adm.currentAgent != nil {
		for _, cmd := range adm.modules {
			suggestions = append(
				suggestions,
				prompt.Suggest{Text: cmd.Name()})
		}
	}

	args := strings.Split(d.TextBeforeCursor(), " ")
	if len(args) == 1 {
		return prompt.FilterHasPrefix(suggestions, args[0], true)
	}

	cmd := command(args[0])
	switch cmd {
	case useCommand:
		// the use command should suggest any available agents as the next
		// argument

		// if there is already a uuid present in the arguments, do not suggest
		// anything else
		if len(args) > 2 {
			return noSuggestions
		}

		// otherwise, suggest any available agents that have a prefix that match
		// what the user is currently typing
		suggestions = []prompt.Suggest{}
		for id := range adm.agents {
			suggestions = append(suggestions, prompt.Suggest{Text: id})
		}
		return prompt.FilterHasPrefix(suggestions, args[1], true)

	case agentsCommand, helpCommand, exitCommand:
		// these commands do not have any additional arguments
		return noSuggestions
	default:
		// if the command isn't found in the main context, pass through and
		// check the agent context
	}

	// none of the agent commands require suggestions
	if _, ok := adm.modules[args[0]]; ok {
		return noSuggestions
	}

	// by default, recommend whatever commands are available
	return prompt.FilterHasPrefix(suggestions, args[0], true)
}
