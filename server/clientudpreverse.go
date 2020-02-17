package server

import (
	"net"
	"time"

	"gogw/common"
	"gogw/common/schema"
	"gogw/logger"
)

type ClientUDPReverse struct {
	Client
	Listener *net.UDPConn
	AddrToConn map[string]schema.ConnectionId
	ConnToAddr map[schema.ConnectionId]string
}

func NewClientUDPReverse (clientId schema.ClientId, clientAddr string, toPort int, sourceAddr string, description string) *ClientUDPReverse {
	client := & ClientUDPReverse {
		Client: Client {
			ClientId: clientId,
			ClientAddr: clientAddr,
			ToPort: toPort,
			Direction: schema.DIRECTION_REVERSE,
			Protocol: "udp",
			SourceAddr: sourceAddr,
			Description: description,
			FromClientChanns: make(map[schema.ConnectionId]chan *schema.PackRequest),
			ToClientChanns: make(map[schema.ConnectionId]chan *schema.PackResponse),
			CmdToClientChann: make(chan *schema.PackResponse),
			SpeedMonitor: NewSpeedMonitor(),
			LastHeartbeat: time.Now(),
		},

		ConnToAddr: make(map[schema.ConnectionId]string),
		AddrToConn: make(map[string]schema.ConnectionId),
	}

	client.CmdHandler = client.cmdHandler
	return client
}

func (client *ClientUDPReverse) Start() (err error) {
	client.Listener, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: client.ToPort})
	if err != nil {
		return err
	}
	
	go func() {
		bs := make([]byte, PACKSIZE)
		for {
			n, remoteAddr, err := client.Listener.ReadFromUDP(bs)
			if err != nil {
				logger.Error(err)
				continue
			}

			pack := & schema.PackResponse {
				ClientId: client.ClientId,
				Type: schema.CLIENT_SEND_PACK,
				Content: string(bs[:n]),
			}

			var connId schema.ConnectionId

			key := remoteAddr.String()
			if v, ok := client.AddrToConn[key]; ok {
				connId = v

			}else{
				client.openConnection(key)
				connId = client.AddrToConn[key]
			}

			pack.ConnId = connId
			client.ToClientChanns[connId] <- pack
		}
	}()

	return nil
}

func (client *ClientUDPReverse) Stop() error {
	if client.Listener == nil {
		return nil
	}
	
	return client.Listener.Close()
}

func (client *ClientUDPReverse) openConnection(remoteAddr string) {
	connId := schema.ConnectionId(common.UUID("connid"))
	toChann, fromChann := make(chan *schema.PackResponse, BUFFSIZE), make(chan *schema.PackRequest, BUFFSIZE)
	client.ToClientChanns[connId] = toChann
	client.FromClientChanns[connId] = fromChann
	client.ConnToAddr[connId] = remoteAddr
	client.AddrToConn[remoteAddr] = connId

	openPack := & schema.PackResponse {
		ClientId: client.ClientId,
		ConnId: connId,
		Type: schema.SERVER_CMD,
		Content: schema.CMD_OPEN_CONN,
	}

	client.CmdToClientChann <- openPack

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
				var addr *net.UDPAddr
				var err error

				if addr, err = net.ResolveUDPAddr("udp", client.ConnToAddr[connId]); err == nil {
					_, err = client.Listener.WriteToUDP([]byte(pack.Content), addr)
				}

				//_, err := io.WriteString(client.Listener, pack.Content)
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

func (client *ClientUDPReverse) closeConnection(connId schema.ConnectionId) {
	defer func(){
		if err := recover(); err != nil {
			logger.Warn(err)
		}
	}()
	
	remoteAddr := client.ConnToAddr[connId]
	delete(client.ConnToAddr, connId)
	delete(client.AddrToConn, remoteAddr)

	close(client.ToClientChanns[connId])
	close(client.FromClientChanns[connId])
	delete(client.ToClientChanns, connId)
	delete(client.FromClientChanns, connId)
}

func (client *ClientUDPReverse) cmdHandler(packRequest *schema.PackRequest) *schema.PackResponse {
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
