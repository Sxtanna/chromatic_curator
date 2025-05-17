package discord

import "emperror.dev/errors"

const (
	tokenRequired = errors.Sentinel("token is required")
)

type BotConfiguration struct {
	Token  string
	Admins string
}

func (c *BotConfiguration) Validate() error {
	if c.Token == "" {
		return tokenRequired
	}

	return nil
}
