package serialization_plugin

import (
	"github.com/vmihailenco/msgpack/v5"
)

type MsgPackCode struct{}

func (m *MsgPackCode) Encode(i interface{}) ([]byte, error) {
	return msgpack.Marshal(i)
}

func (m *MsgPackCode) Decode(data []byte, i interface{}) error {
	return msgpack.Unmarshal(data, i)
}
