[![Build Status](https://travis-ci.org/potato2003/actioncable-client-go.svg?branch=master)](https://travis-ci.org/potato2003/actioncable-client-go)

# Actioncable-Client-Go

Actioncable Client Library for Go.

# Requirements

Go v1.10 or later.

# Usage

## Connections

```go
u, _ := url.Parse("ws://example.org")
consumer := actioncable.CreateConsumer(u, nil)
consumer.Connect()
```

You can also customize the HTTP header by using `NewConsumerOptions`.

```go
u, _ := url.Parse("ws://example.org")

header := http.Header{}
header.Set("User-Agent", "Dummy")

opt := actioncable.NewConsumerOptions()
opt.SetHeader(&header)

consumer := actioncable.CreateConsumer(u, opt)
consumer.Connect()
```

## Subscriptions

```go
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

You can send data from the client side to the server side. For example:

```go
// id := actioncable.NewChannelIdentifier("ChatChannel", params)
// subscription := consumer.Subscriptions.Create(id)

func (h *ChatSubscriptionEventHandler) OnConnected(se *actioncable.SubscriptionEvent) {
    data := map[string]interface{}{
        "fieldA":"valueA",
        "fieldB":"valueB",
    }

    // # Calls `ChatChannel#appear(data)` on the server.
    subscription.Peform("appear", data)
}
```

# License

MIT License

# Test

```
go test -v ./...
```

## Integration Test with ActionCable Server

```bash
(cd test_rails_server; ./bin/setup; bundle exec rails -p 3000 -d) # start actioncable server
TEST_WS="ws://localhost:3000/cable" go test -v ./...
kill $(cat ./test_rails_server/tmp/pids/server.pid) # stop actioncable server

```
