package tests

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"testing"
)

func TestClient(t *testing.T) {
	dial, err := net.Dial("tcp", "127.0.0.1:8471")
	if err != nil {
		log.Fatalln(err)
		return
	}

	go func() {
		for {
			buf := make([]byte, 4)
			_, err := dial.Read(buf)
			if err != nil {
				log.Println(err)
				return
			}
			u := binary.LittleEndian.Uint32(buf)
			buf = make([]byte, u)
			_, err = dial.Read(buf)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println("r: ", string(buf))
		}
	}()

	for {
		//time.Sleep(time.Nanosecond * 1)
		cp := "hello world"
		buf := make([]byte, 4+len(cp))
		binary.LittleEndian.PutUint32(buf[:4], uint32(len(cp)))
		copy(buf[4:], cp)
		_, err := dial.Write(buf)
		if err != nil {
			break
		}
	}
}

func TestServer(t *testing.T) {
	listen, err := net.Listen("tcp", "127.0.0.1:8471")
	if err != nil {
		log.Fatalln(err)
		return
	}

	for {
		accept, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go func() {
			for {
				buf := make([]byte, 4)
				_, err := accept.Read(buf)
				if err != nil {
					log.Println(err)
					return
				}
				u := binary.LittleEndian.Uint32(buf)
				buf = make([]byte, u)
				_, err = accept.Read(buf)
				if err != nil {
					log.Println(err)
					return
				}
				fmt.Println("r: ", string(buf))
			}
		}()

		go func() {
			for {
				//time.Sleep(time.Nanosecond * 1)
				cp := "resp: hello world"
				buf := make([]byte, 4+len(cp))
				binary.LittleEndian.PutUint32(buf[:4], uint32(len(cp)))
				copy(buf[4:], cp)
				accept.Write(buf)
			}
		}()
	}
}
