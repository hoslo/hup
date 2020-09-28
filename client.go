package hup

import (
	"io"
	"net"
	"time"
)

type Client struct {
	Conn    net.Conn
	Name    string
	Addr    string
	Network string
	dp      *DataPack
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
		dp:      NewDataPack(),
	}, nil
}

func (c *Client) Send(route string, msgID uint32, data []byte) error {
	body, err := c.dp.Pack(NewMsgPackage(route, msgID, data))
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
	if msgHead.GetDataLen() > 0 {
		data = make([]byte, msgHead.GetDataLen())
		if _, err := io.ReadFull(c.Conn, data); err != nil {
			Error("read msg data error ", err)
			return nil, err
		}
	}

	var routeData []byte
	if msgHead.GetRouteLen() > 0 {
		routeData = make([]byte, msgHead.GetRouteLen())
		if _, err := io.ReadFull(c.Conn, routeData); err != nil {
			Error("read msg data error ", err)
			return nil, err
		}
	}
	msgRecv := &MsgRecv{
		Route: string(routeData),
		MsgID: msgHead.GetMsgID(),
		Body:  data}
	return msgRecv, nil
}

func (c *Client) Exit() {
	time.Sleep(time.Second * 1)
	c.Conn.Close()

}
