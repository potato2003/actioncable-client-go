package actioncable

import (
	"encoding/json"
	"testing"
)

func TestNewChannelIdentifier(t *testing.T) {
	id := NewChannelIdentifier("ChatRoom", nil)
	if id.channelName != "ChatRoom" {
		t.Fatalf("bad: %s", id.channelName)
	}
}

func TestMarshalJSON(t *testing.T) {
	// no channel params
	id := NewChannelIdentifier("ChatRoom", nil)
	b, err := json.Marshal(id)
	if err != nil {
		t.Fatalf("bad: %s", err)
	}

	actual := string(b)
	expected := `{"channel":"ChatRoom"}`

	if expected != actual {
		t.Fatalf("bad: %s", actual)
	}

	// using channel params (1)
	id = NewChannelIdentifier("ChatRoom", map[string]interface{}{
		"Room": "BestRoom",
	})
	b, err = json.Marshal(id)
	if err != nil {
		t.Fatalf("bad: %s", err)
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("bad: %s", err)
	}

	if data["channel"] != "ChatRoom" {
		t.Fatalf("bad: %s", data)
	}
	if data["Room"] != "BestRoom" {
		t.Fatalf("bad: %s", data)
	}
}
