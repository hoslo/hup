package hup

type Message struct {
	DataLen  uint32 //消息的长度
	routeLen uint32 //消息的ID
	ID       uint32 //消息的ID
	route    []byte
	Data     []byte //消息的内容
}

//创建一个Message消息包
func NewMsgPackage(route string, msgID uint32, data []byte) *Message {
	return &Message{
		DataLen:  uint32(len(data)),
		routeLen: uint32(len(route)),
		ID:       msgID,
		route:    []byte(route),
		Data:     data,
	}
}

//获取消息数据段长度
func (msg *Message) GetDataLen() uint32 {
	return msg.DataLen
}

func (msg *Message) GetMsgID() uint32 {
	return msg.ID
}

//获取消息ID
func (msg *Message) GetRoute() []byte {
	return msg.route
}

func (msg *Message) GetRouteLen() uint32 {
	return msg.routeLen
}

//获取消息内容
func (msg *Message) GetData() []byte {
	return msg.Data
}

//设置消息数据段长度
func (msg *Message) SetDataLen(len uint32) {
	msg.DataLen = len
}

//设计消息ID
func (msg *Message) SetRoute(route []byte) {
	msg.route = route
}

//设计消息内容
func (msg *Message) SetData(data []byte) {
	msg.Data = data
}

func (msg *Message) SetMsgID(msgId uint32) {
	msg.ID = msgId
}
