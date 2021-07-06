package codes

type manager struct {
	codes map[Type]Code
}

type Code interface {
	Encode(i interface{}) ([]byte, error)
	Decode(data []byte, i interface{}) error
}

var Manager = &manager{
	codes: map[Type]Code{},
}

func (m *manager) register(key Type, code Code) {
	m.codes[key] = code
}

func (m *manager) Get(key Type) (Code, bool) {
	code, ex := m.codes[key]
	return code, ex
}

type Type byte

const (
	CodeJson Type = iota
	CodeMsgPack
)
