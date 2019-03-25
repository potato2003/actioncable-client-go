package actioncable

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jpillora/backoff"
)

const (
	// use the same value as actioncable.js
	// 1 ping is 2 sec interval, so detect stale when 2 ping missing.
	DEFAULT_STALE_THRESHOLD time.Duration = 6 * time.Second
)

var (
	staleError error = errors.New("connection is stale")
)

type connection struct {
	url           string
	consumer      *Consumer
	subscriptions *Subscriptions
	disconnected  bool
	dialer        *websocket.Dialer
	ws            *websocket.Conn
	header        *http.Header
	recieveCh     chan Event
	isReady       bool
	readyCh       chan struct{}
	pingedAt      time.Time
	connectedAt   time.Time
	lockForSend   *sync.Mutex
}

func newConnection(url string) *connection {
	return &connection{
		url:          url,
		disconnected: true,
		dialer: &websocket.Dialer{
			HandshakeTimeout: 5 * time.Second,
		},
		header:      &http.Header{},
		recieveCh:   make(chan Event, 1),
		isReady:     false,
		readyCh:     make(chan struct{}, 1),
		lockForSend: new(sync.Mutex),
	}
}

func (c *connection) start() {
	go c.connectionLoop()
	c.waitUntilReady()
}

func (c *connection) stop() error {
	log.Println("close")
	return c.ws.Close()
}

func (c *connection) connectionLoop() {
	b := backoff.Backoff{
		Min:    100 * time.Millisecond,
		Max:    5 * time.Second,
		Factor: 3,
		Jitter: true,
	}
	defer c.stop()

	/*
		Connect to actioncable server.
		If it fails, will retry to connect with backoff.

		After connection establishment, continuously receive messages
		from server, and notify to event handler.
	*/
	for {
		c.isReady = false

		err := c.establishConnection()
		if err != nil {
			logger.Infof("failed to connect, %s\n", err)
		} else {
			b.Reset()
			c.eventHandlerLoop()
		}

		select {
		case <-time.After(b.Duration()): // exponential backoff
			logger.Infof("reconnecting")
		}
	}

	return
}

func (c *connection) establishConnection() error {
	ws, _, err := c.dialer.Dial(c.url, *c.header)
	if err != nil {
		return err
	}

	c.ws = ws

	return nil
}

func (c *connection) eventHandlerLoop() {
	for {
		event, err := c.receive() // wait max `DEFAULT_STALE_THRESHOLD` sec until recive new message

		if err != nil {
			se := createSubscriptionEvent("disconnected", nil)
			c.subscriptions.notifyAll(se)
			logger.Errorf("%s\n", err)
			return // attempts to reconnect
		}

		switch event.Type {
		case "ping": // receive every 2 sec
			c.recordPing()
		case "welcome": // receive after establish connection
			logger.Debug("Received welcome message")
			c.recordConnect()
			c.subscriptions.reload()
			c.ready()
		case "confirm_subscription": // response of subscribe request
			logger.Debug("Received confirm_subscription message")
			se := createSubscriptionEvent(Connected, event)
			c.subscriptions.notify(event.Identifier, se)
		case "rejection":
			c.subscriptions.reject(event.Identifier)
		case "disconnect":
			// close
			se := createSubscriptionEvent(Disconnected, nil)
			c.subscriptions.notifyAll(se)
			return
		default:
			se := createSubscriptionEvent(Received, event)
			c.subscriptions.notify(event.Identifier, se)
		}
	}
}

func (c *connection) receive() (*Event, error) {
	// Note to not call receive() in concurrent, it's not working.
	// because, gorilla/websocket does not support concurrent reading.
	// see: https://godoc.org/github.com/gorilla/websocket#hdr-Concurrency
	ch := make(chan *Event)
	errCh := make(chan error)

	go func() {
		event := &Event{}
		if err := c.ws.ReadJSON(event); err != nil {
			errCh <- err
		}

		ch <- event
	}()

	// using timeout for checking stale of ac server.
	select {
	case event := <-ch:
		return event, nil
	case err := <-errCh:
		return nil, err
	case <-time.After(DEFAULT_STALE_THRESHOLD):
		log.Printf("connection is stale")
		return nil, staleError
	}

	return nil, nil
}

func (c *connection) send(data map[string]interface{}) error {
	// gorilla/websocket does not support concurrent writing.
	// see: https://godoc.org/github.com/gorilla/websocket#hdr-Concurrency
	c.lockForSend.Lock()
	defer c.lockForSend.Unlock()

	err := c.ws.WriteJSON(data)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *connection) ready() {
	c.isReady = true
	clearCh(c.readyCh)
	c.readyCh <- struct{}{}
}

func (c *connection) waitUntilReady() {
	if !c.isReady {
		_ = <-c.readyCh
	}
}

func (c *connection) recordPing() {
	c.pingedAt = time.Now()
}

func (c *connection) recordConnect() {
	c.recordPing()
	c.connectedAt = time.Now()
}

func clearCh(ch chan struct{}) {
	for len(ch) > 0 {
		_ = <-ch
	}
}
