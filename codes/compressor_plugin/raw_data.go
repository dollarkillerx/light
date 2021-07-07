package compressor_plugin

type RawData struct{}

func (r *RawData) Zip(bytes []byte) ([]byte, error) {
	return bytes, nil
}

func (r *RawData) Unzip(bytes []byte) ([]byte, error) {
	return bytes, nil
}
