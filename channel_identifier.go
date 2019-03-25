package actioncable

import (
	"encoding/json"
	"reflect"
)

type ChannelIdentifier struct {
	channelName    string
	params         map[string]interface{}
	marshaledValue []byte
}

func NewChannelIdentifier(channelName string, params map[string]interface{}) *ChannelIdentifier {
	if params == nil {
		params = make(map[string]interface{})
	}

	id := &ChannelIdentifier{
		channelName:channelName,
		params:params,
	}
	m, _ := id.MarshalJSON()
	id.marshaledValue = m

	return id
}

// Implements json.Marshaler#MarshalJSON()
func (c *ChannelIdentifier) MarshalJSON() ([]byte, error) {
	copied := make(map[string]interface{})
	for k, v := range c.params {
		copied[k] = v
	}

	copied["channel"] = c.channelName

	return json.Marshal(copied)
}

// Implements json.Marshaler#UnmarshalJSON()
func (c *ChannelIdentifier) UnmarshalJSON(doubleEncodedData []byte) error {
	str := ""
	if err := json.Unmarshal(doubleEncodedData, &str); err != nil {
		return err
	}

	raw := json.RawMessage(str)
	params := make(map[string]interface{})
	if err := json.Unmarshal(raw, &params); err != nil {
		return err
	}

	if str, ok := params["channel"].(string); ok {
		c.channelName = str
	}
	delete(params, "channel")

	c.params = params
	c.marshaledValue = []byte(str)

	return nil
}

func (self *ChannelIdentifier) Equals(other *ChannelIdentifier) bool {
	// Comapre serialized value of self and other value.
	// because golang does not support comparisons between Structs,
	// so does not work when ChannelIdentifer#params including Struct Type.
	return reflect.DeepEqual(self.marshaledValue, other.marshaledValue)
}
