package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net"

	"gogw/common/schema"
	"gogw/logger"
)

const (
	PACKSIZE = 10240
	BUFFERSIZE = 100
)

type Client struct {
	ServerAddr string

	LocalAddr string
	RemotePort int

	ClientId schema.ClientId
}

func NewClient(serverAddr string, localAddr string, remotePort int) *Client {
	return & Client {
		ServerAddr: serverAddr,
		LocalAddr: localAddr,
		RemotePort: remotePort,
	}
}

func (client *Client) Start() {
	logger.Info("client start")
	if err := client.register(); err != nil {
		logger.Error(err)
		return
	}

	client.recvCmdFromServer()
}

func (client *Client) register() error {
	url := fmt.Sprintf("%s/register?port=%d", client.ServerAddr, client.RemotePort)
	data, err := client.query(url, nil)
	if err != nil {
		return err
	}

	registerResponse := & schema.RegisterResponse{}
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
		defer func(){
			client.closeConnection(connId, conn)
		}()

		bs := make([]byte, PACKSIZE)
		for {
			if n, err := conn.Read(bs); err == nil && n > 0 {
				packRequest := & schema.PackRequest {
					ClientId: client.ClientId,
					ConnId: connId,
					Type: schema.CLIENT_SEND_PACK,
					Content: string(bs[:n]),
				}

				client.sendToServer(packRequest)


			}else if err != nil {
				logger.Warn(err)
				return
			}
		}
	}()

	//read from server, send to conn
	go func(){
		defer func(){
			client.closeConnection(connId, conn)
		}()

		for {
			packResponse, err := client.recvFromServer(connId)
			if err == nil {
				_, err = io.WriteString(conn, packResponse.Content)

			}
			if err != nil {
				return
			}
		}
	}()

	return nil
}

func (client *Client) closeConnection(connId schema.ConnectionId, conn net.Conn) {
	conn.Close()
}

func (client *Client) sendToServer(packRequest *schema.PackRequest) (err error) {
	url := fmt.Sprintf("%s/pack?clientid=%s", client.ServerAddr, client.ClientId)
	var data []byte
	if data, err = packRequest.Marshal(); err == nil {
		_, err = client.query(url, data)
	}

	return err
}

func (client *Client) recvFromServer(connId schema.ConnectionId) (*schema.PackResponse, error) {
	url := fmt.Sprintf("%s/pack?clientid=%s", client.ServerAddr, client.ClientId)
	packRequest := & schema.PackRequest {
		ClientId: client.ClientId,
		ConnId: connId,
		Type: schema.CLIENT_REQUEST_PACK,
	}

	if data, err := packRequest.Marshal(); err == nil {
		rep, err1 := client.query(url, data)
		if err1 == nil {
			packResponse := & schema.PackResponse{}
			if err2 := packResponse.Unmarshal(rep); err2 == nil {
				return packResponse, nil
			}
		}
	}

	return nil, fmt.Errorf("recv error")
}

func (client *Client) cmdHandler(pack *schema.PackResponse) {
	if pack.Content == schema.CMD_OPEN_CONN {
		client.openConnection(pack.ConnId)
	}
}

func (client *Client) recvCmdFromServer() error {
	url := fmt.Sprintf("%s/pack?clientid=%s", client.ServerAddr, client.ClientId)
	for {
		packRequest := & schema.PackRequest {
			ClientId: client.ClientId,
			Type: schema.CLIENT_REQUEST_CMD,
		}

		if data, err := packRequest.Marshal(); err == nil {
			if data, err = client.query(url, data); err == nil {
				packResponse := & schema.PackResponse {}
				if err = packResponse.Unmarshal(data); err == nil {
					client.cmdHandler(packResponse)
				}
			}
		}
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

