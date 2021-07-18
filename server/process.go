package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"time"

	"github.com/dollarkillerx/light"
	"github.com/dollarkillerx/light/codes"
	"github.com/dollarkillerx/light/cryptology"
	"github.com/dollarkillerx/light/pkg"
	"github.com/dollarkillerx/light/protocol"
	"github.com/dollarkillerx/light/utils"
)

func (s *Server) process(conn net.Conn) {

	defer func() {
		// 网络不可靠
		if err := recover(); err != nil {
			utils.PrintStack()
			log.Println("Recover Err: ", err)
		}
	}()

	s.options.Discovery.Add(1)
	defer func() {
		s.options.Discovery.Less(1)
		// 退出 回收句柄
		err := conn.Close()
		if err != nil {
			log.Println(err)
			return
		}

		if s.options.Trace {
			log.Println("close connect: ", conn.RemoteAddr())
		}
	}()

	// 初始化
	if tc, ok := conn.(*net.TCPConn); ok {
		err := tc.SetKeepAlive(true)
		if err != nil {
			log.Println(err)
			return
		}
		err = tc.SetKeepAlivePeriod(s.options.options["TCPKeepAlivePeriod"].(time.Duration))
		if err != nil {
			log.Println(err)
			return
		}
		err = tc.SetLinger(10)
		if err != nil {
			log.Println(err)
			return
		}
	}

	xChannel := utils.NewXChannel(s.options.processChanSize)

	// 握手
	handshake := protocol.Handshake{}
	err := handshake.Handshake(conn)
	if err != nil {
		return
	}

	aesKey, err := cryptology.RsaDecrypt(handshake.Key, s.options.RSAPrivateKey)
	if err != nil {
		encodeHandshake := protocol.EncodeHandshake([]byte(""), []byte(""), []byte(err.Error()))
		conn.Write(encodeHandshake)
		return
	}

	if len(aesKey) != 32 && len(aesKey) != 16 {
		encodeHandshake := protocol.EncodeHandshake([]byte(""), []byte(""), []byte("aes key != 32 && key != 16"))
		conn.Write(encodeHandshake)
		return
	}

	token, err := cryptology.RsaDecrypt(handshake.Token, s.options.RSAPrivateKey)
	if err != nil {
		encodeHandshake := protocol.EncodeHandshake([]byte(""), []byte(""), []byte(err.Error()))
		conn.Write(encodeHandshake)
		return
	}

	if s.options.AuthFunc != nil {
		err := s.options.AuthFunc(light.DefaultCtx(), string(token))
		if err != nil {
			encodeHandshake := protocol.EncodeHandshake([]byte(""), []byte(""), []byte(err.Error()))
			conn.Write(encodeHandshake)
			return
		}
	}

	// limit 限流
	if s.options.Discovery.Limit() {
		// 熔断
		encodeHandshake := protocol.EncodeHandshake([]byte(""), []byte(""), []byte(pkg.ErrCircuitBreaker.Error()))
		conn.Write(encodeHandshake)
		log.Println(s.options.Discovery.Limit())
		return
	}

	encodeHandshake := protocol.EncodeHandshake([]byte(""), []byte(""), []byte(""))
	_, err = conn.Write(encodeHandshake)
	if err != nil {
		return
	}
	// send
	go func() {
	loop:
		for {
			select {
			case msg, ex := <-xChannel.Ch:
				if !ex {
					if s.options.Trace {
						log.Printf("ip: %s  close send server", conn.RemoteAddr())
					}
					break loop
				}
				now := time.Now()
				if s.options.writeTimeout > 0 {
					conn.SetWriteDeadline(now.Add(s.options.writeTimeout))
				}
				// send message
				_, err := conn.Write(msg)
				if err != nil {
					if s.options.Trace {
						log.Printf("ip: %s err: %s", conn.RemoteAddr(), err)
					}
					break loop
				}
			}
		}
	}()

	defer func() {
		xChannel.Close()
	}()
loop:
	for { // 具体消息获取
		now := time.Now()
		if s.options.readTimeout > 0 {
			conn.SetReadDeadline(now.Add(s.options.readTimeout))
		}

		proto := protocol.NewProtocol()
		msg, err := proto.IODecode(conn)
		if err != nil {
			if err == io.EOF {
				if s.options.Trace {
					log.Printf("ip: %s close", conn.RemoteAddr())
				}
				break loop
			}

			// 遇到错误关闭链接
			if s.options.Trace {
				log.Printf("ip: %s err: %s", conn.RemoteAddr(), err)
			}
			break loop
		}

		go s.processResponse(xChannel, msg, conn.RemoteAddr().String(), aesKey)
	}
}

