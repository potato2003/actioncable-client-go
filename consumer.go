package actioncable

import (
	"encoding/json"
	"fmt"
)

type Consumer struct {
	url           string
	Subscriptions *Subscriptions
	connection    *connection
}

func CreateConsumer(url string) (*Consumer, error) {
	connection := newConnection(url)
	consumer := newConsumer(url, connection)
	consumer.connection = connection

	connection.consumer = consumer
	connection.subscriptions = consumer.Subscriptions

	connection.start()

	return consumer, nil
}

func newConsumer(url string, connection *connection) *Consumer {
	c := &Consumer{
		url:        url,
		connection: connection,
	}
	c.Subscriptions = newSubscriptions(c)

	return c
}

func (c *Consumer) send(data map[string]interface{}) error {
	logger.Debugf("send: %+v", data)

	// identifier must be string of object
	if identifier, ok := data["identifier"]; ok {
		copied := make(map[string]interface{})
		for k, v := range data {
			copied[k] = v
		}

		encodedIdentifer, err := json.Marshal(identifier)
		copied["identifier"] = string(encodedIdentifer)

		if err != nil {
			fmt.Println(err)
			return err
		}

		data = copied
	}

	return c.connection.send(data)
}

func (c *Consumer) connect() {
}

func (c *Consumer) disconnect() {
}

func (c *Consumer) ensureActiveConnection() {
	// c.connection.waitUntilReady()
}
