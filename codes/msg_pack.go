package codes

import "github.com/vmihailenco/msgpack/v5"

func init() {
	Manager.register("msg_pack", &MsgPackCode{})
}

type MsgPackCode struct{}

func (m *MsgPackCode) Encode(i interface{}) ([]byte, error) {
	return msgpack.Marshal(i)
}

func (m *MsgPackCode) Decode(data []byte, i interface{}) error {
	return msgpack.Unmarshal(data, i)
}
