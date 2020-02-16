package server

import (
	"io"
	"net"
	"time"

	"gogw/common"
	"gogw/common/schema"
	"gogw/logger"
)

type ClientTCP struct {
	Client
	Conns map[schema.ConnectionId]net.Conn
}

func NewClientTCP(clientId schema.ClientId, clientAddr string, portTo int, sourceAddr string, description string) *ClientTCP {
	client := & ClientTCP {
		Client: Client {
			ClientId: clientId,
			ClientAddr: clientAddr,
			PortTo: portTo,
			Protocol: "tcp",
			Direction: schema.DIRECTION_FORWARD,
			SourceAddr: sourceAddr,
			Description: description,
			FromClientChanns: make(map[schema.ConnectionId]chan *schema.PackRequest),
			ToClientChanns: make(map[schema.ConnectionId]chan *schema.PackResponse),
			CmdToClientChann: make(chan *schema.PackResponse),
			SpeedMonitor: NewSpeedMonitor(),
			LastHeartbeat: time.Now(),
		},

		Conns: make(map[schema.ConnectionId]net.Conn),
	}

	client.CmdHandler = client.cmdHandler

	return client
}

func (client *ClientTCP) Start() (err error) {
	return nil
}

func (client *ClientTCP) Stop() error {
	connIds := []schema.ConnectionId{}
	for connId, _ := range client.ToClientChanns {
		connIds = append(connIds, connId)
	}

	for _, connId := range connIds {
		client.closeConnection(connId)
	}

	return nil
}

func (client *ClientTCP) openConnection() *schema.PackResponse {
	openPack := & schema.PackResponse {
		Code: schema.FAILED,
	}

	var conn net.Conn
	var err error
	conn, err = net.Dial(client.Protocol, client.SourceAddr)
	if err != nil {
		return openPack
	}

	connId := schema.ConnectionId(common.UUID("connid"))
	toChann, fromChann := make(chan *schema.PackResponse, BUFFSIZE), make(chan *schema.PackRequest, BUFFSIZE)
	client.ToClientChanns[connId] = toChann
	client.FromClientChanns[connId] = fromChann
	client.Conns[connId] = conn

	openPack = & schema.PackResponse {
		ClientId: client.ClientId,
		ConnId: connId,
		Type: schema.SERVER_CMD,
		Content: schema.CMD_OPEN_CONN,
		Code: schema.SUCCESS,
	}

	//read from conn, send to client
	go func(){
		defer func(){
			client.closeConnection(connId)
			if err := recover(); err != nil {
				logger.Warn(err)
			}
		}()

		bs := make([]byte, PACKSIZE)
		for {
			if n, err := conn.Read(bs); err == nil && n > 0 {
				pack := & schema.PackResponse {
					ClientId: client.ClientId,
					ConnId: connId,
					Type: schema.CLIENT_SEND_PACK,
					Content: string(bs[:n]),
					Code: schema.SUCCESS,
				}

				toChann <- pack

			}else if err != nil {
				logger.Warn(err)
				return
			}
		}
	}()

	//read from client, send to conn
	go func(){
		defer func() {
			client.closeConnection(connId)
			if err := recover(); err != nil {
				logger.Warn(err)
			}
		}()

		for {
			pack, ok := <- fromChann

			if ok {
				_, err := io.WriteString(conn, pack.Content)
				if err != nil {
					logger.Warn(err)
					return
				}

			}else{
				return
			}
		}
	}()

	return openPack
}

func (client *ClientTCP) closeConnection(connId schema.ConnectionId) {
	defer func(){
		if err := recover(); err != nil {
			logger.Warn(err)
		}
	}()

	client.Conns[connId].Close()
	delete(client.Conns, connId)

	close(client.ToClientChanns[connId])
	close(client.FromClientChanns[connId])
	delete(client.ToClientChanns, connId)
	delete(client.FromClientChanns, connId)
}

func (client *ClientTCP) cmdHandler(packRequest *schema.PackRequest) *schema.PackResponse {
	switch  packRequest.Content {
	case schema.CMD_CLOSE_CONN:
		connId := packRequest.ConnId
		client.closeConnection(connId)
	case schema.CMD_OPEN_CONN:
		return client.openConnection()
	}

	packResponse := & schema.PackResponse{
		Code: schema.SUCCESS,
	}

	return packResponse
}
