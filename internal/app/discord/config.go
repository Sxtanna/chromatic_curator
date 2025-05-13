package discord

import "emperror.dev/errors"

const (
	tokenRequired = errors.Sentinel("token is required")
)

type BotConfiguration struct {
	Token   string
	GuildID string // Optional: If provided, slash commands will be registered for this guild only
}

func (c *BotConfiguration) Validate() error {
	if c.Token == "" {
		return tokenRequired
	}

	return nil
}
