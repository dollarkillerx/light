package client

import "github.com/dollarkillerx/light"

type Client struct {
	options *Options
}

func NewClient(options ...Option) *Client {
	client := &Client{
		options: defaultOptions(),
	}

	for _, fn := range options {
		fn(client.options)
	}

	return client
}

type Connect struct {
	Client     *Client
	pool       *connectPool
	close      chan struct{}
	serverName string
}

func (c *Client) NewConnect(serverName string) (conn *Connect, err error) {
	connect := &Connect{
		Client:     c,
		serverName: serverName,
		close:      make(chan struct{}),
	}

	connect.pool, err = initPool(connect)
	return connect, err
}

func (c *Connect) Call(ctx *light.Context, serviceMethod string, request interface{}, response interface{}) error {
	return nil
}

func (c *Connect) Close() {
	close(c.close)
}
