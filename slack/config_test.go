package slack

import (
	"testing"

	"github.com/slack-go/slack"
)

func Test_Client(t *testing.T) {
	config := &Config{
		Token: "test token",
	}

	client, err := config.Client()

	if err != nil {
		t.Fatal(err)
	}

	if client.(*slack.Client) == nil {
		t.Fatalf("required non-nil client")
	}
}
