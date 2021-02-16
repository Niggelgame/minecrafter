package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

func sendNoPermission(ctx *Ctx) {

}

func sendError(ctx *Ctx, exception string) error {
	e := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s - Error", ctx.Command.Name),
		Description: exception,
		Color:       0x7749a0,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	fields := make([]*discordgo.MessageEmbedField, 0)
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "What to do now?",
		Value:  "Do not hesitate to contact me on [Github](https://github.com/niggelgame)",
		Inline: true,
	})

	e.Fields = fields

	_, err := ctx.Session.ChannelMessageSendEmbed(ctx.ChannelID, e)

	return err
}
