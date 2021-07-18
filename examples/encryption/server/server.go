package main

import (
	"fmt"
	"log"

	"github.com/dollarkillerx/light"
	"github.com/dollarkillerx/light/server"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	ser := server.NewServer()
	err := ser.RegisterName(&HelloWorld{}, "helloworld")
	if err != nil {
		log.Fatalln(err)
	}

	publicKey := []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDMacmioOUq4HTMKVutxsrWimQO
vFOZIU93NKJYusRFV8lN8NB0dSg5AcbZgyYegY07mWXBBg8zlI+4PphUj40kn0F3
aOnvZ6WrsyYlPi1ZnXBaTFXxC6YN2LH9Lb9KaWrtZH4AM+6PoXIjmtWhpZr1JuuD
7J540DLMsuoEGzLoRQIDAQAB
-----END PUBLIC KEY-----`)
	privateKey := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDMacmioOUq4HTMKVutxsrWimQOvFOZIU93NKJYusRFV8lN8NB0
dSg5AcbZgyYegY07mWXBBg8zlI+4PphUj40kn0F3aOnvZ6WrsyYlPi1ZnXBaTFXx
C6YN2LH9Lb9KaWrtZH4AM+6PoXIjmtWhpZr1JuuD7J540DLMsuoEGzLoRQIDAQAB
AoGBAIjyut8U0lQWismZT82t6IkxsHVa4NsvwPCJN9cgUYxvkvN/yfir4SXINfPV
9LztaLSsQcq/B4I0HtF+Rkoo3pKe+9FlDk+cnylNcGPznSdYH4oI6Z5lvaprnthA
yc4KKCBTn+UT06d80im4HjzzIstUjRSgtPTuS5yqZCLL+8ZlAkEA4QAzmaVge59w
x3kbTL9kcbVm7aLVb8UHjiu4yej2xjGOL4H4xSSG0bbwrTk8l4QyV932/TV91933
I72HRyJGKwJBAOiTdg1FZ9ioHtNjSVgOtfduaySzJZgunbItex+CWddq9p7t50j0
5IOHUOI35lamqrKZAOY8O7vcLBEDs77GQ08CQQDOEfgwdWWbc5jAKKwXK4ecGQ9O
//7JYkQcMvEIg8RYGxTlb/1e2ahctFdT34MeJiZRkWpf2DkMly99XV1jigGHAkAN
pvJMDyHsZtoAYJiikaJ+1r11VwrC5yGcnuzWSamKap31cFOeRbnQOrY1wUBFH91v
RGn4GdsLyP3RNd1sOmkjAkEAk/LxZ0oq4ViOtmMT2BrZTK5b9E0p/XchnJv+u7FR
Ni4dv5khSlSOf1z9ySCrYaYvHP+LPoUx+KWa7xyhehtpFw==
-----END RSA PRIVATE KEY-----`)

	if err := ser.Run(server.UseTCP("0.0.0.0:8074"), server.SetRSAKey(publicKey, privateKey)); err != nil {
		log.Fatalln(err)
	}
}

type HelloWorld struct{}

type HelloWorldRequest struct {
	Name string
}

type HelloWorldResponse struct {
	RPName string
}

func (s *HelloWorld) HelloWorld(ctx *light.Context, req *HelloWorldRequest, resp *HelloWorldResponse) error {
	resp.RPName = fmt.Sprintf("hello world by: %s", req.Name)
	return nil
}
