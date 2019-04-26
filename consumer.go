package actioncable

import (
	"encoding/json"
	"net/url"
)

type Consumer struct {
	url           *url.URL
	Subscriptions *Subscriptions
	connection    *connection
	opts          *ConsumerOptions
}

func CreateConsumer(url *url.URL, opts *ConsumerOptions) (*Consumer, error) {
	if opts == nil {
		opts = NewConsumerOptions()
	}

	connection := newConnection(url.String())
	consumer := newConsumer(url, connection)
	consumer.opts = opts

	return consumer, nil
}

func newConsumer(url *url.URL, connection *connection) *Consumer {
	c := &Consumer{
		url:        url,
		connection: connection,
	}
	c.Subscriptions = newSubscriptions(c)

	return c
}

func (c *Consumer) send(data map[string]interface{}) error {
	logger.Debugf("send command: %+v", data)

	if identifier, ok := data["identifier"]; ok {
		copied := make(map[string]interface{})
		for k, v := range data {
			copied[k] = v
		}

		encodedIdentifer, err := json.Marshal(identifier)
		copied["identifier"] = string(encodedIdentifer)

		if err != nil {
			logger.Errorf("failed to send: %s", err.Error())
			return err
		}

		data = copied
	}

	return c.connection.send(data)
}

func (c *Consumer) Connect() {
	connection := newConnection(c.url.String())
	connection.consumer = c

	//if c.header != nil {
	//	connection.header = c.header
	//}

	connection.subscriptions = c.Subscriptions
	connection.start()

	c.connection = connection
}

func (c *Consumer) Disconnect() {
	c.connection.stop()
}
