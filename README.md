# Actioncable-Client-Go

Actioncable Client Library for Go

# Usage

## Connections

```
u, _ := url.Parse("ws://example.org")
consumer := actioncable.CreateConsumer(u, nil)
consumer.Connect()
```

**Adding custom header**

```
u, _ := url.Parse("ws://example.org")

header := http.Header{}
header.Set("User-Agent", "Dummy")

opt := actioncable.NewConsumerOptions()
opt.SetHeader(header)

consumer := actioncable.CreateConsumer(u, opt)
consumer.Connect()
```

## Subscriptions

```
params := map[string]interface{}{
    "room":"Best Room"
}

id := actioncable.NewChannelIdentifier("ChatChannel", params)
subscription := consumer.Subscriptions.Create(id)
subscription.SetHandler(&ChatSubscriptionEventHandler{})


type ChatSubscriptionEventHandler {
    actioncable.SubscriptionEventHandler
}

func (h *ChatSubscriptionEventHandler) OnConnected(se *actioncable.SubscriptionEvent) {
    fmt.Println("on connected")
}

func (h *ChatSubscriptionEventHandler) OnDisconnected(se *actioncable.SubscriptionEvent) {
    fmt.Println("on disconnected")
}

func (h *ChatSubscriptionEventHandler) OnRejected(se *actioncable.SubscriptionEvent) {
    fmt.Println("on rejected")
}

func (h *ChatSubscriptionEventHandler) OnReceived(se *actioncable.SubscriptionEvent) {
    data := map[string]interface{}{}
    se.ReadJSON(&data)
    fmt.Println(data)
}
```
