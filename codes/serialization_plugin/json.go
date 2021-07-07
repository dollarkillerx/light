package serialization_plugin

import (
	"encoding/json"
)

type JsonCode struct{}

func (j *JsonCode) Encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (j *JsonCode) Decode(data []byte, i interface{}) error {
	return json.Unmarshal(data, i)
}
