package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/niggelgame/minecrafter/pkg/minecraft"
	"strings"
	"time"
)

func NewWhitelistCommand(connection *minecraft.Connection) *Command {
	return &Command{
		Name:        "Whitelist",
		Description: "Check or edit the whitelist",
		Usage:       "",
		Invocations: []string{"wl", "whitelist", "whl"},
		SubCommands: []*Command{addWhitelist(connection), removeWhitelist(connection), listWhitelist(connection)},
		Process: func(ctx *Ctx) error {
			ctx.SendHelp()
			return nil
		},
	}
}

func addWhitelist(connection *minecraft.Connection) *Command {
	return &Command{
		Name:        "Add whitelist",
		Description: "Add a user to the Whitelist",
		Usage:       "",
		Invocations: []string{"add"},
		Process: func(ctx *Ctx) error {
			return addUsersToWhitelist(ctx, connection)
		},
	}
}

func addUsersToWhitelist(ctx *Ctx, connection *minecraft.Connection) (e error) {
	if len(ctx.Args) == 0 {
		ctx.SendHelp()
		return
	}

	if connection == nil {
		e = sendError(ctx, "Connection to MC not existent")
		return &CommandError{Reason: "Connection to MC not existent"}
	}

	successUsers := make([]string, 0)
	alreadyUsers := make([]string, 0)
	failedUsers := make([]string, 0)

	for _, user := range ctx.Args {
		cmd := fmt.Sprintf("whitelist add %s", user)
		res, err := connection.ExecuteCommand(cmd)
		if err != nil {
			e = sendError(ctx, err.Error())
			return
		}
		if strings.HasPrefix(res, "Added") {
			successUsers = append(successUsers, user)
		} else if strings.HasSuffix(res, "already whitelisted") {
			alreadyUsers = append(alreadyUsers, user)
		} else {
			failedUsers = append(failedUsers, user)
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s", ctx.Command.Name),
		Description: "Successfully worked through users",
		Color:       0x7749a0,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	addedStr := fmt.Sprintf("Adding users `%s`", strings.Join(successUsers, "`, `"))
	alreadyStr := fmt.Sprintf("Already added: `%s`", strings.Join(alreadyUsers, "`, `"))
	failedStr := fmt.Sprintf("Not found: `%s`", strings.Join(failedUsers, "`, `"))

	fields := make([]*discordgo.MessageEmbedField, 0)

	if len(successUsers) > 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Success",
			Value:  addedStr,
			Inline: true,
		})
	}
	if len(alreadyUsers) > 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Already",
			Value:  alreadyStr,
			Inline: true,
		})
	}
	if len(failedUsers) > 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Failed",
			Value:  failedStr,
			Inline: true,
		})
	}

	if len(fields) < 1 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Action",
			Value:  "No action completed",
			Inline: true,
		})
	}

	embed.Fields = fields

	_, e = ctx.Session.ChannelMessageSendEmbed(ctx.ChannelID, embed)

	return
}

func removeWhitelist(connection *minecraft.Connection) *Command {
	return &Command{
		Name:        "Remove whitelist",
		Description: "Remove a user / users from the Whitelist",
		Usage:       "",
		Invocations: []string{"rm", "remove"},
		Process: func(ctx *Ctx) error {
			return removeUsersFromWhitelist(ctx, connection)
		},
	}
}

func removeUsersFromWhitelist(ctx *Ctx, connection *minecraft.Connection) (e error) {
	if len(ctx.Args) == 0 {
		ctx.SendHelp()
		return
	}

	if connection == nil {
		e = sendError(ctx, "Connection to MC not existent")
		return &CommandError{Reason: "Connection to MC not existent"}
	}

	successUsers := make([]string, 0)
	alreadyRemovedUsers := make([]string, 0)
	failedUsers := make([]string, 0)

	for _, user := range ctx.Args {
		cmd := fmt.Sprintf("whitelist remove %s", user)
		res, err := connection.ExecuteCommand(cmd)
		if err != nil {
			e = sendError(ctx, err.Error())
			return
		}

		if strings.HasPrefix(res, "Removed") {
			successUsers = append(successUsers, user)
		} else if strings.HasSuffix(res, "not whitelisted") {
			alreadyRemovedUsers = append(alreadyRemovedUsers, user)
		} else {
			failedUsers = append(failedUsers, user)
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s", ctx.Command.Name),
		Description: "Successfully worked through users",
		Color:       0x7749a0,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	removedStr := fmt.Sprintf("Removing users `%s`", strings.Join(successUsers, "`, `"))
	alreadyStr := fmt.Sprintf("Already removed: `%s`", strings.Join(alreadyRemovedUsers, "`, `"))
	failedStr := fmt.Sprintf("Not found: `%s`", strings.Join(failedUsers, "`, `"))

	fields := make([]*discordgo.MessageEmbedField, 0)

	if len(successUsers) > 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Success",
			Value:  removedStr,
			Inline: true,
		})
	}
	if len(alreadyRemovedUsers) > 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Already Removed / Not on whitelist",
			Value:  alreadyStr,
			Inline: true,
		})
	}
	if len(failedUsers) > 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Failed",
			Value:  failedStr,
			Inline: true,
		})
	}

	if len(fields) < 1 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Action",
			Value:  "No action completed",
			Inline: true,
		})
	}

	embed.Fields = fields

	_, e = ctx.Session.ChannelMessageSendEmbed(ctx.ChannelID, embed)

	return
}

func listWhitelist(connection *minecraft.Connection) *Command {
	return &Command{
		Name:        "List whitelist",
		Description: "List the Whitelist",
		Usage:       "",
		Invocations: []string{"list", "l"},
		Process: func(ctx *Ctx) error {
			return listUsersFromWhitelist(ctx, connection)
		},
	}
}

func listUsersFromWhitelist(ctx *Ctx, connection *minecraft.Connection) (e error) {
	if connection == nil {
		e = sendError(ctx, "Connection to MC not existent")
		return &CommandError{Reason: "Connection to MC not existent"}
	}

	uList := make([]string, 0)

	cmd := "whitelist list"
	res, err := connection.ExecuteCommand(cmd)
	if err != nil {
		e = sendError(ctx, err.Error())
		return
	}

	if !strings.HasSuffix(res, "no whitelisted players") {
		split := strings.Split(res, ": ")

		if len(split) > 1 {
			split = strings.Split(split[1], ", ")
			uList = split
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s", ctx.Command.Name),
		Description: fmt.Sprintf("Users: `%s`", strings.Join(uList, "`, `")),
		Color:       0x7749a0,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	_, e = ctx.Session.ChannelMessageSendEmbed(ctx.ChannelID, embed)

	return
}
