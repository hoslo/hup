package hup

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"sync"
)

type Connection struct {
	Conn        *net.TCPConn
	MsgHandler  map[string]HandlerFunc
	msgBuffChan chan []byte
	ctx         context.Context
	cancel      context.CancelFunc
	sync.RWMutex
	isClosed bool
}

func NewConnection(conn *net.TCPConn, msgHandler map[string]HandlerFunc) *Connection {
	//初始化Conn属性
	c := &Connection{
		Conn:        conn,
		isClosed:    false,
		MsgHandler:  msgHandler,
		msgBuffChan: make(chan []byte, 1024),
	}

	//将新创建的Conn添加到链接管理中

	return c
}

func (c *Connection) SendMsgBuffChan(route string, msgID uint32, msgData []byte) error {
	//c.Lock()
	if c.isClosed == true {
		return errors.New("Connection closed when send buff msg")
	}
	//c.RUnlock()

	//将data封包，并且发送
	dp := NewPack()
	msg, err := dp.Pack(NewMessage(route, msgID, msgData))
	if err != nil {
		Error("Pack error msg route = ", route)
		return errors.New("Pack error msg ")
	}

	//写回客户端
	c.msgBuffChan <- msg

	return nil
}

func (c *Connection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())

	go c.Read()
	go c.Write()

}

func (c *Connection) Stop() {
	Info("Conn Stop:", c.Conn.RemoteAddr())
	//如果当前链接已经关闭
	c.Lock()
	defer c.Unlock()

	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 关闭socket链接
	c.Conn.Close()
	//关闭Writer
	c.cancel()

	//关闭该链接全部管道
	close(c.msgBuffChan)
}

func (c *Connection) Read() {
	Debug("read goroutine is running")
	defer Debug("conn read exit")
	defer c.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			dp := NewPack()

			headData := make([]byte, dp.GetHeadLen())
			if _, err := io.ReadFull(c.Conn, headData); err != nil {
				if err.Error() == "EOF" {
					return
				} else if strings.Contains(err.Error(), "wsarecv: An existing connection was forcibly closed by the remote host") {
					return
				}
				Error("read msg head err:", err)
				return
			}

			msg, err := dp.Unpack(headData)
			if err != nil {
				Error("unpack err:", err)
			}

			var data []byte
			if msg.DataLen > 0 {
				data = make([]byte, msg.DataLen)
				if _, err := io.ReadFull(c.Conn, data); err != nil {
					Error("read msg data error ", err)
					return
				}
			}

			var routeData []byte
			if msg.RouteLen > 0 {
				routeData = make([]byte, msg.RouteLen)
				if _, err := io.ReadFull(c.Conn, routeData); err != nil {
					Error("read msg data error ", err)
					return
				}
			}
			msg.Data = data
			msg.Route = routeData
			if HandFunc, ok := c.MsgHandler[string(msg.Route)]; ok {
				go HandFunc(c, msg)
			} else {
				return
			}

		}
	}
}

func (c *Connection) Write() {
	Debug("[Writer Goroutine is running]")
	defer Debug(c.Conn.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					Error("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				Debug("msgBuffChan is Closed")
				break
			}
		case <-c.ctx.Done():
			return
		}
	}
}
