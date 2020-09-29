package hup

type Message struct {
	DataLen  uint32
	RouteLen uint32
	ID       uint32
	Route    []byte
	Data     []byte
}

func NewMessage(route string, msgID uint32, data []byte) *Message {
	return &Message{
		DataLen:  uint32(len(data)),
		RouteLen: uint32(len(route)),
		ID:       msgID,
		Route:    []byte(route),
		Data:     data,
	}
}
