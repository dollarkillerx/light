package compressor_plugin

import (
	"bytes"
	"io/ioutil"

	"github.com/golang/snappy"
)

type Snappy struct{}

func (s *Snappy) Zip(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}

	var buffer bytes.Buffer
	writer := snappy.NewBufferedWriter(&buffer)
	_, err := writer.Write(data)
	if err != nil {
		writer.Close()
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (s *Snappy) Unzip(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}

	reader := snappy.NewReader(bytes.NewReader(data))
	out, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return out, err
}
