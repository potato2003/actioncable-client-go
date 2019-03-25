package actioncable

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestCreateConsumer(t *testing.T) {
	if os.Getenv("TEST_WS") == "" {
		fmt.Println("skip TestCreateConsumer")
		return
	}

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGCONT)
		buf := make([]byte, 1<<20)
		for {
			sig := <-sigs
			fmt.Printf("%+v", sig)

			stacklen := runtime.Stack(buf, true)
			log.Printf("=== received SIGCOUNT ===\n*** goroutine dump...\n%s\n*** end\n", buf[:stacklen])
		}
	}()

	url := "ws://localhost:3000/cable"
	if strings.HasSuffix(os.Getenv("TEST_WS"), "ws") {
		url = os.Getenv("TEST_WS")
	}

	consumer, err := CreateConsumer(url)

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

	time.Sleep(50 * time.Second)
}

func TestSlice(t *testing.T) {
	s := []*Subscription{}

	s = append(s, &Subscription{})
	fmt.Printf("%p: %d\n", &s, len(s))

	s = append(s, &Subscription{})
	fmt.Printf("%p: %d\n", &s, len(s))

	b := append(s, &Subscription{})
	fmt.Printf("%p: %d\n", &b, len(b))
}
