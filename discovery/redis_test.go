package discovery

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/dollarkillerx/light/transport"
)

func TestRedis(t *testing.T) {
	discovery, err := NewRedisDiscovery("127.0.0.1:6379", 3, nil)
	if err != nil {
		log.Fatalln(err)
	}

	err = discovery.Registry("Ser", "127.0.0.1:8654", 10, transport.TCP, nil)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		servers, err := discovery.Discovery("Ser")
		if err != nil {
			log.Fatalln(err)
			return
		}
		fmt.Println(len(servers))
		time.Sleep(time.Second)
	}
}
