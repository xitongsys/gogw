package main

import (
	"fmt"
	"bytes"
	"gogw/schema"
)

func main(){
	msgPack := & schema.MsgPack {
		MsgType: schema.MSG_TYPE_REGISTER_RESPONSE,
		Msg: & schema.RegisterResponse {
			ClientId: "ABC",
			Status: schema.STATUS_SUCCESS,
		},
	}

	var b bytes.Buffer
	schema.WriteMsg(&b, msgPack)

	fmt.Println("=====", b.String())


	msgPack, err := schema.ReadMsg(&b)

	fmt.Println(msgPack, err)
}
