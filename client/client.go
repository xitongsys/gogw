package client

import (
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

func (c *Client) reverseNewConn() {
	url := fmt.Sprintf("http://%v/reversenewconn?client=%v", c.ClientId)
	if conn, err := net.Dial(c.Protocol, c.SourceAddr); err == nil {
		if request, err := http.NewRequest("GET", url, conn); err == nil {
			httpClient := &http.Client{}
			if response, err := httpClient.Do(request); err == nil {
				go func(){
					defer func(){
						if err := recover(); err != nil {
							logger.Error(err)
						}
					}()

					io.Copy(conn, response.Body)
				}()
			}
		}
	}
}

func (c *Client) msgHandler(msg *schema.Msg) {
	if msg.MsgType == schema.MSG_OPEN_CONN {
		c.reverseNewConn()

	}else if msg.MsgType == schema.MSG_SET_CLIENT_ID {
		c.ClientId = msg.MsgContent
	}
}

func (c *Client) register() error {
	url := fmt.Sprintf(
		"http://%v/register?sourceaddr=%v&toport=%v&direction=%v&protocol=%v&description=%v", 
		c.ServerAddr,
		c.SourceAddr,
		c.ToPort,
		c.Direction,
		c.Protocol,
		c.Description,
	)

	request, err := http.NewRequest("GET", url, c.MsgPipeReader)
	if err != nil {
		return err
	}

	httpClient := &http.Client{}
	response, err2 := httpClient.Do(request)
	if err2 != nil {
		return err2
	}

	//read msg 
	for {
		msg, err := common.ReadMsg(response.Body)
		if err != nil {
			return err
		}

		c.msgHandler(msg)
	}
}
