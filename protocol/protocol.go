package protocol

import (
	"encoding/binary"
	"hash/crc32"

	"github.com/dollarkillerx/light/pkg"
)

var MaxPayloadMemory = 10 << 20 // 每个请求体最大 10M 超过进行拆分
var Crc32 bool

/**
	协议设计
	// 每个请求体最大 10M 超过进行拆分
	// crc32校验, 当前消息总数,  当前消息offset , key大小, 当前请求ID (github.com/rs/xid go客户端使用xid生成), 请求体
	crc32	:	total	:	offset	: magicNumberSize: serverNameSize : serverMethodSize:  respType : compressorType: serializationType :   magicNumber : serverName : serverMethod :  payload
    4 		:	4 		: 	4 	    :     4          :       4        :         4        :     1    :        1      :          1        :     xxx       : xxx        :      xxx     :  xxx
*/

type Message struct {
	Total             uint32
	Offset            uint32
	MagicNumber       string
	RespType          byte
	CompressorType    byte
	SerializationType byte
	ServiceName       string
	ServiceMethod     string
	Payload           []byte
}

func DecodeMessage(data []byte) (*Message, error) {
	var result Message

	c32 := binary.LittleEndian.Uint32(data[4:])
	if Crc32 {
		if crc32.ChecksumIEEE(data[4:]) != c32 {
			return nil, pkg.ErrCrc32
		}
	}

	result.Total = binary.LittleEndian.Uint32(data[4:8])
	result.Offset = binary.LittleEndian.Uint32(data[8:12])
	magicNumberSize := binary.LittleEndian.Uint32(data[8:12])

}
