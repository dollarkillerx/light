package serialization_plugin

import (
	"encoding/json"

	"github.com/dollarkillerx/light/codes"
)

func init() {
	codes.SerializationManager.Register(codes.Json, &JsonCode{})
}

type JsonCode struct{}

func (j *JsonCode) Encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (j *JsonCode) Decode(data []byte, i interface{}) error {
	return json.Unmarshal(data, i)
}
