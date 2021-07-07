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
	message, err := EncodeMessage(serverName, serverPath, byte(Request), byte(codes.Byte), byte(codes.Json), encode)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Printf("%+v \n", message)
}
