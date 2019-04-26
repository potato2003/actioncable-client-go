package actioncable

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCreateConsumer(t *testing.T) {
	url := TestWebsocketURL()
	if url == nil {
		fmt.Println("skip TestCreateConsumer")
		return
	}

	consumer, err := CreateConsumer(url, nil)
	if err != nil {
		t.Fatalf("bad: %s", err)
	}

	consumer.Connect()
	if consumer.connection == nil {
		t.Fatalf("bad: consumer.connection is nil")
	}

	consumer.Disconnect()
}

// TODO move to subscription_test
func TestCreateConsumer2(t *testing.T) {
	if os.Getenv("TEST_WS") == "" {
		fmt.Println("skip TestCreateConsumer")
		return
	}

	url, _ := url.Parse("ws://localhost:3000/cable")
	if strings.HasSuffix(os.Getenv("TEST_WS"), "ws") {
		url, _ = url.Parse(os.Getenv("TEST_WS"))
	}

	consumer, err := CreateConsumer(url, nil)
	consumer.Connect()

	channelIdentifier := NewChannelIdentifier("ChatChannel", map[string]interface{}{"room": "BestRoom"})
	subscription, err := consumer.Subscriptions.Create(channelIdentifier)

	if err != nil {
		t.Fatalf("bad: %s", err)
	}
	if consumer.connection == nil {
		t.Fatalf("bad: consumer.connection is nil")
	}

	handler := NewTestChatSubscriptionEventHandler()
	subscription.SetHandler(handler)

	select {
	case <-handler.ConnectedEvents:
	case <-time.After(1 * time.Second):
		t.Fatal("bad: connection timeout")
	}

	// data := map[string]interface{}{}
	// se.ReadJSON(&data)
}

type TestChatSubscriptionEventHandler struct {
	NumConnected    int
	NumDisconnected int
	NumRejected     int
	NumReceived     int

	ConnectedEvents    chan *SubscriptionEvent
	DisconnectedEvents chan *SubscriptionEvent
	RejectedEvents     chan *SubscriptionEvent
	ReceivedEvents     chan *SubscriptionEvent

	SubscriptionEventHandler
}

func NewTestChatSubscriptionEventHandler() *TestChatSubscriptionEventHandler {
	h := TestChatSubscriptionEventHandler{}

	h.NumConnected = 0
	h.NumDisconnected = 0
	h.NumRejected = 0
	h.NumReceived = 0

	h.ConnectedEvents = make(chan *SubscriptionEvent, 128)
	h.DisconnectedEvents = make(chan *SubscriptionEvent, 128)
	h.RejectedEvents = make(chan *SubscriptionEvent, 128)
	h.ReceivedEvents = make(chan *SubscriptionEvent, 128)

	return &h
}

func (h *TestChatSubscriptionEventHandler) Reset() {
	h.NumConnected = 0
	h.NumDisconnected = 0
	h.NumRejected = 0
	h.NumReceived = 0

	h.ConnectedEvents = make(chan *SubscriptionEvent, 128)
	h.DisconnectedEvents = make(chan *SubscriptionEvent, 128)
	h.RejectedEvents = make(chan *SubscriptionEvent, 128)
	h.ReceivedEvents = make(chan *SubscriptionEvent, 128)
}

func (h *TestChatSubscriptionEventHandler) Close() {
	close(h.ConnectedEvents)
	close(h.DisconnectedEvents)
	close(h.RejectedEvents)
	close(h.ReceivedEvents)
}

func (h *TestChatSubscriptionEventHandler) OnConnected(se *SubscriptionEvent) {
	h.NumConnected++
	h.ConnectedEvents <- se
}

func (h *TestChatSubscriptionEventHandler) OnDisconnected(se *SubscriptionEvent) {
	h.NumDisconnected++
	h.DisconnectedEvents <- se
}

func (h *TestChatSubscriptionEventHandler) OnRejected(se *SubscriptionEvent) {
	h.NumReceived++
	h.RejectedEvents <- se
}

func (h *TestChatSubscriptionEventHandler) OnReceived(se *SubscriptionEvent) {
	h.NumReceived++
	h.ReceivedEvents <- se
}
