package client

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/dollarkillerx/light"
	"github.com/dollarkillerx/light/codes"
	"github.com/dollarkillerx/light/cryptology"
	"github.com/dollarkillerx/light/pkg"
	"github.com/dollarkillerx/light/protocol"
	"github.com/dollarkillerx/light/transport"
	"github.com/dollarkillerx/light/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var compressorMin = 10 << 20 // > 10M
var compressorMax = 50 << 20 // < 50M 进行压缩

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

	err   error
	close chan struct{}
}

type respMessage struct {
	response interface{}
	ctx      *light.Context
	respChan chan error
}

func newBaseClient(serverName string, options *Options) (*BaseClient, error) {
	service, err := options.loadBalancing.GetService()
	if err != nil {
		return nil, err
	}
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

	// 握手
	encrypt, err := cryptology.RsaEncrypt([]byte(options.AUTH), options.rsaPublicKey)
	if err != nil {
		return nil, err
	}

	aesKey := []byte(strings.ReplaceAll(uuid.New().String(), "-", ""))

	// 交换秘钥
	aesKey2, err := cryptology.RsaEncrypt(aesKey, options.rsaPublicKey)
	if err != nil {
		return nil, err
	}
	handshake := protocol.EncodeHandshake(aesKey2, encrypt, []byte(""))
	_, err = con.Write(handshake)
	if err != nil {
		con.Close()
		return nil, err
	}

	hsk := &protocol.Handshake{}
	err = hsk.Handshake(con)
	if err != nil {
		con.Close()
		return nil, err
	}
	if hsk.Error != nil && len(hsk.Error) > 0 {
		con.Close()
		err := string(hsk.Error)
		return nil, errors.New(err)
	}

	bc := &BaseClient{
		serverName:    serverName,
		conn:          con,
		options:       options,
		serialization: serialization,
		compressor:    compressor,
		respInterMap:  map[string]*respMessage{},
		aesKey:        aesKey,
		close:         make(chan struct{}),
	}

	go bc.heartBeat()
	go bc.processMessageManager()

	return bc, nil
}

func (b *BaseClient) Call(ctx *light.Context, serviceMethod string, request interface{}, response interface{}) (err error) {
	now := time.Now()
	defer func() {
		if err := recover(); err != nil {
			utils.PrintStack()
			log.Println("Recover Err: ", err)
		}
		if b.options.Trace {
			log.Printf("Call %s %s time: %d/mill", b.serverName, serviceMethod, time.Since(now).Milliseconds())
		}
	}() // 网络不可靠

	// call
	var magic string
	respChan := make(chan error, 0)
	magic, err = b.call(ctx, serviceMethod, request, response, respChan)
	if err != nil {
		return errors.WithStack(err)
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

	if b.options.Trace {
		log.Printf("Call %s %s  magic: %s", b.serverName, serviceMethod, magic)
	}
	// 调度
	select {
	case err, ex := <-respChan:
		if !ex {
			return nil
		}

		return errors.WithStack(err)
	case <-time.After(timeout):
		return errors.WithStack(pkg.ErrTimeout)
	case <-b.close:
		return errors.New("net close")
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

	compressorType := b.options.compressorType
	if len(metaDataBytes) > compressorMin && len(metaDataBytes) < compressorMax {
		// 1.3 压缩
		metaDataBytes, err = b.compressor.Zip(metaDataBytes)
		if err != nil {
			return "", err
		}

		requestBytes, err = b.compressor.Zip(requestBytes)
		if err != nil {
			return "", err
		}
	} else {
		compressorType = codes.RawData
	}

	// 1.4 封装消息
	magic, message, err := protocol.EncodeMessage("", serviceNameByte, serviceMethodByte, metaDataBytes, byte(protocol.Request), byte(compressorType), byte(b.options.serializationType), requestBytes)
	if err != nil {
		return "", err
	}
	// 2. 发送消息
	if b.options.writeTimeout > 0 {
		now := time.Now()
		timeout := ctx.GetTimeout() // 如果ctx 存在设置 则采用 返之使用默认配置
		if timeout > 0 {
			b.conn.SetDeadline(now.Add(timeout))
			b.conn.SetWriteDeadline(now.Add(timeout))
		} else {
			b.conn.SetDeadline(now.Add(b.options.writeTimeout))
			b.conn.SetWriteDeadline(now.Add(b.options.writeTimeout))
		}
	}
	// 写MAP
	b.respInterRM.Lock()
	b.respInterMap[magic] = &respMessage{
		response: response,
		ctx:      ctx,
		respChan: respChan,
	}
	b.respInterRM.Unlock()

	// 有点暴力呀 直接上锁
	b.writeMu.Lock()
	_, err = b.conn.Write(message)
	b.writeMu.Unlock()
	if err != nil {
		if b.options.Trace {
			log.Println(err)
		}
		b.err = err
		return "", errors.WithStack(err)
	}

	return magic, nil
}

func (b *BaseClient) heartBeat() {
	defer func() {
		fmt.Println("heartBeat Close")
	}()

loop:
	for {
		select {
		case <-b.close:
			break loop
		case <-time.After(b.options.heartBeat):
			_, i, err := protocol.EncodeMessage("x", []byte(""), []byte(""), []byte(""), byte(protocol.HeartBeat), byte(b.options.compressorType), byte(b.options.serializationType), []byte(""))
			if err != nil {
				log.Println(err)
				break
			}
			now := time.Now()
			b.conn.SetDeadline(now.Add(b.options.writeTimeout))
			b.conn.SetWriteDeadline(now.Add(b.options.writeTimeout))
			b.writeMu.Lock()
			_, err = b.conn.Write(i)
			b.writeMu.Unlock()
			if err != nil {
				b.err = err
				break loop
			}
		}
	}
}

func (b *BaseClient) processMessageManager() {
	defer func() {
		fmt.Println("processMessageManager Close")
	}()

	for {
		magic, respChan, err := b.processMessage()
		if err == nil && magic == "" {
			continue
		}

		if err != nil && magic == "" {
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
		b.err = err
		close(b.close)
		return "", nil, err
	}

	// heartbeat
	if msg.Header.RespType == byte(protocol.HeartBeat) {
		if b.options.Trace {
			log.Println("is HeartBeat")
		}
		return "", nil, nil
	}

	b.respInterRM.RLock()
	message, ex := b.respInterMap[msg.MagicNumber]
	b.respInterRM.RUnlock()
	if !ex { // 不存在 代表消息已经失效
		if b.options.Trace {
			log.Println("Not Ex", msg.MagicNumber)
		}
		return "", nil, nil
	}

	comp, ex := codes.CompressorManager.Get(codes.CompressorType(msg.Header.CompressorType))
	if !ex {
		return "", nil, nil
	}

	// 1. 解压缩
	msg.MetaData, err = comp.Unzip(msg.MetaData)
	if err != nil {
		return "", nil, err
	}
	msg.Payload, err = comp.Unzip(msg.Payload)
	if err != nil {
		return "", nil, err
	}
	// 2. 解密
	msg.MetaData, err = cryptology.AESDecrypt(b.aesKey, msg.MetaData)
	if err != nil {
		if len(msg.MetaData) != 0 {
			return "", nil, err
		}
		msg.Payload = []byte("")
	}

	msg.Payload, err = cryptology.AESDecrypt(b.aesKey, msg.Payload)
	if err != nil {
		if len(msg.Payload) != 0 {
			return "", nil, err
		}
		msg.Payload = []byte("")
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

func (b *BaseClient) Error() error {
	return b.err
}
