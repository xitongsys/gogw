package client

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"gogw/common"
	"gogw/logger"
	"gogw/schema"
)

type Client struct {
	ServerAddr string
	SourceAddr  string
	ToPort int
	Direction string
	Protocol string
	Description string
	ClientId string

	Conns *sync.Map
}

func NewClient(
	serverAddr string,
	sourceAddr string,
	toPort int,
	direction string,
	protocol string,
	description string,
) *Client {
	return &Client {
		ServerAddr: serverAddr,
		SourceAddr: sourceAddr,
		ToPort: toPort,
		Direction: direction,
		Protocol: protocol,
		Description: description,
		ClientId: "",
		Conns: &sync.Map{},
	}
}

func (c *Client) Start() {
	logger.Info(fmt.Sprintf("\nclient start\nServer: %v\nSourceAddr: %v\nToPort: %v\nDirection: %v\nProtocol: %v\nDescription: %v\n", 
	c.ServerAddr, c.SourceAddr, c.ToPort, c.Direction, c.Protocol, c.Description))

	for {
		if err := c.register(); err != nil {
			logger.Error(err)
			time.Sleep(2 * time.Second)
			continue
		}
	}
}

func (c *Client) register() error {
	url := fmt.Sprintf("http://%v/register", c.ServerAddr)
	msgPack := &schema.MsgPack{
		MsgType: schema.MSG_TYPE_REGISTER_REQUEST,
		Msg: & schema.RegisterRequest {
			SourceAddr: c.SourceAddr,
			ToPort: c.ToPort,
			Direction: c.Direction,
			Protocol: c.Protocol,
			Description: c.Description,
		},
	}

	data, err := schema.MarshalMsg(msgPack)
	if err != nil {
		return err
	}

	response, err := http.Post(url, "", bytes.NewReader(data))
	msgPack, err = schema.ReadMsg(response.Body)
	if err != nil {
		return err
	}

	msg, _ := msgPack.Msg.(*schema.RegisterResponse)
	if msg.Status == schema.STATUS_SUCCESS {
		c.ClientId = msg.ClientId
	}else{
		err = fmt.Errorf("register failed")
	}

	return err
}

func (c *Client) msgRequestLoop(){
	url := fmt.Sprintf("http://%v/msg?clientid=%v", c.ServerAddr, c.ClientId)
	msgPack := &schema.MsgPack{
		MsgType: schema.MSG_TYPE_MSG_COMMON_REQUEST,
	}

	data, _ := schema.MarshalMsg(msgPack)

	for {
		response, err := http.Post(url, "", bytes.NewReader(data))
		if err != nil {
			logger.Error(err)
			return
		}

		msgPackResponse, err := schema.ReadMsg(response.Body)
		
	}
}
