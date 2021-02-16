package commands

type CommandRegistry struct {
	Commands []*Command
	// For easier access to these commands and to prevent a for loop for each
	Invokes map[string]*Command
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		Commands: make([]*Command, 0),
		Invokes: make(map[string]*Command),
	}
}

func (h *CommandRegistry) RegisterCommand(cmd *Command) {
	h.Commands = append(h.Commands, cmd)
	for _, invoke := range cmd.Invocations {
		h.Invokes[invoke] = cmd
	}
}