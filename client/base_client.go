package client

import (
	"log"
	"net"
	"time"

	"github.com/dollarkillerx/light"
	"github.com/dollarkillerx/light/codes"
	"github.com/dollarkillerx/light/cryptology"
	"github.com/dollarkillerx/light/pkg"
	"github.com/dollarkillerx/light/protocol"
	"github.com/dollarkillerx/light/utils"
)

type BaseClient struct {
	conn       net.Conn
	options    *Options
	serverName string

	aesKey        []byte
	serialization codes.Serialization
	compressor    codes.Compressor
}

func newBaseClient(serverName string, con net.Conn, options *Options) (*BaseClient, error) {

	serialization, ex := codes.SerializationManager.Get(options.serializationType)
	if !ex {
		return nil, pkg.ErrSerialization404
	}

	compressor, ex := codes.CompressorManager.Get(options.compressorType)
	if !ex {
		return nil, pkg.ErrCompressor404
	}

	return &BaseClient{
		serverName:    serverName,
		conn:          con,
		options:       options,
		serialization: serialization,
		compressor:    compressor,
		aesKey:        options.aesKey,
	}, nil
}

func (b *BaseClient) Call(ctx *light.Context, serviceMethod string, request interface{}, response interface{}) (err error) {
	metaData := ctx.GetMetaData()
	defer func() {
		if err := recover(); err != nil {
			utils.PrintStack()
			log.Println("Recover Err: ", err)
		}
	}() // 网络不可靠

	// 1. 构造请求
	// 1.1 序列化
	serviceNameByte := []byte(b.serverName)
	serviceMethodByte := []byte(serviceMethod)
	var metaDataBytes []byte
	var requestBytes []byte
	metaDataBytes, err = b.serialization.Encode(metaData)
	if err != nil {
		return err
	}
	requestBytes, err = b.serialization.Encode(request)
	if err != nil {
		return err
	}
	// 1.2 加密
	metaDataBytes, err = cryptology.AESEncrypt(metaDataBytes, b.aesKey)
	if err != nil {
		return err
	}

	requestBytes, err = cryptology.AESEncrypt(requestBytes, b.aesKey)
	if err != nil {
		return err
	}
	// 1.3 压缩
	metaDataBytes, err = b.compressor.Zip(metaDataBytes)
	if err != nil {
		return err
	}

	requestBytes, err = b.compressor.Zip(requestBytes)
	if err != nil {
		return err
	}
	// 1.4 封装消息
	message, err := protocol.EncodeMessage(serviceNameByte, serviceMethodByte, metaDataBytes, byte(protocol.Request), byte(b.options.compressorType), byte(b.options.serializationType), requestBytes)
	if err != nil {
		return err
	}
	// 2. 发送消息
	if b.options.writeTimeout > 0 {
		now := time.Now()
		b.conn.SetWriteDeadline(now.Add(b.options.writeTimeout))
	}
	_, err = b.conn.Write(message)
	if err != nil {
		return err
	}
	// 3.封装回执

	return nil
}
