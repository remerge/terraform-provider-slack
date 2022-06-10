package slack

import (
	"github.com/slack-go/slack"
)

const (
	ctxId = 1
)

type Config struct {
	Token string
}

func (c *Config) Client() (interface{}, error) {
	return slack.New(c.Token), nil
}
