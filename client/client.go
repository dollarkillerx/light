package client

import (
	"context"
	"github.com/dollarkillerx/light"
	"github.com/dollarkillerx/light/discovery"
	"log"
	"time"
)

type Client struct {
	options *Options
}

func NewClient(discover discovery.Discovery, options ...Option) *Client {
	client := &Client{
		options: defaultOptions(),
	}

	client.options.Discovery = discover

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
	ctxT, _ := context.WithTimeout(context.TODO(), time.Second*6)
	var err error
	client, err := c.pool.Get(ctxT)
	if err != nil {
		return err
	}
	defer func() {
		go c.pool.Put(client, err)
	}()

	err = client.Call(ctx, serviceMethod, request, response)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}

func (c *Connect) Close() {
	close(c.close)
}
