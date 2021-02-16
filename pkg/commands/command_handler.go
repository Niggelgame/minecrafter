package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

type CommandHandler struct {
	Registry *CommandRegistry
	Prefix   string
}

func NewHandler(registry *CommandRegistry, prefix string) *CommandHandler {
	return &CommandHandler{
		Registry: registry,
		Prefix:   prefix,
	}
}

func (h *CommandHandler) HandleMessageEvent(s *discordgo.Session, evt *discordgo.MessageCreate) {
	// Deny Bots
	if evt.Message.Author.Bot {
		return
	}

	raw := evt.Message.Content
	lower := strings.ToLower(raw)

	// Check for prefix
	if !strings.HasPrefix(lower, h.Prefix) {
		return
	}

	// [ , invoke, args]
	parts := strings.Split(lower, h.Prefix)
	if len(parts) < 2 {
		return
	}

	args := strings.Split(parts[1], " ")

	cmd := h.Registry.Invokes[args[0]]
	if cmd == nil {
		return
	}
	h.processCommand(s, evt, cmd, args)
}

func (h *CommandHandler) processCommand(s *discordgo.Session, evt *discordgo.MessageCreate, cmd *Command, args []string) {
	invoke := args[0]
	if len(args) > 1 {
		subInvoke := args[1]
		for _, sub := range cmd.SubCommands {
			for _, inv := range sub.Invocations {
				if inv == subInvoke {
					h.processCommand(s, evt, sub, args[1:])
					return
				}
			}
		}
	}
	ctx := &Ctx{
		Invoke:     invoke,
		FullInvoke: fmt.Sprintf("%s%s", strings.Split(strings.ToLower(evt.Message.Content), strings.Join(args, " "))[0], invoke),
		Args:       args[1:],
		Session:    s,
		Message:    evt.Message,
		ChannelID:  evt.Message.ChannelID,
		Author:     evt.Message.Author,
		Handler:    h,
		Command:    cmd,
	}
	err := cmd.Process(ctx)
	if err != nil {
		h.handleError(ctx, err)
	}
}

func (h *CommandHandler) handleError(ctx *Ctx, err error) {

}