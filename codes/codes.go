package codes

type manager struct {
	codes map[string]Code
}

type Code interface {
	Encode(i interface{}) ([]byte, error)
	Decode(data []byte, i interface{}) error
}

var Manager = &manager{
	codes: map[string]Code{},
}

func (m *manager) register(key string, code Code) {
	m.codes[key] = code
}

func (m *manager) Get(key string) (Code, bool) {
	code, ex := m.codes[key]
	return code, ex
}
