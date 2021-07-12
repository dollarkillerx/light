package utils

import (
	"bytes"
	"io"

	gzip "github.com/klauspost/pgzip"
)

// Unzip unzips data.
func Unzip(data []byte) ([]byte, error) {
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

// Zip zips data.
func Zip(data []byte) ([]byte, error) {
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