func (s *Server) processResponse(xChannel *utils.XChannel, msg *protocol.Message, addr string, aesKey []byte) {
	var err error
	s.options.Discovery.Add(1)
	defer func() {
		s.options.Discovery.Less(1)
		if err != nil {
			if s.options.Trace {
				log.Println("ProcessResponse Error: ", err, "  ID: ", addr)
			}
			xChannel.Close()
		}
	}()

	// heartBeat 判断
	if msg.Header.RespType == byte(protocol.HeartBeat) {
		// 心跳返回
		if s.options.Trace {
			log.Println("HeartBeat: ", addr)
		}

		// 4. 打包
		_, message, err := protocol.EncodeMessage(msg.MagicNumber, []byte(msg.ServiceName), []byte(msg.ServiceMethod), []byte(""), byte(protocol.HeartBeat), msg.Header.CompressorType, msg.Header.SerializationType, []byte(""))
		if err != nil {
			return
		}
		// 5. 回写
		err = xChannel.Send(message)
		if err != nil {
			return
		}

		return
	}

	// 限流
	if s.options.Discovery.Limit() {
		serialization, _ := codes.SerializationManager.Get(codes.MsgPack)
		metaData := make(map[string]string)
		metaData["RespError"] = pkg.ErrCircuitBreaker.Error()
		meta, err := serialization.Encode(metaData)
		if err != nil {
			return
		}
		decrypt, err := cryptology.AESDecrypt(aesKey, meta)
		if err != nil {
			return
		}
		_, message, err := protocol.EncodeMessage(msg.MagicNumber, []byte(msg.ServiceName), []byte(msg.ServiceMethod), decrypt, byte(protocol.Response), byte(codes.RawData), byte(codes.MsgPack), []byte(""))
		if err != nil {
			return
		}
		// 5. 回写
		err = xChannel.Send(message)
		if err != nil {
			return
		}

		log.Println(s.options.Discovery.Limit())
		log.Println("限流/////////////")

		return
	}

	// 1. 解压缩
	compressor, ex := codes.CompressorManager.Get(codes.CompressorType(msg.Header.CompressorType))
	if !ex {
		err = errors.New("compressor 404")
		return
	}
	msg.MetaData, err = compressor.Unzip(msg.MetaData)
	if err != nil {
		return
	}

	msg.Payload, err = compressor.Unzip(msg.Payload)
	if err != nil {
		return
	}
	// 2. 解密
	msg.MetaData, err = cryptology.AESDecrypt(aesKey, msg.MetaData)
	if err != nil {
		return
	}

	msg.Payload, err = cryptology.AESDecrypt(aesKey, msg.Payload)
	if err != nil {
		return
	}

	// 3. 反序列化
	serialization, ex := codes.SerializationManager.Get(codes.SerializationType(msg.Header.SerializationType))
	if !ex {
		err = errors.New("serialization 404")
		return
	}

	metaData := make(map[string]string)
	err = serialization.Decode(msg.MetaData, &metaData)
	if err != nil {
		return
	}

	ctx := light.DefaultCtx()
	ctx.SetMetaData(metaData)

	// 1.3 auth
	if s.options.AuthFunc != nil {
		auth := metaData["Light_AUTH"]
		err := s.options.AuthFunc(ctx, auth)
		if err != nil {
			ctx.SetValue("RespError", err.Error())
			var metaDataByte []byte
			metaDataByte, _ = serialization.Encode(ctx.GetMetaData())
			metaDataByte, _ = cryptology.AESEncrypt(aesKey, metaDataByte)
			metaDataByte, _ = compressor.Zip(metaDataByte)
			// 4. 打包
			_, message, err := protocol.EncodeMessage(msg.MagicNumber, []byte(msg.ServiceName), []byte(msg.ServiceMethod), metaDataByte, byte(protocol.Response), msg.Header.CompressorType, msg.Header.SerializationType, []byte(""))
			if err != nil {
				return
			}
			// 5. 回写
			err = xChannel.Send(message)
			if err != nil {
				return
			}
			return
		}
	}

	ser, ex := s.serviceMap[msg.ServiceName]
	if !ex {
		err = errors.New("service does not exist")
		return
	}

	method, ex := ser.methodType[msg.ServiceMethod]
	if !ex {
		err = errors.New("method does not exist")
		return
	}

	req := utils.RefNew(method.RequestType)
	resp := utils.RefNew(method.ResponseType)

	err = serialization.Decode(msg.Payload, req)
	if err != nil {
		return
	}

	path := fmt.Sprintf("%s.%s", msg.ServiceName, msg.ServiceMethod)
	ctx.SetPath(path)

	// 前置middleware
	if len(s.beforeMiddleware) != 0 {
		for idx := range s.beforeMiddleware {
			err := s.beforeMiddleware[idx](ctx, req, resp)
			if err != nil {
				return
			}
		}
	}
	funcs, ex := s.beforeMiddlewarePath[path]
	if ex {
		if len(funcs) != 0 {
			for idx := range funcs {
				err := funcs[idx](ctx, req, resp)
				if err != nil {
					return
				}
			}
		}
	}

	// 核心调用
	callErr := ser.call(ctx, method, reflect.ValueOf(req), reflect.ValueOf(resp))
	if callErr != nil {
		ctx.SetValue("RespError", callErr.Error())
	}

	// 后置middleware
	if len(s.afterMiddleware) != 0 {
		for idx := range s.afterMiddleware {
			err := s.afterMiddleware[idx](ctx, req, resp)
			if err != nil {
				return
			}
		}
	}
	funcs, ex = s.afterMiddlewarePath[path]
	if ex {
		if len(funcs) != 0 {
			for idx := range funcs {
				err := funcs[idx](ctx, req, resp)
				if err != nil {
					return
				}
			}
		}
	}
	// response

	// 1. 序列化
	var respBody []byte
	respBody, err = serialization.Encode(resp)

	var metaDataByte []byte
	metaDataByte, _ = serialization.Encode(ctx.GetMetaData())
	// 2. 加密
	metaDataByte, err = cryptology.AESEncrypt(aesKey, metaDataByte)
	if err != nil {
		return
	}
	respBody, err = cryptology.AESEncrypt(aesKey, respBody)
	if err != nil {
		return
	}
	// 3. 压缩
	metaDataByte, err = compressor.Zip(metaDataByte)
	if err != nil {
		return
	}
	respBody, err = compressor.Zip(respBody)
	if err != nil {
		return
	}
	// 4. 打包
	_, message, err := protocol.EncodeMessage(msg.MagicNumber, []byte(msg.ServiceName), []byte(msg.ServiceMethod), metaDataByte, byte(protocol.Response), msg.Header.CompressorType, msg.Header.SerializationType, respBody)
	if err != nil {
		return
	}
	// 5. 回写
	err = xChannel.Send(message)
	if err != nil {
		return
	}
}
