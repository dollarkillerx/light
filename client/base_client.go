package client

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/dollarkillerx/light"
	"github.com/dollarkillerx/light/codes"
	"github.com/dollarkillerx/light/cryptology"
	"github.com/dollarkillerx/light/pkg"
	"github.com/dollarkillerx/light/protocol"
	"github.com/dollarkillerx/light/transport"
	"github.com/dollarkillerx/light/utils"
	"github.com/pkg/errors"
)

type BaseClient struct {
	conn       net.Conn
	options    *Options
	serverName string

	aesKey        []byte
	serialization codes.Serialization
	compressor    codes.Compressor

	respInterMap map[string]*respMessage
	respInterRM  sync.RWMutex
	writeMu      sync.Mutex
}

type respMessage struct {
	response interface{}
	ctx      *light.Context
	respChan chan error
}

func newBaseClient(serverName string, options *Options) (*BaseClient, error) {
	service := options.loadBalancing.GetService()
	con, err := transport.Client.Gen(service.Protocol, service.Addr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	serialization, ex := codes.SerializationManager.Get(options.serializationType)
	if !ex {
		return nil, pkg.ErrSerialization404
	}

	compressor, ex := codes.CompressorManager.Get(options.compressorType)
	if !ex {
		return nil, pkg.ErrCompressor404
	}

	bc := &BaseClient{
		serverName:    serverName,
		conn:          con,
		options:       options,
		serialization: serialization,
		compressor:    compressor,
		aesKey:        options.aesKey,
		respInterMap:  map[string]*respMessage{},
	}

	go bc.heartBeat()
	go bc.processMessageManager()

	return bc, nil
}

func (b *BaseClient) Call(ctx *light.Context, serviceMethod string, request interface{}, response interface{}) (err error) {
	defer func() {
		if err := recover(); err != nil {
			utils.PrintStack()
			log.Println("Recover Err: ", err)
		}
	}() // 网络不可靠

	// call
	var magic string
	respChan := make(chan error, 0)
	magic, err = b.call(ctx, serviceMethod, request, response, respChan)
	if err != nil {
		return err
	}
	defer func() {
		// 回收内存
		b.respInterRM.Lock()
		delete(b.respInterMap, magic)
		b.respInterRM.Unlock()
	}()

	timeout := ctx.GetTimeout()
	if timeout <= 0 {
		timeout = b.options.readTimeout
	}

	// 调度
	select {
	case err, ex := <-respChan:
		if !ex {
			return nil
		}

		return err
	case <-time.After(timeout):
		return pkg.ErrTimeout
	}
}

func (b *BaseClient) call(ctx *light.Context, serviceMethod string, request interface{}, response interface{}, respChan chan error) (magic string, err error) {
	metaData := ctx.GetMetaData()

	// 1. 构造请求
	// 1.1 序列化
	serviceNameByte := []byte(b.serverName)
	serviceMethodByte := []byte(serviceMethod)
	var metaDataBytes []byte
	var requestBytes []byte
	metaDataBytes, err = b.serialization.Encode(metaData)
	if err != nil {
		return "", err
	}
	requestBytes, err = b.serialization.Encode(request)
	if err != nil {
		return "", err
	}

	// 1.2 加密
	metaDataBytes, err = cryptology.AESEncrypt(b.aesKey, metaDataBytes)
	if err != nil {
		return "", err
	}

	requestBytes, err = cryptology.AESEncrypt(b.aesKey, requestBytes)
	if err != nil {
		return "", err
	}
	// 1.3 压缩
	metaDataBytes, err = b.compressor.Zip(metaDataBytes)
	if err != nil {
		return "", err
	}

	requestBytes, err = b.compressor.Zip(requestBytes)
	if err != nil {
		return "", err
	}
	// 1.4 封装消息
	magic, message, err := protocol.EncodeMessage("", serviceNameByte, serviceMethodByte, metaDataBytes, byte(protocol.Request), byte(b.options.compressorType), byte(b.options.serializationType), requestBytes)
	if err != nil {
		return "", err
	}
	// 2. 发送消息
	if b.options.writeTimeout > 0 {
		now := time.Now()
		timeout := ctx.GetTimeout() // 如果ctx 存在设置 则采用 返之使用默认配置
		if timeout > 0 {
			b.conn.SetWriteDeadline(now.Add(timeout))
		} else {
			b.conn.SetWriteDeadline(now.Add(b.options.writeTimeout))
		}
	}

	// 有点暴力呀 直接上锁
	b.writeMu.Lock()
	b.writeMu.Unlock()
	_, err = b.conn.Write(message)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// 写MAP
	b.respInterRM.Lock()
	b.respInterMap[magic] = &respMessage{
		response: response,
		ctx:      ctx,
		respChan: respChan,
	}
	b.respInterRM.Unlock()

	return magic, nil
}

func (b *BaseClient) heartBeat() {
	for {
		_, i, err := protocol.EncodeMessage("x", []byte(""), []byte(""), []byte(""), byte(protocol.HeartBeat), byte(b.options.compressorType), byte(b.options.serializationType), []byte(""))
		if err != nil {
			log.Println(err)
			break
		}

		b.writeMu.Lock()
		b.conn.Write(i)
		b.writeMu.Unlock()
		time.Sleep(b.options.heartBeat)
	}
}

func (b *BaseClient) processMessageManager() {
	for {
		magic, respChan, err := b.processMessage()
		if err == nil && magic == "" {
			continue
		}

		if err != nil && magic == "" {
			log.Println(err)
			// 重zhi
			break
		}

		if err != nil && magic != "" && respChan != nil {
			respChan <- err
		}

		if err == nil && magic != "" && respChan != nil {
			close(respChan)
		}
	}
}

func (b *BaseClient) processMessage() (magic string, respChan chan error, err error) {
	// 3.封装回执
	now := time.Now()
	b.conn.SetReadDeadline(now.Add(b.options.readTimeout))

	proto := protocol.NewProtocol()
	msg, err := proto.IODecode(b.conn)
	if err != nil {
		return "", nil, err
	}

	// heartbeat
	if msg.Header.RespType == byte(protocol.HeartBeat) {
		fmt.Println("is HeartBeat")
		return "", nil, nil
	}

	b.respInterRM.RLock()
	message, ex := b.respInterMap[msg.MagicNumber]
	b.respInterRM.RUnlock()
	if !ex { // 不存在 代表消息已经失效
		fmt.Println("Not Ex", msg.MagicNumber)
		fmt.Println(b.respInterMap)
		return "", nil, nil
	}

	// 1. 解压缩
	msg.MetaData, err = b.compressor.Unzip(msg.MetaData)
	if err != nil {
		return "", nil, err
	}
	msg.Payload, err = b.compressor.Unzip(msg.Payload)
	if err != nil {
		return "", nil, err
	}
	// 2. 解密
	msg.MetaData, err = cryptology.AESDecrypt(b.options.aesKey, msg.MetaData)
	if err != nil {
		return "", nil, err
	}

	msg.Payload, err = cryptology.AESDecrypt(b.options.aesKey, msg.Payload)
	if err != nil {
		return "", nil, err
	}
	// 3. 反序列化 RespError
	mtData := make(map[string]string)
	err = b.serialization.Decode(msg.MetaData, &mtData)
	if err != nil {
		return "", nil, err
	}

	message.ctx.SetMetaData(mtData)

	value := message.ctx.Value("RespError")
	if value != "" {
		return msg.MagicNumber, message.respChan, errors.New(value)
	}

	return msg.MagicNumber, message.respChan, b.serialization.Decode(msg.Payload, message.response)
}
