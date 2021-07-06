package codes

import "github.com/dollarkillerx/light/codes/serialization_plugin"

func init() {
	// 初始化所有插件
	serialization_plugin.InitSerialization()
}

type SerializationType byte

const (
	Json SerializationType = iota
	MsgPack
	Byte
)

type serializationManager struct {
	codes map[SerializationType]Serialization
}

type Serialization interface {
	Encode(i interface{}) ([]byte, error)
	Decode(data []byte, i interface{}) error
}

var SerializationManager = &serializationManager{
	codes: map[SerializationType]Serialization{},
}

func (m *serializationManager) Register(key SerializationType, code Serialization) {
	m.codes[key] = code
}

func (m *serializationManager) Get(key SerializationType) (Serialization, bool) {
	code, ex := m.codes[key]
	return code, ex
}
