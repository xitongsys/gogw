package client

import (
	"io"
	"net"
	"time"
	"fmt"

	"gogw/common/schema"
	"gogw/logger"
)

type ClientForward struct {
	Client
	Listener net.Listener
}

func NewClientForward(serverAddr string, sourceAddr string, toPort int, protocol string, description string, timeoutSecond int) *ClientForward {
	client := & ClientForward{
		Client: Client{
			ServerAddr: serverAddr,
			SourceAddr:  sourceAddr,
			ToPort: toPort,
			Direction: schema.DIRECTION_FORWARD,
			Protocol: protocol,
			Description: description,
			TimeoutSecond: time.Duration(timeoutSecond) * time.Second,
			LastHeartbeat: time.Now(),
			ClientId: "",
		},
	}
	client.CmdHandler = client.cmdHandler

	return client
}

func (client *ClientForward) Start() {
	var err error 
	client.Listener, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", client.ToPort))
	if err != nil {
		logger.Error(err)
		return
	}

	go client.Client.Start()
	
	for {
		conn, err := client.Listener.Accept()
		if err != nil {
			logger.Warn(err)
			return
		}

		if err = client.openConnection(conn); err != nil {
			logger.Error(err)
		}
	}
}

func (client *ClientForward) openConnection(conn net.Conn) error {
	packResponse, err := client.sendCmdToServer("", schema.CMD_OPEN_CONN)
	if err != nil {
		return err
	}

	if packResponse.Code != schema.SUCCESS {
		return fmt.Errorf("open connection error from server")
	}

	connId := packResponse.ConnId
	
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

func (client *ClientForward) closeConnection(connId schema.ConnectionId, conn net.Conn) {
	conn.Close()
	_, err := client.sendCmdToServer(connId, schema.CMD_CLOSE_CONN)
	if err != nil {
		logger.Error(err)
	}
}

func (client *ClientForward) cmdHandler(pack *schema.PackResponse) {
}
