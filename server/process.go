package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/dollarkillerx/light/protocol"
	"github.com/dollarkillerx/light/utils"
)

func (s *Server) process(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			utils.PrintStack()
			log.Println("Recover Err: ", err)
		}
	}() // 网络不可靠

	defer func() {
		conn.Close() // 退出 回收句柄
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
	rChannel := xChannel.Read()
loop:
	for { // 具体消息获取
		select {
		case msg, ex := <-rChannel:
			if !ex {
				break loop
			}
			now := time.Now()
			if s.options.writeTimeout != 0 {
				conn.SetWriteDeadline(now.Add(s.options.writeTimeout))
			}
			// send message
			msg = msg
		default:
			now := time.Now()
			if s.options.readTimeout != 0 {
				conn.SetReadDeadline(now.Add(s.options.readTimeout))
			}

			// TODO: heartBeat 逻辑

			proto := protocol.NewProtocol()
			msg, err := proto.IODecode(conn)
			if err != nil {
				if err == io.EOF {
					break loop
				}
				continue
			}
			marshal, err := json.Marshal(msg)
			if err != nil {
				log.Fatalln(err)
				return
			}
			fmt.Println(string(marshal))
		}
	}
}
