package client

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
	serverName string
}

func (c *Client) NewConnect(serverName string) *Connect {
	return &Connect{
		Client:     c,
		serverName: serverName,
	}
}

func (c *Client) Call() error {
	return nil
}
