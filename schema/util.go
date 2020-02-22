package schema

import (
	"encoding/binary"
	"io"

	"github.com/vmihailenco/msgpack/v4"
)

//msg with a length(4 byte) header
func ReadMsg(r io.Reader) (*MsgPack, error) {
	lengthHeaderData := make([]byte, 4)
	_, err := io.ReadAtLeast(r, lengthHeaderData, 4)
	if err != nil {
		return nil, err
	}

	length := binary.LittleEndian.Uint32(lengthHeaderData)
	data := make([]byte, length)

	_, err = io.ReadAtLeast(r, data, int(length))
	if err != nil {
		return nil, err
	}

	return unmarshalMsg(data)
}

func unmarshalMsg(data []byte) (*MsgPack, error) {
	msgPack := &MsgPack{}
	var err error
	if err = msgpack.Unmarshal(data, msgPack); err != nil {
		return nil, err
	}

	data = []byte(msgPack.MsgContent)
	if msgPack.MsgType == MSG_TYPE_REGISTER_REQUEST {
		msg := & RegisterRequest{}
		err = msgpack.Unmarshal(data, msg)
		msgPack.Msg = msg

	}else if msgPack.MsgType == MSG_TYPE_REGISTER_RESPONSE {
		msg := & RegisterResponse{}
		err = msgpack.Unmarshal(data, msg)
		msgPack.Msg = msg

	}else if msgPack.MsgType == MSG_TYPE_OPEN_CONN_REQUEST {
		msg := & OpenConnRequest{}
		err = msgpack.Unmarshal(data, msg)
		msgPack.Msg = msg

	}else if msgPack.MsgType == MSG_TYPE_OPEN_CONN_RESPONSE {
		msg := & OpenConnResponse{}
		err = msgpack.Unmarshal(data, msg)
		msgPack.Msg = msg
	}else if msgPack.MsgType == MSG_TYPE_MSG_COMMON_REQUEST {
		//do nothing
	}

	return msgPack, err
}

func WriteMsg(w io.Writer, msgPack *MsgPack) error {
	data, err := marshalMsg(msgPack)
	if err != nil {
		return err
	}

	length := uint32(len(data))
	lengthHeaderData := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthHeaderData, length)
	data = append(lengthHeaderData, data...)
	_, err = w.Write(data)
	return err
}

func marshalMsg(msgPack *MsgPack) ([]byte, error) {
	data, err := msgpack.Marshal(msgPack.Msg)
	if err != nil {
		return nil, err
	}

	msgPack1 := & MsgPack {
		MsgType: msgPack.MsgType,
		MsgContent: string(data),
	}

	return msgpack.Marshal(msgPack1)
}
