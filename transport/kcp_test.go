package transport

import (
	"fmt"
	"github.com/xtaci/kcp-go/v5"
	"time"

	"log"
	"testing"
)

func TestKCP(t *testing.T) {
	dial, err := kcp.Dial("0.0.0.0:8641")
	if err != nil {
		log.Fatalln(err)
		return
	}

	n := time.Now()
	dial.SetWriteDeadline(n.Add(time.Second))
	_, err = dial.Write([]byte(" hell world"))
	if err != nil {
		log.Fatalln(err)
		return
	}

	buf := make([]byte, 10)
	dial.SetReadDeadline(n.Add(time.Second))
	_, err = dial.Read(buf)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Println(dial.RemoteAddr())
}
