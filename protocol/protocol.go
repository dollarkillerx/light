package protocol

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"log"

	"github.com/rs/xid"
)

type RequestType byte

const (
	Request RequestType = iota
	Response
	HeartBeat
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
	Header        *Header
	MagicNumber   string
	ServiceName   string
	ServiceMethod string
	MetaData      []byte
	Payload       []byte
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

type Protocol struct{}

func NewProtocol() *Protocol {
	return &Protocol{}
}

func (m *Protocol) IODecode(r io.Reader) (*Message, error) {
	headerByte := make([]byte, HeadSize)
	// 读取标志位
	_, err := io.ReadFull(r, headerByte[:1])
	if err != nil {
		return nil, err
	}

	if headerByte[0] != LightSt {
		log.Println(headerByte)
		return nil, errors.New("NO Light")
	}

	// 读取剩下的
	_, err = io.ReadFull(r, headerByte[1:])
	if err != nil {
		return nil, err
	}

	// 解析 header
	header, err := DecodeHeader(headerByte)
	if err != nil {
		return nil, err
	}

	bodyLen := header.MagicNumberSize + header.ServerNameSize + header.ServerMethodSize + header.MetaDataSize + header.PayloadSize
	bodyData := make([]byte, bodyLen)
	_, err = io.ReadFull(r, bodyData)
	if err != nil {
		return nil, err
	}

	msg, err := DecodeMessageV2(bodyData, header, 0)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func DecodeHeader(data []byte) (*Header, error) {
	var header Header
	header.St = data[0]
	header.Version = data[1]
	header.Crc32 = binary.BigEndian.Uint32(data[2:6])
	if Crc32 {
		u := crc32.ChecksumIEEE(data[6:])
		if header.Crc32 != u {
			return nil, errors.New("CRC Calibration")
		}
	}
	header.MagicNumberSize = binary.BigEndian.Uint32(data[6:10])
	header.ServerNameSize = binary.BigEndian.Uint32(data[10:14])
	header.ServerMethodSize = binary.BigEndian.Uint32(data[14:18])
	header.MetaDataSize = binary.BigEndian.Uint32(data[18:22])
	header.PayloadSize = binary.BigEndian.Uint32(data[22:26])
	header.RespType = data[26]
	header.CompressorType = data[27]
	header.SerializationType = data[28]

	return &header, nil
}

// DecodeMessage 完整Decode
func DecodeMessage(data []byte) (*Message, error) {
	header, err := DecodeHeader(data)
	if err != nil {
		return nil, err
	}
	return DecodeMessageV2(data, header, HeadSize)
}

func DecodeMessageV2(data []byte, header *Header, headSize uint32) (*Message, error) {
	var result Message
	result.Header = header
	var st uint32 = headSize
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
func EncodeMessage(magicStr string, server, method, metaData []byte, respType, compressorType, serializationType byte, payload []byte) (magic string, data []byte, err error) {
	var magicNumber = []byte(magicStr)
	if magicStr == "" {
		magicNumber = []byte(xid.New().String())
	}

	bufSize := HeadSize + len(server) + len(method) + len(metaData) + len(payload) + len(magicNumber)
	buf := make([]byte, bufSize)

	buf[0] = LightSt
	buf[1] = byte(V1)
	binary.BigEndian.PutUint32(buf[6:10], uint32(len(magicNumber)))
	binary.BigEndian.PutUint32(buf[10:14], uint32(len(server)))
	binary.BigEndian.PutUint32(buf[14:18], uint32(len(method)))
	binary.BigEndian.PutUint32(buf[18:22], uint32(len(metaData)))
	binary.BigEndian.PutUint32(buf[22:26], uint32(len(payload)))
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
	copy(buf[st:endI], payload)

	if Crc32 {
		u := crc32.ChecksumIEEE(buf[6:])
		binary.BigEndian.PutUint32(buf[2:6], u)
	}

	return string(magicNumber), buf, nil
}

const LightHandshakeSt = 0x09 // Handshake 起始符

/**
	Handshake 协议设计
	起始符 :  keySize  :  tokenSize :  errorSize :  key  : token :  err :
    0x09  :     4     :     4      :    4       :   xxx  :  xxx  :  xxx :
*/

type Handshake struct {
	Key   []byte
	Token []byte
	Error []byte
}

// EncodeHandshake 编码握手函数
func EncodeHandshake(key, token, err []byte) []byte {
	buf := make([]byte, 13+len(key)+len(token)+len(err))

	buf[0] = LightHandshakeSt
	binary.BigEndian.PutUint32(buf[1:5], uint32(len(key)))
	binary.BigEndian.PutUint32(buf[5:9], uint32(len(token)))
	binary.BigEndian.PutUint32(buf[9:13], uint32(len(err)))

	st := 13
	endI := st + len(key)
	copy(buf[st:endI], key)

	st = endI
	endI = st + len(token)
	copy(buf[st:endI], token)

	st = endI
	endI = st + len(err)
	copy(buf[st:endI], err)

	return buf
}

// Handshake 解码握手函数
func (h *Handshake) Handshake(r io.Reader) error {
	headerByte := make([]byte, 13)

	// 读取标志位
	_, err := io.ReadFull(r, headerByte[:1])
	if err != nil {
		return err
	}

	if headerByte[0] != LightHandshakeSt {
		return errors.New("NO Light Handshake")
	}

	// 读取剩下的
	_, err = io.ReadFull(r, headerByte[1:])
	if err != nil {
		return err
	}

	// 解析头
	keySize := binary.BigEndian.Uint32(headerByte[1:5])
	tokenSize := binary.BigEndian.Uint32(headerByte[5:9])
	errorSize := binary.BigEndian.Uint32(headerByte[9:13])

	bodyData := make([]byte, keySize+tokenSize+errorSize)
	_, err = io.ReadFull(r, bodyData)
	if err != nil {
		return err
	}

	var st uint32 = 0
	endIdx := keySize
	key := make([]byte, keySize)
	copy(key, bodyData[st:endIdx])
	h.Key = key

	st = endIdx
	endIdx = st + tokenSize
	token := make([]byte, tokenSize)
	copy(token, bodyData[st:endIdx])
	h.Token = token

	st = endIdx
	endIdx = st + errorSize
	errStr := make([]byte, errorSize)
	copy(errStr, bodyData[st:endIdx])
	h.Error = errStr

	return nil
}
