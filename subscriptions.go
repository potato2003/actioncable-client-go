package actioncable

import (
	"log"
)

type Subscriptions struct {
	consumer      *Consumer
	subscriptions []*Subscription
}

func newSubscriptions(consumer *Consumer) *Subscriptions {
	return &Subscriptions{
		consumer:      consumer,
		subscriptions: []*Subscription{},
	}
}

func (s *Subscriptions) Create(channelIdentifier *ChannelIdentifier) (*Subscription, error) {
	subscription, err := s.add(newSubscription(s.consumer, channelIdentifier))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return subscription, nil
}

func (s *Subscriptions) add(subscription *Subscription) (*Subscription, error) {
	s.subscriptions = append(s.subscriptions, subscription)
	s.sendCommand(subscription, "subscribe")

	return subscription, nil
}

func (s *Subscriptions) remove(subscription *Subscription) *Subscription {
	s.forget(subscription)
	s.sendCommand(subscription, "unsubscribe")

	return subscription
}

func (s *Subscriptions) reject(identifier *ChannelIdentifier) []*Subscription {
	matches := s.findAll(identifier)

	for _, subscription := range matches {
		s.forget(subscription)

		rejectedEvent := createSubscriptionEvent(Rejected, nil)
		s.notify(identifier, rejectedEvent)
	}

	return matches
}

func (s *Subscriptions) forget(subscription *Subscription) {
	s.subscriptions = s.filter(subscription.Identifier)
}

func (s *Subscriptions) findAll(identifier *ChannelIdentifier) []*Subscription {
	result := make([]*Subscription, 0, len(s.subscriptions))

	for _, subscription := range s.subscriptions {
		if identifier.Equals(subscription.Identifier) {
			result = append(result, subscription)
		}
	}

	return result
}

func (s *Subscriptions) filter(identifier *ChannelIdentifier) []*Subscription {
	result := make([]*Subscription, 0, len(s.subscriptions))

	for _, subscription := range s.subscriptions {
		if !identifier.Equals(subscription.Identifier) {
			result = append(result, subscription)
		}
	}

	return result
}

func (s *Subscriptions) reload() {
	logger.Debug("reloading Subscriptions")

	for _, subscription := range s.subscriptions {
		s.sendCommand(subscription, "subscribe")
	}
}

func (s *Subscriptions) notifyAll(event *SubscriptionEvent) {
	for _, subscription := range s.subscriptions {
		log.Printf("notify %+v", event)
		subscription.NotifyCh <- event
	}
}

func (s *Subscriptions) notify(identifier *ChannelIdentifier, event *SubscriptionEvent) {
	for _, subscription := range s.findAll(identifier) {
		log.Printf("notify %+v", event)
		subscription.NotifyCh <- event
	}
}

func (s *Subscriptions) sendCommand(subscription *Subscription, command string) {
	data := map[string]interface{}{
		"command":    command,
		"identifier": subscription.Identifier,
	}
	s.consumer.send(data)

	return
}
