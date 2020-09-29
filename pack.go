package hup

import (
	"bytes"
	"encoding/binary"
)

type Pack struct{}

func NewPack() *Pack {
	return &Pack{}
}

func (dp *Pack) GetHeadLen() uint32 {
	//route uint32(4字节) +  DataLen uint32(4字节) + MsgID uint32(4字节)
	return 12
}

func (dp *Pack) Pack(msg *Message) ([]byte, error) {

	dataBuff := bytes.NewBuffer([]byte{})

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.DataLen); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.RouteLen); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.ID); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.Data); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.Route); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (dp *Pack) Unpack(binaryData []byte) (*Message, error) {

	dataBuff := bytes.NewReader(binaryData)

	msg := &Message{}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.RouteLen); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}

	return msg, nil
}
