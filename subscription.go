package actioncable

import (
	"encoding/json"
	"sync"
)

type Subscription struct {
	consumer       *Consumer
	Identifier     *ChannelIdentifier
	NotifyCh       chan *SubscriptionEvent
	stopHandleCh   chan chan struct{}
	handler        SubscriptionEventHandler
	lockForHandler *sync.Mutex
}

func newSubscription(consumer *Consumer, identifier *ChannelIdentifier) *Subscription {
	return &Subscription{
		consumer:       consumer,
		Identifier:     identifier,
		NotifyCh:       make(chan *SubscriptionEvent, 32),
		stopHandleCh:   make(chan chan struct{}, 1),
		lockForHandler: &sync.Mutex{},
	}
}

func (s *Subscription) Perform(action string, data map[string]interface{}) {
	copied := map[string]interface{}{}
	for k, v := range data {
		copied[k] = v
	}
	copied["action"] = action

	s.Send(copied)
}

func (s *Subscription) Send(data map[string]interface{}) {
	encodedData, _ := json.Marshal(data)

	s.consumer.send(map[string]interface{}{
		"command":    "message",
		"identifier": s.Identifier,
		"data":       string(encodedData),
	})
}

func (s *Subscription) SetHandler(h SubscriptionEventHandler) {
	s.lockForHandler.Lock()
	defer s.lockForHandler.Unlock()

	s.stopHandle()
	s.handler = h
	if h == nil {
		return
	}

	go func() {
		for {
			select {
			case se := <-s.NotifyCh:
				switch se.Type {
				case Connected:
					h.OnConnected(se)
				case Disconnected:
					h.OnDisconnected(se)
				case Rejected:
					h.OnRejected(se)
				case Received:
					h.OnReceived(se)
				default:
					logger.Warnf("unknown subscription event: %v", se)
				}
			case doneCh := <-s.stopHandleCh:
				doneCh <- struct{}{}
				return
			}
		}
	}()
}

func (s *Subscription) stopHandle() {
	if s.handler == nil {
		return
	}

	doneCh := make(chan struct{}, 1)
	defer close(doneCh)

	s.stopHandleCh <- doneCh
	<-doneCh
}

func (s *Subscription) Unsubscribe() {
	s.lockForHandler.Lock()
	defer s.lockForHandler.Unlock()

	s.stopHandle()
	s.consumer.Subscriptions.remove(s)
}
