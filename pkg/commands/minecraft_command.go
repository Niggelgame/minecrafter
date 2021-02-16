package commands

import (
	"github.com/niggelgame/minecrafter/pkg/minecraft"
	"strings"
)

func NewMinecraftCommand(connection *minecraft.Connection) *Command {
	return &Command{
		Name:        "CommandExecution",
		Description: "Execute a minecraft command",
		Usage:       "",
		Invocations: []string{"e", "execute", "run", "eval"},
		SubCommands: []*Command{},
		Process: func(ctx *Ctx) error {
			return minecraftCommandHandler(ctx, connection)
		},
	}
}

func minecraftCommandHandler(ctx *Ctx, connection *minecraft.Connection) (e error) {
	if len(ctx.Args) == 0 {
		ctx.SendHelp()
		return
	}

	cmd := strings.Join(ctx.Args, " ")

	if connection == nil {
		e = sendError(ctx, "Connection to MC not existent")
		return &CommandError{Reason: "Connection to MC not existent"}
	}

	res, e := connection.ExecuteCommand(cmd)
	if e != nil {
		e = sendError(ctx, e.Error())
		return
	}

	_, e = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, res)
	if e != nil {
		return
	}

	return
}
