package messenger

import (
	"github.com/Maelkum/b7s/models/request"
	"github.com/Maelkum/b7s/models/response"
)

type Message struct {
	Name               string // name of the message
	Definition         any    // b7s type to be sent
	HasResponse        bool   // should we expect response for this message
	ResponseDefinition any    // b7s type we expect back
}

var messages = []Message{
	{
		Name:               "execute",
		Definition:         request.Execute{},
		HasResponse:        true,
		ResponseDefinition: response.Execute{},
	},
	{
		Name:               "work-order",
		Definition:         request.WorkOrder{},
		HasResponse:        true,
		ResponseDefinition: response.WorkOrder{},
	},
}

func SupportedMessages() []string {
	msgs := make([]string, len(messages))

	for i, msg := range messages {
		msgs[i] = msg.Name
	}

	return msgs
}
