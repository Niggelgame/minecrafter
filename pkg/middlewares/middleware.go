package middlewares

import (
	"github.com/niggelgame/minecrafter/pkg/commands"
)

type CommandPermissions struct {
	RequiredPermissions []int64
	// UserIDs
	AlwaysAllowedUsers        map[string]bool
}

func NewCommandPermissions(required []int64, alwaysAllowed []string) *CommandPermissions {
	allowedMap := make(map[string]bool)
	for _, u := range alwaysAllowed {
		allowedMap[u] = true
	}

	return &CommandPermissions{
		RequiredPermissions: required,
		AlwaysAllowedUsers:  allowedMap,
	}
}

func (c *CommandPermissions) CheckPermissions(ctx *commands.Ctx) bool {
	member := ctx.Message.Member

	_, exists := c.AlwaysAllowedUsers[member.User.ID]
	if exists {
		return true
	}

	return false
}
