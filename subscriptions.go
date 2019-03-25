package actioncable

import (
	"log"
)

type Subscriptions struct {
	consumer *Consumer
	subscriptions []*Subscription
}

func newSubscriptions(consumer *Consumer) *Subscriptions {
	return &Subscriptions{
		consumer: consumer,
		subscriptions: []*Subscription{},
	}
}

func (s *Subscriptions) Create(channelIdentifier *ChannelIdentifier) (*Subscription, error) {
	subscription, err := s.add(newSubscription(s.consumer, channelIdentifier))
	if err != nil {
		log.Println(err)
	}

	return subscription, nil
}

func (s *Subscriptions) add(subscription *Subscription) (*Subscription, error) {
	s.subscriptions = append(s.subscriptions, subscription)
	s.consumer.ensureActiveConnection()
	s.sendCommand(subscription, "subscribe")
	// notify initialize

	return subscription, nil
}

func (s *Subscriptions) remove(subscription *Subscription) *Subscription {
	s.forget(subscription)
	s.sendCommand(subscription, "unsubscribe")

	return subscription
}

func (s *Subscriptions) reject(identifier *ChannelIdentifier) []*Subscription {
	matches := s.findAll(identifier)
	s.forgetAll(matches)

	for _, subscription := range matches {
		_ = subscription
		// notify initialize
	}

	return matches
}

func (s *Subscriptions) forget(subscription *Subscription) []*Subscription {
	_ = s.findAll(subscription.Identifier)

	return nil
}

func (s *Subscriptions) forgetAll(subscriptions []*Subscription) []*Subscription {
	//_ = s.findAll(identifier)

	return nil
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

func (s *Subscriptions) reload() {
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
		"command":command,
		"identifier":subscription.Identifier,
	}
	s.consumer.send(data)

	return
}
