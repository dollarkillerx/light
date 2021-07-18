package client

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"sync"
	"time"
)

type connectPool struct {
	connect *Connect
	pool    chan LightClient

	mu sync.Mutex
}

func initPool(c *Connect) (*connectPool, error) {
	cp := &connectPool{
		connect: c,
		pool:    make(chan LightClient, c.Client.options.pool),
	}

	return cp, cp.initPool()
}

func (c *connectPool) initPool() error {
	hosts, err := c.connect.Client.options.Discovery.Discovery(c.connect.serverName)
	if err != nil {
		return err
	}

	if len(hosts) == 0 {
		return errors.New(fmt.Sprintf("%s server 404", c.connect.serverName))
	}

	c.connect.Client.options.loadBalancing.InitBalancing(hosts)

	// 初始化连接池
	for i := 0; i < c.connect.Client.options.pool; i++ {
		client, err := newBaseClient(c.connect.serverName, c.connect.Client.options)
		if err != nil {
			return errors.WithStack(err)
		}
		c.pool <- client
	}

	return nil
}

func (c *connectPool) Get(ctx context.Context) (LightClient, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("pool get timeout")
	case r := <-c.pool:
		return r, nil
	}
}

func (c *connectPool) Put(client LightClient) {
	if client.Error() == nil {
		c.pool <- client
		return
	}

	go func() {
		fmt.Println("The server starts to restore")
		for {
			time.Sleep(time.Second)
			hosts, err := c.connect.Client.options.Discovery.Discovery(c.connect.serverName)
			if err != nil {
				log.Println(err)
				continue
			}

			if len(hosts) == 0 {
				err := errors.New(fmt.Sprintf("%s server 404", c.connect.serverName))
				log.Println(err)
				continue
			}

			c.connect.Client.options.loadBalancing.InitBalancing(hosts)
			baseClient, err := newBaseClient(c.connect.serverName, c.connect.Client.options)
			if err != nil {
				log.Println(err)
				continue
			}

			c.pool <- baseClient
			fmt.Println("Service recovery success")
			break
		}
	}()
}
