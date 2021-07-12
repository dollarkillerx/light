package compressor_plugin

import (
	"github.com/dollarkillerx/light/utils"
)

type GZIP struct{}

func (G *GZIP) Zip(bytes []byte) ([]byte, error) {
	return utils.Zip(bytes)
}

func (G *GZIP) Unzip(bytes []byte) ([]byte, error) {
	return utils.Unzip(bytes)
}
