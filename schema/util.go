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

	msgPack := &MsgPack{}
	if err = msgpack.Unmarshal(data, msgPack); err != nil {
		return nil, err
	}

	data = []byte(msgPack.MsgContent)
	if msgPack.MsgType == MSG_TYPE_OPEN_CONN_REQUEST {
		msg := & OpenConnRequest{}
		err = msgpack.Unmarshal(data, msg)
		msgPack.Msg = msg

	}else if msgPack.MsgType == MSG_TYPE_OPEN_CONN_RESPONSE {
		msg := & OpenConnResponse{}
		err = msgpack.Unmarshal(data, msg)
		msgPack.Msg = msg
	}

	return msgPack, err
}

func WriteMsg(w io.Writer, msgPack *MsgPack) error {
	data, err := msgpack.Marshal(msgPack.Msg)
	if err != nil {
		return err
	}

	msgPack1 := & MsgPack {
		MsgType: msgPack.MsgType,
		MsgContent: string(data),
	}

	data, err = msgpack.Marshal(msgPack1)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}