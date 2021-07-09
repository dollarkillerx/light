package client

import (
	"errors"
	"fmt"

	"github.com/dollarkillerx/light/transport"
)

type connectPool struct {
	connect *Connect
	pool    chan LightClient
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
		con, err := transport.Client.Gen(c.connect.Client.options.protocol, c.connect.Client.options.loadBalancing.GetService())
		if err != nil {
			return err
		}

		c.pool <- newBaseClient(con)
	}

	return nil
}
