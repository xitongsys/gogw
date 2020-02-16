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

const (
	PACKSIZE   = 1024 * 1024
	BUFFERSIZE = 100
)

type Client struct {
	ServerAddr string

	SourceAddr  string
	ToPort int
	Direction string

	Protocol string

	Description string
	TimeoutSecond time.Duration

	LastHeartbeat time.Time
	ClientId schema.ClientId

	CmdHandler func (pack *schema.PackResponse)
}


func (client *Client) Start() {
	logger.Info(fmt.Sprintf("\nclient start\nServer: %v\nSourceAddr: %v\nToPort: %v\nDirection: %v\nProtocol: %v\nDescription: %v\nTimeoutSecond: %v\n", 
	client.ServerAddr, client.SourceAddr, client.ToPort, client.Direction, client.Protocol, client.Description, int(client.TimeoutSecond.Seconds())))

	//start heartbeat
	go client.heartbeat()

	//recv cmd from server
	go client.recvCmdFromServer()

	for {
		if err := client.register(); err != nil {
			logger.Error(err)
			time.Sleep(2 * time.Second)
			continue
		}

		for {
			t := time.Now()
			if t.Sub(client.LastHeartbeat).Milliseconds() > client.TimeoutSecond.Milliseconds() {
				logger.Error("timeout")
				break
			}
			time.Sleep(2 * time.Second)
		}
	}
}

func (client *Client) register() error {
	client.ClientId = schema.ClientId("")

	url := fmt.Sprintf("http://%s/register", client.ServerAddr)
	registerRequest := & schema.RegisterRequest {
		SourceAddr: client.SourceAddr,
		ToPort: client.ToPort,
		Direction: client.Direction,
		Protocol: client.Protocol,
		Description: client.Description,
	}

	data, err := registerRequest.Marshal()
	if err != nil {
		return err
	}

	data, err = client.query(url, data)
	if err != nil {
		return err
	}

	registerResponse := &schema.RegisterResponse{}
	if err := registerResponse.Unmarshal(data); err != nil || registerResponse.Code == schema.FAILED {
		return fmt.Errorf("Register failed")
	}

	client.ClientId = registerResponse.ClientId
	client.LastHeartbeat = time.Now()
	client.ToPort = registerResponse.ToPort
	return nil
}

func (client *Client) sendCmdToServer(connId schema.ConnectionId, cmd string) (*schema.PackResponse, error) {
	packRequest := &schema.PackRequest{
		ClientId: client.ClientId,
		ConnId:   connId,
		Type:     schema.CLIENT_SEND_CMD,
		Content:  cmd,
	}

	var data []byte
	var err error
	if data, err = packRequest.Marshal(); err != nil {
		logger.Error(err)
		return nil, err
	}

	url := fmt.Sprintf("http://%s/pack?clientid=%s", client.ServerAddr, client.ClientId)
	if data, err = client.query(url, data); err != nil {
		return nil, err
	}

	packResponse := &schema.PackResponse{}
	err = packResponse.Unmarshal(data)	
	return packResponse, err
}

func (client *Client) sendToServer(packRequest *schema.PackRequest) (err error) {
	url := fmt.Sprintf("http://%s/pack?clientid=%s", client.ServerAddr, client.ClientId)
	var data []byte
	if data, err = packRequest.Marshal(); err == nil {
		_, err = client.query(url, data)
	}

	return err
}

func (client *Client) recvFromServer(connId schema.ConnectionId) (*schema.PackResponse, error) {
	url := fmt.Sprintf("http://%s/pack?clientid=%s", client.ServerAddr, client.ClientId)
	packRequest := &schema.PackRequest{
		ClientId: client.ClientId,
		ConnId:   connId,
		Type:     schema.CLIENT_REQUEST_PACK,
	}
	var data []byte
	var err error
	if data, err = packRequest.Marshal(); err == nil {
		data, err = client.query(url, data)
		if err == nil {
			packResponse := &schema.PackResponse{}
			if err = packResponse.Unmarshal(data); err == nil {
				return packResponse, nil
			}
		}
	}

	return nil, err
}

//recv open conn cmd
func (client *Client) recvCmdFromServer() error {
	for {
		if client.ClientId != "" {
			url := fmt.Sprintf("http://%s/pack?clientid=%s", client.ServerAddr, client.ClientId)
			packRequest := &schema.PackRequest{
				ClientId: client.ClientId,
				Type:     schema.CLIENT_REQUEST_CMD,
			}

			if data, err := packRequest.Marshal(); err == nil {
				if data, err = client.query(url, data); err == nil {
					packResponse := &schema.PackResponse{}
					if err = packResponse.Unmarshal(data); err == nil {
						client.CmdHandler(packResponse)
					}
				}
			}

		}else{
			time.Sleep(time.Second)
		}
	}
}

func (client *Client) heartbeat() {
	for {
		if client.ClientId != "" {
			url := fmt.Sprintf("http://%s/heartbeat?clientid=%s", client.ServerAddr, client.ClientId)
			data, err := client.query(url, nil)
			if string(data) == "ok" {
				client.LastHeartbeat = time.Now()
			}
			if err != nil {
				logger.Error(err)
			}
		}
		time.Sleep(3 * time.Second)
	}
}

func (client *Client) query(url string, body []byte) ([]byte, error) {
	rep, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	defer rep.Body.Close()
	return ioutil.ReadAll(rep.Body)
}
