package client

import (
	"io"
	"net"
	"time"

	"gogw/common/schema"
	"gogw/logger"
)

type ClientReverse struct {
	Client
}

func NewClientReverse(serverAddr string, sourceAddr string, toPort int, protocol string, description string, timeoutSecond int) *ClientReverse {
	client := & ClientReverse{
		Client: Client{
			ServerAddr: serverAddr,
			SourceAddr:  sourceAddr,
			ToPort: toPort,
			Direction: schema.DIRECTION_REVERSE,
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

func (client *ClientReverse) Start() {
	client.Client.Start()
}

func (client *ClientReverse) openConnection(connId schema.ConnectionId) error {
	var conn net.Conn
	var err error
	conn, err = net.Dial(client.Protocol, client.SourceAddr)
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

func (client *ClientReverse) closeConnection(connId schema.ConnectionId, conn net.Conn) {
	conn.Close()
	_, err := client.sendCmdToServer(connId, schema.CMD_CLOSE_CONN)
	if err != nil {
		logger.Error(err)
	}
}

func (client *ClientReverse) cmdHandler(pack *schema.PackResponse) {
	if pack.Content == schema.CMD_OPEN_CONN {
		client.openConnection(pack.ConnId)
	}
}
