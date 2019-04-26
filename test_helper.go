package actioncable

import (
	"net/url"
	"os"
	"strings"
)

func TestWebsocketURL() *url.URL {
	if os.Getenv("TEST_WS") == "" {
		return nil
	}

	url, _ := url.Parse("ws://localhost:3000/cable")
	if strings.HasSuffix(os.Getenv("TEST_WS"), "ws") {
		url, _ = url.Parse(os.Getenv("TEST_WS"))
	}

	return url
}
