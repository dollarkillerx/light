package protocol

import (
	"encoding/binary"
	"github.com/rs/xid"
	"hash/crc32"

	"github.com/dollarkillerx/light/pkg"
)

type RequestType byte

const (
	Request RequestType = iota
	Response
)

var MaxPayloadMemory = 10 << 20 // 每个请求体最大 10M 超过进行拆分
var Crc32 = true

/**
	协议设计
	// 每个请求体最大 10M 超过进行拆分
	// crc32校验, 当前消息总数,  当前消息offset , key大小, 当前请求ID (github.com/rs/xid go客户端使用xid生成), 请求体
	crc32	:	total	:	offset	: magicNumberSize: magicNumber: serverNameSize : serverMethodSize:  respType : compressorType: serializationType : serverName : serverMethod :  payload
    4 		:	4 		: 	4 	    :     4          :     xxxx   :       4        :         4        :     1    :        1      :          1         : xxx        :      xxx     :  xxx
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

type BaseMessage struct {
	Total       uint32
	Offset      uint32
	MagicNumber string
	Data        []byte
}

// BaseDecodeMsg 基础Decode
func BaseDecodeMsg(data []byte) (*BaseMessage, error) {
	var result BaseMessage

	c32 := binary.LittleEndian.Uint32(data[:4])
	if Crc32 {
		if crc32.ChecksumIEEE(data[4:]) != c32 {
			return nil, pkg.ErrCrc32
		}
	}
	// [34 247 28 173 $ 1 0 0 0 $ 1 0 0 0 $ 12 0 0 0 $ 96 229 97 155 153 81 68 217 86 52 80 135 $ 1 0 0 0 $ 1 0 0 0 $ 0 2 0 $ 97 $ 97 $ 97]

	result.Total = binary.LittleEndian.Uint32(data[4:8])
	result.Offset = binary.LittleEndian.Uint32(data[8:12])
	magicNumberSize := binary.LittleEndian.Uint32(data[12:16])

	magicNumber := make([]byte, magicNumberSize)

	magicEndOffset := 16 + magicNumberSize
	copy(magicNumber, data[16:magicEndOffset])
	result.MagicNumber = string(magicNumber)
	result.Data = data[magicEndOffset:]

	return &result, nil
}

// DecodeMessage 完整Decode
func DecodeMessage(msg *BaseMessage) (*Message, error) {
	var result Message
	serverNameSize := binary.LittleEndian.Uint32(msg.Data[0:4])
	serverMethodSize := binary.LittleEndian.Uint32(msg.Data[4:8])
	result.RespType = msg.Data[8]
	result.CompressorType = msg.Data[9]
	result.SerializationType = msg.Data[10]
	serverName := make([]byte, serverNameSize)
	serverMethod := make([]byte, serverMethodSize)
	payload := make([]byte, len(msg.Data)-int(11+serverNameSize+serverMethodSize))

	// [1 0 0 0 $ 1 0 0 0 $ 0 2 0 $ 97 $ 97 $ 97]

	copy(serverName, msg.Data[11:11+serverNameSize])
	copy(serverMethod, msg.Data[(11+serverNameSize):(11+serverNameSize+serverMethodSize)])
	copy(payload, msg.Data[(11+serverNameSize+serverMethodSize):])

	result.ServiceName = string(serverName)
	result.ServiceMethod = string(serverMethod)
	result.Payload = payload

	return &result, nil
}

// BaseEncodeMessage 基础编码
func BaseEncodeMessage(server, method []byte, respType, compressorType, serializationType byte, payload []byte) ([]byte, error) {
	/**
	  	 serverNameSize : serverMethodSize:  respType : compressorType: serializationType : serverName : serverMethod :  payload
	           4        :         4       :     1     :        1      :          1        : xxx        :      xxx     :  xxx
	*/
	bufSize := 11 + len(server) + len(method) + len(payload)
	buf := make([]byte, bufSize)

	binary.LittleEndian.PutUint32(buf[0:4], uint32(len(server)))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(len(method)))
	buf[8] = respType
	buf[9] = compressorType
	buf[10] = serializationType
	copy(buf[11:11+len(server)], server)
	copy(buf[11+len(server):11+len(server)+len(method)], method)
	copy(buf[11+len(server)+len(method):], payload)

	return buf, nil
}

// EncodeMessage 基础编码
func EncodeMessage(server, method []byte, respType, compressorType, serializationType byte, payload []byte) ([]byte, error) {
	/**
	crc32	:	total	:	offset	: magicNumberSize: magicNumber: serverNameSize : serverMethodSize:  respType : compressorType: serializationType : serverName : serverMethod :  payload
	4 		:	4 		: 	4 	    :     4          :     xxxx   :       4        :         4        :     1    :        1      :          1        : xxx        :      xxx     :  xxx
	*/
	magicNumber := []byte(xid.New().String())
	// 如果 payload 大小 < MaxPayloadMemory 则不分包  [ 现阶段 不设 包大小限制 ]
	//if len(payload) <= MaxPayloadMemory {
	//	// total
	//	var total uint32 = 1
	//	var offset uint32 = 1
	//	bufSize := 16 + len(magicNumber)
	//	buf := make([]byte, bufSize)
	//	// 直接分装 不 分页
	//	message, err := BaseEncodeMessage(server, method, respType, compressorType, serializationType, payload)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	buf = append(buf, message...)
	//
	//	binary.LittleEndian.PutUint32(buf[4:8], total)
	//	binary.LittleEndian.PutUint32(buf[8:12], offset)
	//	binary.LittleEndian.PutUint32(buf[12:16], uint32(len(magicNumber)))
	//	copy(buf[16:16 + len(magicNumber)], magicNumber)
	//
	//	if Crc32 {
	//		u := crc32.ChecksumIEEE(buf[4:])
	//		binary.LittleEndian.PutUint32(buf[:4], u)
	//	}
	//	return buf, nil
	//}

	// total
	var total uint32 = 1
	var offset uint32 = 1
	bufSize := 16 + len(magicNumber)
	buf := make([]byte, bufSize)
	// 直接分装 不 分页
	message, err := BaseEncodeMessage(server, method, respType, compressorType, serializationType, payload)
	if err != nil {
		return nil, err
	}

	buf = append(buf, message...)

	binary.LittleEndian.PutUint32(buf[4:8], total)
	binary.LittleEndian.PutUint32(buf[8:12], offset)
	binary.LittleEndian.PutUint32(buf[12:16], uint32(len(magicNumber)))
	copy(buf[16:16+len(magicNumber)], magicNumber)

	if Crc32 {
		u := crc32.ChecksumIEEE(buf[4:])
		binary.LittleEndian.PutUint32(buf[0:4], u)
	}

	return buf, nil
}
