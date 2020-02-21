package schema

import (
	"io"
	"io/ioutil"

	"github.com/vmihailenco/msgpack/v4"
)

func ReadMsg(r io.Reader) (*MsgPack, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return UnmarshalMsg(data)
}

func UnmarshalMsg(data []byte) (*MsgPack, error) {
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
	data, err := MarshalMsg(msgPack)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

func MarshalMsg(msgPack *MsgPack) ([]byte, error) {
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
