package protocol

import (
	"fmt"
	"log"
	"testing"

	"github.com/dollarkillerx/light/codes"
)

type user struct {
	Name string `json:"name"`
	Psw  string `json:"psw"`
}

func TestProtocol(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	serverName := []byte("dp1")
	serverPath := []byte("com")

	usr := user{
		Name: "name...",
		Psw:  "psw...",
	}

	js, bx := codes.SerializationManager.Get(codes.Json)
	if !bx {
		log.Fatalln("what fuck?")
	}

	encode, err := js.Encode(usr)
	if err != nil {
		log.Fatalln(err)
		return
	}

	metaData := map[string]string{
		"a": "aa",
	}
	mt, err := js.Encode(metaData)
	if err != nil {
		log.Fatalln(err)
		return
	}

	_, message, err := EncodeMessage("", serverName, serverPath, mt, byte(Request), byte(codes.Byte), byte(codes.Json), encode)
	if err != nil {
		log.Fatalln(err)
		return
	}

	decodeMessage, err := DecodeMessage(message)
	if err != nil {
		log.Fatalln(err)
		return
	}

	if len(mt) != len(decodeMessage.MetaData) {
		panic("err ...")
	}

	if len(encode) != len(decodeMessage.Payload) {
		panic("err ...")
	}

	fmt.Println(decodeMessage.ServiceName)
	fmt.Println(decodeMessage.ServiceMethod)
}
