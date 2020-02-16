package server

import (
	"fmt"
	"io"
	"net"
	"time"

	"gogw/common"
	"gogw/common/schema"
	"gogw/logger"
)

type ClientTCPReverse struct {
	Client
	Listener net.Listener
	Conns map[schema.ConnectionId]net.Conn
}

func NewClientTCPReverse(clientId schema.ClientId, clientAddr string, portTo int, sourceAddr string, description string) *ClientTCPReverse {
	client := & ClientTCPReverse {
		Client: Client {
			ClientId: clientId,
			ClientAddr: clientAddr,
			PortTo: portTo,
			Protocol: "tcp",
			Direction: schema.DIRECTION_REVERSE,
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

func (client *ClientTCPReverse) Start() (err error) {
	client.Listener, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", client.PortTo))
	if err != nil {
		return err
	}
	
	go func() {
		for {
			conn, err := client.Listener.Accept()
			if err != nil {
				logger.Warn(err)
				return
			}

			client.openConnection(conn)
		}
	}()

	return nil
}

func (client *ClientTCPReverse) Stop() error {
	if client.Listener == nil {
		return nil
	}
	
	return client.Listener.Close()
}

func (client *ClientTCPReverse) openConnection(conn net.Conn) {
	connId := schema.ConnectionId(common.UUID("connid"))
	toChann, fromChann := make(chan *schema.PackResponse, BUFFSIZE), make(chan *schema.PackRequest, BUFFSIZE)
	client.ToClientChanns[connId] = toChann
	client.FromClientChanns[connId] = fromChann
	client.Conns[connId] = conn

	openPack := & schema.PackResponse {
		ClientId: client.ClientId,
		ConnId: connId,
		Type: schema.SERVER_CMD,
		Content: schema.CMD_OPEN_CONN,
	}
	client.CmdToClientChann <- openPack

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
}

func (client *ClientTCPReverse) closeConnection(connId schema.ConnectionId) {
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

func (client *ClientTCPReverse) cmdHandler(packRequest *schema.PackRequest) *schema.PackResponse {
	switch  packRequest.Content {
	case schema.CMD_CLOSE_CONN:
		connId := packRequest.ConnId
		client.closeConnection(connId)
		
	}

	packResponse := & schema.PackResponse{
		Code: schema.SUCCESS,
	}
	return packResponse
}
