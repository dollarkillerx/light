package codes

import "encoding/json"

func init() {
	Manager.register(CodeJson, &JsonCode{})
}

type JsonCode struct{}

func (j *JsonCode) Encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (j *JsonCode) Decode(data []byte, i interface{}) error {
	return json.Unmarshal(data, i)
}
