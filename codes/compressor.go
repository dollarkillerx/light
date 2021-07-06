package codes

type CompressorType byte

const (
	GZIP CompressorType = iota
	RawData
	Snappy
)

// Compressor 压缩解压
type Compressor interface {
	Zip([]byte) ([]byte, error)
	Unzip([]byte) ([]byte, error)
}

type compressorManager struct {
	codes map[CompressorType]Compressor
}

var CompressorManager = &compressorManager{
	codes: map[CompressorType]Compressor{},
}

func (m *compressorManager) Register(key CompressorType, code Compressor) {
	m.codes[key] = code
}

func (m *compressorManager) Get(key CompressorType) (Compressor, bool) {
	code, ex := m.codes[key]
	return code, ex
}
