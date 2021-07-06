package serialization_plugin

import (
	"github.com/dollarkillerx/light/codes"
	"github.com/vmihailenco/msgpack/v5"
)

func init() {
	codes.SerializationManager.Register(codes.CodeMsgPack, &MsgPackCode{})
}

type MsgPackCode struct{}

func (m *MsgPackCode) Encode(i interface{}) ([]byte, error) {
	return msgpack.Marshal(i)
}

func (m *MsgPackCode) Decode(data []byte, i interface{}) error {
	return msgpack.Unmarshal(data, i)
}
