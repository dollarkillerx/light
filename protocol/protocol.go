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
	crc32	:	total	:	offset	: magicNumberSize: magicNumber: serverNameSize : serverMethodSize:  metaDataSize: respType : compressorType: serializationType : metaDataValue : serverName : serverMethod :  payload
    4 		:	4 		: 	4 	    :     4          :     xxxx   :       4        :         4        :     4       :    1    :        1      :          1         :      xxx      :      xxx   :      xxx     :   xxxx
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
	MetaData          []byte
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
	metaDataSize := binary.LittleEndian.Uint32(msg.Data[8:12])
	result.RespType = msg.Data[12]
	result.CompressorType = msg.Data[13]
	result.SerializationType = msg.Data[14]
	serverName := make([]byte, serverNameSize)
	serverMethod := make([]byte, serverMethodSize)
	metaData := make([]byte, metaDataSize)
	payload := make([]byte, len(msg.Data)-int(15+serverNameSize+serverMethodSize+metaDataSize))

	// [1 0 0 0 $ 1 0 0 0 $ 0 2 0 $ 97 $ 97 $ 97]
	copy(metaData, msg.Data[15:15+metaDataSize])
	copy(serverName, msg.Data[15+metaDataSize:15+metaDataSize+serverNameSize])
	copy(serverMethod, msg.Data[(15+metaDataSize+serverNameSize):(15+metaDataSize+serverNameSize+serverMethodSize)])
	copy(payload, msg.Data[(15+metaDataSize+serverNameSize+serverMethodSize):])

	result.ServiceName = string(serverName)
	result.ServiceMethod = string(serverMethod)
	result.Payload = payload
	result.MetaData = metaData

	return &result, nil
}

// BaseEncodeMessage 基础编码
func BaseEncodeMessage(server, method, metaData []byte, respType, compressorType, serializationType byte, payload []byte) ([]byte, error) {
	/**
	  	 serverNameSize : serverMethodSize:  metaDataSize : respType : compressorType: serializationType:  metaDataVal : serverName : serverMethod :  payload
	           4        :         4       :     4         :    1     :        1      :          1       :    xxx       : xxx        :      xxx     :  xxx
	*/
	bufSize := 15 + len(server) + len(method) + len(metaData) + len(payload)
	buf := make([]byte, bufSize)

	binary.LittleEndian.PutUint32(buf[0:4], uint32(len(server)))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(len(method)))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(len(method)))
	binary.LittleEndian.PutUint32(buf[8:12], uint32(len(metaData)))
	buf[12] = respType
	buf[13] = compressorType
	buf[14] = serializationType
	copy(buf[15:15+len(metaData)], metaData)
	copy(buf[15+len(metaData):15+len(metaData)+len(server)], server)
	copy(buf[15+len(metaData)+len(server):15+len(metaData)+len(server)+len(method)], method)
	copy(buf[15+len(metaData)+len(server)+len(method):], payload)

	return buf, nil
}

// EncodeMessage 基础编码
func EncodeMessage(server, method, metaData []byte, respType, compressorType, serializationType byte, payload []byte) ([]byte, error) {
	/**
	crc32	:	total	:	offset	: magicNumberSize: magicNumber: serverNameSize : serverMethodSize:  metaDataSize: respType : compressorType: serializationType : metaDataValue : serverName : serverMethod :  payload
	4 		:	4 		: 	4 	    :     4          :     xxxx   :       4        :         4        :     4       :    1    :        1      :          1         :      xxx      :      xxx   :      xxx     :   xxxx
	*/
	magicNumber := []byte(xid.New().String())
	// 如果 payload 大小 < MaxPayloadMemory 则不分包  [ 现阶段 不设 包大小限制 ]
	//if len(payload) <= MaxPayloadMemory {

	var total uint32 = 1
	var offset uint32 = 1
	bufSize := 16 + len(magicNumber)
	buf := make([]byte, bufSize)
	// 直接分装 不 分页
	message, err := BaseEncodeMessage(server, method, metaData, respType, compressorType, serializationType, payload)
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
