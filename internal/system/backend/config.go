package backend

import "emperror.dev/errors"

const (
	hostIsRequired     = errors.Sentinel("host is required")
	portIsRequired     = errors.Sentinel("port is required")
	usernameIsRequired = errors.Sentinel("username is required")
	passwordIsRequired = errors.Sentinel("password is required")
)

type Config struct {
	Host string
	Port int
}

type AuthenticatedConfig struct {
	*Config
	Username string
	Password string
}

func (c *Config) Validate() error {
	if c.Host == "" {
		return hostIsRequired
	}

	if c.Port == 0 {
		return portIsRequired
	}

	return nil
}

func (c *AuthenticatedConfig) Validate() error {
	if err := c.Config.Validate(); err != nil {
		return err
	}

	if c.Username == "" {
		return usernameIsRequired
	}

	if c.Password == "" {
		return passwordIsRequired
	}

	return nil
}
