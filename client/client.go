package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gogw/common/schema"
	"gogw/logger"
)

type Client struct {
	ServerAddr string
	SourceAddr  string
	ToPort int
	Direction string
	Protocol string
	Description string
	ClientId schema.ClientId

	MsgReader io.Reader
	MsgWriter io.Writer
}


func (client *Client) Start() {
	logger.Info(fmt.Sprintf("\nclient start\nServer: %v\nSourceAddr: %v\nToPort: %v\nDirection: %v\nProtocol: %v\nDescription: %v\nTimeoutSecond: %v\n", 
	client.ServerAddr, client.SourceAddr, client.ToPort, client.Direction, client.Protocol, client.Description, int(client.TimeoutSecond.Seconds())))

	for {
		if err := client.register(); err != nil {
			logger.Error(err)
			time.Sleep(2 * time.Second)
			continue
		}
	}
}

func (client *Client) register() error {
	url := fmt.Sprintf(
		"http://%v/register?sourceaddr=%v&toport=%v&direction=%v&protocol=%v&description=%v", 
		client.ServerAddr,
		client.SourceAddr,
		client.ToPort,
		client.Direction,
		client.Protocol,
		clinet.Description,
	)

	request, err := http.NewRequest("GET", url, )
	if err != nil {
		return err
	}

	client := &http.Client()
	response, err2 := client.Do(request)
	if err2 != nil {
		return err2
	}


	

}
