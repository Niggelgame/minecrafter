package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

type Command struct {
	Name        string
	Description string
	Usage       string
	Invocations []string
	SubCommands []*Command
	Process     func(ctx *Ctx) error
}

type Ctx struct {
	Invoke     string
	FullInvoke string
	Args       []string
	Message    *discordgo.Message
	Author     *discordgo.User
	ChannelID  string
	Session    *discordgo.Session
	Handler    *CommandHandler
	Command    *Command
}

func (c *Ctx) SendHelp() {
	e := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s - Hilfe", c.Command.Name),
		Description: c.Command.Description,
		Color: 0x7749a0,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	fields := make([]*discordgo.MessageEmbedField, 0)
	fields = append(fields, &discordgo.MessageEmbedField{
		Name: "Aliase",
		Value: strings.Join(c.Command.Invocations, ", "),
		Inline: true,
	})
	for _, subcommand := range c.Command.SubCommands {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name: fmt.Sprintf("%s - %s", subcommand.Name, subcommand.Description),
			Value: fmt.Sprintf("%s %s", c.FullInvoke, subcommand.Invocations[0]),
		})
	}

	e.Fields = fields

	_, err := c.Session.ChannelMessageSendEmbed(c.ChannelID, e)

	if err != nil {
		fmt.Println(err)
	}
}

type CommandError struct {
	Reason string
}

func (c *CommandError) Error() string {
	return c.Reason
}
