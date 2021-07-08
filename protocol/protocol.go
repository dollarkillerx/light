package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"

	"github.com/rs/xid"
)

type RequestType byte

const (
	Request RequestType = iota
	Response
)

//var MaxPayloadMemory = 10 << 20 // 每个请求体最大 10M 超过进行拆分
var Crc32 = false

type LightVersion byte

const (
	V1 LightVersion = iota
)

const LightSt = 0x05 // 起始符
const HeadSize = 29

/**
	协议设计
	起始符 :  版本号 :  crc32校验 :   magicNumberSize:    serverNameSize :   serverMethodSize :  metaDataSize : payloadSize:  respType :   compressorType :    serializationType :    magicNumber :  serverName :   serverMethod :  metaData :  payload
    0x05  :  0x01  :     4     :        4         :         4         :         4          :       4       :      4     :      1    :          1       :           1          :        xxx     :       xxx   :        xxx     :    xxx    :    xxx
*/

type Message struct {
	Header            *Header
	MagicNumber       string
	RespType          byte
	CompressorType    byte
	SerializationType byte
	ServiceName       string
	ServiceMethod     string
	MetaData          []byte
	Payload           []byte
}

type Header struct {
	St                byte
	Version           byte
	Crc32             uint32
	MagicNumberSize   uint32
	ServerNameSize    uint32
	ServerMethodSize  uint32
	MetaDataSize      uint32
	PayloadSize       uint32
	RespType          byte
	CompressorType    byte
	SerializationType byte
}

func DecodeHeader(data []byte) (*Header, error) {
	var header Header
	header.St = data[0]
	header.Version = data[1]
	header.Crc32 = binary.LittleEndian.Uint32(data[2:6])
	if Crc32 {
		u := crc32.ChecksumIEEE(data[6:])
		if header.Crc32 != u {
			return nil, errors.New("CRC Calibration")
		}
	}
	header.MagicNumberSize = binary.LittleEndian.Uint32(data[6:10])
	header.ServerNameSize = binary.LittleEndian.Uint32(data[10:14])
	header.ServerMethodSize = binary.LittleEndian.Uint32(data[14:18])
	header.MetaDataSize = binary.LittleEndian.Uint32(data[18:22])
	header.PayloadSize = binary.LittleEndian.Uint32(data[22:26])
	header.RespType = data[26]
	header.CompressorType = data[27]
	header.SerializationType = data[28]

	return &header, nil
}

// DecodeMessage 完整Decode
func DecodeMessage(data []byte) (*Message, error) {
	var result Message
	header, err := DecodeHeader(data)
	if err != nil {
		return nil, err
	}
	result.Header = header

	var st uint32 = HeadSize
	endI := st + header.MagicNumberSize
	les := endI - st
	magicNumber := make([]byte, les)
	copy(magicNumber, data[st:endI])
	result.MagicNumber = string(magicNumber)

	st = endI
	endI = st + header.ServerNameSize
	les = endI - st
	serverName := make([]byte, les)
	copy(serverName, data[st:endI])
	result.ServiceName = string(serverName)

	st = endI
	endI = st + header.ServerMethodSize
	les = endI - st
	serverMethodSize := make([]byte, les)
	copy(serverMethodSize, data[st:endI])
	result.ServiceMethod = string(serverMethodSize)

	st = endI
	endI = st + header.MetaDataSize
	les = endI - st
	metaDataSize := make([]byte, les)
	copy(metaDataSize, data[st:endI])
	result.MetaData = metaDataSize

	st = endI
	endI = st + header.PayloadSize
	les = endI - st
	payloadSize := make([]byte, les)
	copy(payloadSize, data[st:endI])
	result.Payload = payloadSize

	return &result, nil
}

// EncodeMessage 基础编码
func EncodeMessage(server, method, metaData []byte, respType, compressorType, serializationType byte, payload []byte) ([]byte, error) {
	magicNumber := []byte(xid.New().String())

	bufSize := HeadSize + len(server) + len(method) + len(metaData) + len(payload) + len(magicNumber)
	buf := make([]byte, bufSize)

	buf[0] = LightSt
	buf[1] = byte(V1)
	binary.LittleEndian.PutUint32(buf[6:10], uint32(len(magicNumber)))
	binary.LittleEndian.PutUint32(buf[10:14], uint32(len(server)))
	binary.LittleEndian.PutUint32(buf[14:18], uint32(len(method)))
	binary.LittleEndian.PutUint32(buf[18:22], uint32(len(metaData)))
	binary.LittleEndian.PutUint32(buf[22:26], uint32(len(payload)))
	buf[26] = respType
	buf[27] = compressorType
	buf[28] = serializationType

	st := HeadSize
	endI := st + len(magicNumber)
	copy(buf[st:endI], magicNumber)

	st = endI
	endI = st + len(server)
	copy(buf[st:endI], server)

	st = endI
	endI = st + len(method)
	copy(buf[st:endI], method)

	st = endI
	endI = st + len(metaData)
	copy(buf[st:endI], metaData)

	st = endI
	endI = st + len(payload)
	fmt.Println(bufSize)
	copy(buf[st:endI], payload)

	if Crc32 {
		u := crc32.ChecksumIEEE(buf[6:])
		binary.LittleEndian.PutUint32(buf[2:6], u)
	}

	return buf, nil
}
