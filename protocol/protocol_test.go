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

	//serverName = []byte("a")
	//serverPath = []byte("a")

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
	//encode = []byte("a")
	fmt.Println(encode)
	fmt.Printf("req: %+v  byt: %+v  json: %+v \n", byte(Request), byte(codes.Byte), byte(codes.Json))

	metaData := map[string]string{
		"a": "aa",
	}
	bytes, err := js.Encode(metaData)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Println("meta data: ", bytes)
	message, err := EncodeMessage(serverName, serverPath, bytes, byte(Request), byte(codes.Byte), byte(codes.Json), encode)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Printf("%+v \n", message)

	msg, err := BaseDecodeMsg(message)
	if err != nil {
		log.Fatalln(err)
		return
	}

	decodeMessage, err := DecodeMessage(msg)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Println(decodeMessage)
}
