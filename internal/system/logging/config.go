package logging

import "emperror.dev/errors"

const (
	encodingInvalid = errors.Sentinel("invalid encoding, must be one of [json, console]")
)

type Config struct {
	Level    string
	Dev      bool
	Encoding string
	Output   []string
}

func (c *Config) Validate() error {

	if c.Encoding != "" && (c.Encoding != "json" && c.Encoding != "console") {
		return encodingInvalid
	}

	return nil
}
