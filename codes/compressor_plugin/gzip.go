package compressor_plugin

import (
	"github.com/dollarkillerx/light/utils"
)

type GZIP struct{}

func (G *GZIP) Zip(bytes []byte) ([]byte, error) {
	zip, err := utils.Zip(bytes)
	if err != nil {
		panic(err)
		return nil, err
	}
	return zip, err
}

func (G *GZIP) Unzip(bytes []byte) ([]byte, error) {
	unzip, err := utils.Unzip(bytes)
	if err != nil {
		panic(err)
		return nil, err
	}

	return unzip, err
}
