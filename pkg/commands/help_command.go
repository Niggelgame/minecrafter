package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

func NewHelpCommand() *Command {
	cmds := []string{"add", "help", "list", "remove", "execute"}

	return &Command{
		Name:        "Help",
		Description: "List all available commands",
		Usage:       "help",
		Invocations: []string{"help", "h"},
		SubCommands: []*Command{},
		Process: func(ctx *Ctx) error {
			e := &discordgo.MessageEmbed{
				Title:       fmt.Sprintf("%s - Usage of this Bot", ctx.Command.Name),
				Description: fmt.Sprintf("`%s`", strings.Join(cmds, "`, `")),
				Color:       0x7749a0,
				Timestamp:   time.Now().Format(time.RFC3339),
			}

			_, err := ctx.Session.ChannelMessageSendEmbed(ctx.ChannelID, e)

			return err
		},
	}
}
