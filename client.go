package hup

import (
	"io"
	"net"
)

type Client struct {
	Conn    net.Conn
	Name    string
	Addr    string
	Network string
	dp      *Pack
}

type MsgRecv struct {
	Route string
	MsgID uint32
	Body  []byte
}

func NewClient(name string, network string, addr string) (*Client, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		Conn:    conn,
		Network: network,
		Name:    name,
		Addr:    addr,
		dp:      NewPack(),
	}, nil
}

func (c *Client) Send(route string, msgID uint32, data []byte) error {
	body, err := c.dp.Pack(NewMessage(route, msgID, data))
	if err != nil {
		return err
	}
	_, err = c.Conn.Write(body)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Recv() (*MsgRecv, error) {
	headData := make([]byte, c.dp.GetHeadLen())
	if _, err := io.ReadFull(c.Conn, headData); err != nil {
		Error("read msg head err:", err.Error())
		return nil, err
	}

	msgHead, err := c.dp.Unpack(headData)
	if err != nil {
		Error("unpack err:", err)
		return nil, err
	}

	var data []byte
	if msgHead.DataLen > 0 {
		data = make([]byte, msgHead.DataLen)
		if _, err := io.ReadFull(c.Conn, data); err != nil {
			Error("read msg data error ", err)
			return nil, err
		}
	}

	var routeData []byte
	if msgHead.RouteLen > 0 {
		routeData = make([]byte, msgHead.RouteLen)
		if _, err := io.ReadFull(c.Conn, routeData); err != nil {
			Error("read msg data error ", err)
			return nil, err
		}
	}
	msgRecv := &MsgRecv{
		Route: string(routeData),
		MsgID: msgHead.ID,
		Body:  data}
	return msgRecv, nil
}

func (c *Client) Exit() {
	if c != nil {
		c.Conn.Close()
	}
	return
}
