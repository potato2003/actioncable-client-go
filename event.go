package actioncable

import (
	"encoding/json"
)

type Event struct {
	Type string `json:"type"`
	Message    json.RawMessage    `json:"message"`
	Reason     json.RawMessage    `json:"reason"`
	Reconnect  json.RawMessage    `json:"reconnect"`
	Data       json.RawMessage    `json:"data"`
	Identifier *ChannelIdentifier `json:"identifier"`
}

func NewNilEvent() *Event {
	return &Event{}
}
