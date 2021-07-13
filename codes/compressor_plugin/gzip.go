package compressor_plugin

import (
	"bytes"
	"io"

	gzip "github.com/klauspost/pgzip"
)

type GZIP struct{}

func (G *GZIP) Zip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, err := gw.Write(data)
	if err != nil {
		return nil, err
	}
	err = gw.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (G *GZIP) Unzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.Write(data)
	if err != nil {
		return nil, err
	}
	reader, err := gzip.NewReader(&buf)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	buffer := bytes.NewBuffer(nil)
	_, err = io.Copy(buffer, reader)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
