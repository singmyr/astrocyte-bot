package slack

import (
	"io"
	"strings"
)

// RequestData contains the information regarding the command called.
type RequestData struct {
	Token       string
	TeamID      string
	TeamDomain  string
	ChannelID   string
	ChannelName string
	UserID      string
	UserName    string
	Command     string
	Text        string
	ResponseURL string
	TriggerID   string
}

// DataFromBytes creates a slice of bytes into RequestData.
func DataFromBytes(bytes []byte) (*RequestData, error) {
	// Split it up into a map.
	data := map[string]string{}
	keyValues := strings.Split(string(bytes), "&")
	for _, value := range keyValues {
		d := strings.Split(value, "=")
		data[d[0]] = d[1]
	}
	return &RequestData{
		Token:       data["token"],
		TeamID:      data["team_id"],
		TeamDomain:  data["team_domain"],
		ChannelID:   data["channel_id"],
		ChannelName: data["channel_name"],
		UserID:      data["user_id"],
		UserName:    data["user_name"],
		Command:     cleanCommand(data["command"]),
		Text:        data["text"],
		ResponseURL: data["response_url"],
		TriggerID:   data["trigger_id"],
	}, nil
}

// cleanCommand trims off "%2F" (/) from the beginning of the string.
func cleanCommand(s string) string {
	return s[3:]
}

// Command defines a command.
type Command struct {
	Command string
	Handler func(w io.Writer, d *RequestData)
}

var commands = []*Command{}

// Handle takes the data from an incoming requests and executes the correct command.
func Handle(w io.Writer, d *RequestData) bool {
	// Check which command signature matches the incoming data.

	handled := false
	for _, command := range commands {
		if command.Command == d.Command {
			handled = true
			command.Handler(w, d)
			break
		}
	}

	if !handled {
		return false
	}

	return true
}

// RegisterCommand registers a command to be regognized.
func RegisterCommand(command *Command) error {
	commands = append(commands, command)
	return nil
}
