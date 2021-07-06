package compressor_plugin

import "github.com/dollarkillerx/light/codes"

func init() {
	codes.CompressorManager.Register(codes.RawData, &RawData{})
}

type RawData struct{}

func (r *RawData) Zip(bytes []byte) ([]byte, error) {
	return bytes, nil
}

func (r *RawData) Unzip(bytes []byte) ([]byte, error) {
	return bytes, nil
}
