package actioncable

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConsumerConnect(t *testing.T) {
	url := TestWebsocketURL()
	if url == nil {
		fmt.Printf("Skip %s because os.Getend(\"TEST_WS\") is not set\n", t.Name())
		return
	}

	Convey("Given a valid URL", t, func() {
		consumer, _ := CreateConsumer(url, nil)

		Convey("consumer.connection should not be nil", func() {
			consumer.Connect()
			So(consumer.connection, ShouldNotBeNil)
		})

		Convey("welcome message should be received", func() {
			done := make(chan struct{})
			go func() {
				consumer.Connect()
				consumer.connection.waitUntilReady()
				done <- struct{}{}
			}()

			select {
			case <-done:
				So(consumer.connection.connectedAt, ShouldNotBeNil)
			case <-time.After(1 * time.Second):
				t.Fatal("Test didn't finish in time")
			}
		})

		Convey("with ConsumerOptions struct", func() {
			opt := NewConsumerOptions()

			header := http.Header{}
			header.Set("User-Agent", "Dummy")
			opt.SetHeader(&header)

			consumer, _ := CreateConsumer(url, opt)

			Convey("consumer.connection should not be nil", func() {
				consumer.Connect()
				So(consumer.connection, ShouldNotBeNil)
			})
		})
	})
}
