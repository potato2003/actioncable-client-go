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
	if os.Getenv("TEST_WS") == "" {
		fmt.Println("skip TestCreateConsumer")
		return
	}

	url, _ := url.Parse("ws://localhost:3000/cable")
	if strings.HasSuffix(os.Getenv("TEST_WS"), "ws") {
		url, _ = url.Parse(os.Getenv("TEST_WS"))
	}

	consumer, err := CreateConsumer(url, nil)

	channelIdentifier := NewChannelIdentifier("AgentChannel", map[string]interface{}{"room": "BestRoom"})
	//time.Sleep(100 * time.Millisecond)
	subscription, _ := consumer.Subscriptions.Create(channelIdentifier)

	go func() {
		for {
			select {
			case event := <-subscription.NotifyCh:
				switch event.Type {
				case "connected":
					fmt.Printf("connected: %+v", event)
				case "received":
					fmt.Printf("received: %+v", event)
					subscription.Perform("hogehoge", map[string]interface{}{"data": "TEST"})
				default:
					fmt.Printf("subscribe event: %+v", event)
				}
			}
		}
	}()

	if err != nil {
		t.Fatalf("bad: %s", err)
	}
	if consumer.connection == nil {
		t.Fatalf("bad: consumer.connection is nil")
	}

	time.Sleep(1 * time.Second)
}
