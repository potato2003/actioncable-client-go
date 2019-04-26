package actioncable

type SubscriptionEventHandler interface {
	OnConnected(*SubscriptionEvent)
	OnDisconnected(*SubscriptionEvent)
	OnRejected(*SubscriptionEvent)
	OnReceived(*SubscriptionEvent)
}
