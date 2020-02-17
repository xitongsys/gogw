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
	//tcp
	Listener net.Listener

	//udp
	ListenerUDP *net.UDPConn
	AddrToConn map[string]schema.ConnectionId
	ConnToAddr map[schema.ConnectionId]string
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

		AddrToConn: make(map[string]schema.ConnectionId),
		ConnToAddr: make(map[schema.ConnectionId]string),
	}
	client.CmdHandler = client.cmdHandler

	return client
}

func (client *ClientForward) Start() {
	if client.Protocol == "tcp" {
		client.startTCP()

	}else if client.Protocol == "udp" {
		client.startUDP()

	}else{
		logger.Error("Unsupported protocol: ", client.Protocol)
	}
}

func (client *ClientForward) startTCP() {
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

		connId, err := client.getConnectionId()
		if err != nil {
			logger.Error(err)
			conn.Close()
			continue
		}

		if err = client.openConnectionTCP(connId, conn); err != nil {
			logger.Error(err)
		}
	}
}

func (client *ClientForward) startUDP() {
	var err error
	client.ListenerUDP, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: client.ToPort})
	if err != nil {
		logger.Error(err)
		return
	}

	go client.Client.Start()
	
	bs := make([]byte, PACKSIZE)
	for {
		n, remoteAddr, err := client.ListenerUDP.ReadFromUDP(bs)
		if err != nil {
			logger.Error(err)
			continue
		}

		var connId schema.ConnectionId

		key := remoteAddr.String()
		if v, ok := client.AddrToConn[key]; ok {
			connId = v

		}else{
			connId, err = client.getConnectionId()
			if err != nil {
				logger.Error(err)
				continue
			}

			client.AddrToConn[key] = connId
			client.ConnToAddr[connId] = key

			client.openConnectionUDP(connId)
		}

		packRequest := &schema.PackRequest{
			ClientId: client.ClientId,
			ConnId:   connId,
			Type:     schema.CLIENT_SEND_PACK,
			Content:  string(bs[:n]),
		}

		client.sendToServer(packRequest)
	}
}

func (client *ClientForward) getConnectionId() (schema.ConnectionId, error) {
	packResponse, err := client.sendCmdToServer("", schema.CMD_OPEN_CONN)
	if err != nil {
		return "", err
	}

	if packResponse.Code != schema.SUCCESS {
		return "", fmt.Errorf("open connection error from server")
	}

	return packResponse.ConnId, nil
}

func (client *ClientForward) openConnectionTCP(connId schema.ConnectionId, conn net.Conn) error {
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


func (client *ClientForward) openConnectionUDP(connId schema.ConnectionId) error {
	//read from server, send to conn
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Warn(err)
			}

			client.closeConnection(connId, nil)
		}()

		for {
			packResponse, err := client.recvFromServer(connId)

			if err == nil && len(packResponse.Content) > 0 {
				logger.Debug("from server", *packResponse)

				var addr *net.UDPAddr
				var err error

				if addr, err = net.ResolveUDPAddr("udp", client.ConnToAddr[connId]); err == nil {
					_, err = client.ListenerUDP.WriteToUDP([]byte(packResponse.Content), addr)
				}
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
	if conn != nil {
		conn.Close()
	}

	_, err := client.sendCmdToServer(connId, schema.CMD_CLOSE_CONN)
	if err != nil {
		logger.Error(err)
	}
}

func (client *ClientForward) cmdHandler(pack *schema.PackResponse) {
}
