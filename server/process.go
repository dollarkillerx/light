package server

import (
	"errors"
	"io"
	"log"
	"net"
	"reflect"
	"time"

	"github.com/dollarkillerx/light"
	"github.com/dollarkillerx/light/codes"
	"github.com/dollarkillerx/light/cryptology"
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

	defer func() {
		// 退出 回收句柄
		err := conn.Close()
		if err != nil {
			log.Println(err)
			return
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

		go s.processResponse(xChannel, msg, conn.RemoteAddr().String())
	}
}

func (s *Server) processResponse(xChannel *utils.XChannel, msg *protocol.Message, addr string) {
	var err error
	defer func() {
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
	msg.MetaData, err = cryptology.AESDecrypt(s.options.AesKey, msg.MetaData)
	if err != nil {
		return
	}

	msg.Payload, err = cryptology.AESDecrypt(s.options.AesKey, msg.Payload)
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
		auth := metaData["AUTH"]
		err := s.options.AuthFunc(ctx, auth)
		if err != nil {
			ctx.SetValue("RespError", err.Error())
			var metaDataByte []byte
			metaDataByte, _ = serialization.Encode(ctx.GetMetaData())
			metaDataByte, _ = cryptology.AESEncrypt(s.options.AesKey, metaDataByte)
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

	callErr := ser.call(ctx, method, reflect.ValueOf(req), reflect.ValueOf(resp))
	if callErr != nil {
		ctx.SetValue("RespError", callErr.Error())
	}
	// 1. 序列化
	var respBody []byte
	respBody, err = serialization.Encode(resp)

	var metaDataByte []byte
	metaDataByte, _ = serialization.Encode(ctx.GetMetaData())
	// 2. 加密
	metaDataByte, err = cryptology.AESEncrypt(s.options.AesKey, metaDataByte)
	if err != nil {
		return
	}
	respBody, err = cryptology.AESEncrypt(s.options.AesKey, respBody)
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
