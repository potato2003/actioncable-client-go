package actioncable

import (
	"encoding/json"
)

type SubscriptionEvent struct {
	Type       SubscriptionEventType
	Message    map[string]interface{}
	MessageRaw *json.RawMessage
}

type SubscriptionEventType string

const (
	Connected    SubscriptionEventType = SubscriptionEventType("connected")
	Disconnected SubscriptionEventType = SubscriptionEventType("disconnected")
	Rejected     SubscriptionEventType = SubscriptionEventType("rejected")
	Received     SubscriptionEventType = SubscriptionEventType("received")
)

func createSubscriptionEvent(eventType SubscriptionEventType, event *Event) *SubscriptionEvent {
	if event == nil {
		event = NewNilEvent()
	}
	e := newSubscriptionEvent(eventType, &event.Message)

	return e
}

func newSubscriptionEvent(eventType SubscriptionEventType, raw *json.RawMessage) *SubscriptionEvent {
	return &SubscriptionEvent{
		Type:       eventType,
		Message:    map[string]interface{}{},
		MessageRaw: raw,
	}
}

func (s *SubscriptionEvent) ReadJSON(v interface{}) error {
	b, err := json.Marshal(s.MessageRaw)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}
