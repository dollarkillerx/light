package protocol

var MaxPayloadMemory = 10 << 20 // 每个请求体最大 10M 超过进行拆分

/**
	协议设计
	// 每个请求体最大 10M 超过进行拆分
	// crc32校验, 当前消息总数,  当前消息offset , key大小, 当前请求ID (github.com/rs/xid go客户端使用xid生成), 请求体
	crc32	:	total	:	offset	:	key_size : magicNumber : payload
    4 		:	4 		: 	4 	:  4 :  xxx	 :   xxx
*/

// Payload 具体 payload
type Payload struct {
	MN string `json:"mn"` // magic number
	SN string `json:"sn"` // service name
	SM string `json:"sm"` // service method
	RT byte   `json:"rt"` // resp type
	CT byte   `json:"ct"` // compressor type
	ST byte   `json:"st"` // serialization type
	P  []byte `json:"p"`  // payload 内部payload
}
