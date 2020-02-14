package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
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

	LocalAddr  string
	RemotePort int

	Description string
	TimeoutSecond time.Duration

	LastHeartbeat time.Time
	ClientId schema.ClientId
}

func NewClient(serverAddr string, localAddr string, remotePort int, description string, timeoutSecond int) *Client {
	return &Client{
		ServerAddr: serverAddr,
		LocalAddr:  localAddr,
		RemotePort: remotePort,
		Description: description,
		TimeoutSecond: time.Duration(timeoutSecond) * time.Second,
		LastHeartbeat: time.Now(),
		ClientId: "",
	}
}

func (client *Client) Start() {
	logger.Info(fmt.Sprintf("\nclient start\nServer: %v\nLocal: %v\nRemotePort: %v\nDescription: %v\nTimeoutSecond: %v\n", 
	client.ServerAddr, client.LocalAddr, client.RemotePort, client.Description, int(client.TimeoutSecond.Seconds())))

	//start heartbeat
	go client.heartbeat()

	go client.recvCmdFromServer()

	for {
		if err := client.register(); err != nil {
			logger.Error(err)
			time.Sleep(5 * time.Second)
			continue
		}

		for {
			t := time.Now()
			if t.Sub(client.LastHeartbeat).Milliseconds() > client.TimeoutSecond.Milliseconds() {
				logger.Error("timeout")
				break
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func (client *Client) register() error {
	client.ClientId = schema.ClientId("")
	
	url := fmt.Sprintf("http://%s/register", client.ServerAddr)
	registerRequest := & schema.RegisterRequest {
		SourceAddr: client.LocalAddr,
		ToPort: client.RemotePort,
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
	return nil
}

func (client *Client) openConnection(connId schema.ConnectionId) error {
	conn, err := net.Dial("tcp", client.LocalAddr)
	if err != nil {
		return err
	}

	//read from conn, send to server
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Warn(err)
			}

			client.closeConnection(connId, conn)
		}()

		bs := make([]byte, PACKSIZE)
		for {
			if n, err := conn.Read(bs); err == nil && n > 0 {
				packRequest := &schema.PackRequest{
					ClientId: client.ClientId,
					ConnId:   connId,
					Type:     schema.CLIENT_SEND_PACK,
					Content:  string(bs[:n]),
				}

				logger.Debug("to server", *packRequest)

				client.sendToServer(packRequest)

			} else if err != nil {
				logger.Warn(err)
				return
			}
		}
	}()

	//read from server, send to conn
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Warn(err)
			}

			client.closeConnection(connId, conn)
		}()

		for {
			packResponse, err := client.recvFromServer(connId)

			if err == nil && len(packResponse.Content) > 0 {

				logger.Debug("from server", *packResponse)

				_, err = io.WriteString(conn, packResponse.Content)

			}

			if err != nil {
				logger.Warn(err)
				return
			}
		}
	}()

	return nil
}

func (client *Client) closeConnection(connId schema.ConnectionId, conn net.Conn) {
	conn.Close()
	err := client.sendCmdToServer(connId, schema.CMD_CLOSE_CONN)
	if err != nil {
		logger.Error(err)
	}
}

func (client *Client) sendCmdToServer(connId schema.ConnectionId, cmd string) (err error) {
	packRequest := &schema.PackRequest{
		ClientId: client.ClientId,
		ConnId:   connId,
		Type:     schema.CLIENT_SEND_CMD,
		Content:  cmd,
	}

	var data []byte
	if data, err = packRequest.Marshal(); err != nil {
		logger.Error(err)
		return err
	}

	url := fmt.Sprintf("http://%s/pack?clientid=%s", client.ServerAddr, client.ClientId)
	_, err = client.query(url, data)
	return err
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

func (client *Client) cmdHandler(pack *schema.PackResponse) {
	if pack.Content == schema.CMD_OPEN_CONN {
		client.openConnection(pack.ConnId)
	}
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
						client.cmdHandler(packResponse)
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
