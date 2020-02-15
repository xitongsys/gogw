package server

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"gogw/common"
	"gogw/common/schema"
	"gogw/logger"
)

type ClientUDP struct {
	ClientId schema.ClientId
	ClientAddr string
	PortTo int
	Protocol string
	SourceAddr string
	Description string

	Listener *net.UDPConn
	FromClientChanns map[schema.ConnectionId]chan *schema.PackRequest
	ToClientChanns map[schema.ConnectionId]chan *schema.PackResponse

	AddrToConn map[string]schema.ConnectionId
	ConnToAddr map[schema.ConnectionId]string


	CmdToClientChann chan *schema.PackResponse

	SpeedMonitor *SpeedMonitor
	LastHeartbeat time.Time
}

func NewClientUDP(clientId schema.ClientId, clientAddr string, portTo int, sourceAddr string, description string) *ClientUDP {
	return & ClientUDP {
		ClientId: clientId,
		ClientAddr: clientAddr,
		PortTo: portTo,
		Protocol: "udp",
		SourceAddr: sourceAddr,
		Description: description,
		FromClientChanns: make(map[schema.ConnectionId]chan *schema.PackRequest),
		ToClientChanns: make(map[schema.ConnectionId]chan *schema.PackResponse),
		ConnToAddr: make(map[schema.ConnectionId]string),
		AddrToConn: make(map[string]schema.ConnectionId),
		CmdToClientChann: make(chan *schema.PackResponse),
		SpeedMonitor: NewSpeedMonitor(),
		LastHeartbeat: time.Now(),
	}
}

func (client *ClientUDP) Start() (err error) {
	client.Listener, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: client.PortTo})
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

			key := remoteAddr.String()
			if connId, ok := client.AddrToConn[key]; ok {
				pack := & schema.PackResponse {
					ClientId: client.ClientId,
					ConnId: connId,
					Type: schema.CLIENT_SEND_PACK,
					Content: string(bs[:n]),
				}

				client.ToClientChanns[connId] <- pack

			}else{
				client.openConnection(key)
			}
		}
	}()

	return nil
}

func (client *ClientUDP) Stop() error {
	if client.Listener == nil {
		return nil
	}
	
	return client.Listener.Close()
}

func (client *ClientUDP) GetClientId() schema.ClientId {
	return client.ClientId
}

func (client *ClientUDP) GetClientAddr() string {
	return client.ClientAddr
}

func (client *ClientUDP) GetPortTo() int {
	return client.PortTo
}

func (client *ClientUDP) GetProtocol() string {
	return client.Protocol
}

func (client *ClientUDP) GetSourceAddr() string {
	return client.SourceAddr
}

func (client *ClientUDP) GetDescription() string {
	return client.Description
}

func (client *ClientUDP) GetConnectionNumber() int {
	return len(client.ConnToAddr)
}

func (client *ClientUDP) GetSpeedMonitor() *SpeedMonitor {
	return client.SpeedMonitor
}

func (client *ClientUDP) GetLastHeartbeat() time.Time {
	return client.LastHeartbeat
}

func (client *ClientUDP) SetLastHeartbeat(t time.Time) {
	client.LastHeartbeat = t
}

func (client *ClientUDP) openConnection(remoteAddr string) {
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
				_, err := io.WriteString(client.Listener, pack.Content)
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

func (client *ClientUDP) closeConnection(connId schema.ConnectionId) {
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

func (client *ClientUDP) cmdHandler(packRequest *schema.PackRequest) *schema.PackResponse {
	switch  packRequest.Content {
	case schema.CMD_CLOSE_CONN:
		connId := packRequest.ConnId
		client.closeConnection(connId)
		
	}

	packResponse := & schema.PackResponse{}
	return packResponse
}

func (client *ClientUDP) RequestHandler(w http.ResponseWriter, req *http.Request) {
	defer func(){
		if err := recover(); err != nil {
			logger.Warn(err)
		}
	}()

	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Debug("from client ", string(bs))
	client.SpeedMonitor.Add(-1, int64(len(bs)))

	packRequest := &schema.PackRequest{}
	if err = packRequest.Unmarshal(bs); err != nil {
		logger.Error(err)
		return
	}

	if packRequest.Type == schema.CLIENT_SEND_PACK {
		client.FromClientChanns[packRequest.ConnId] <- packRequest

	}else if packRequest.Type == schema.CLIENT_REQUEST_PACK {
		if packResponse, ok := <- client.ToClientChanns[packRequest.ConnId]; ok {
			data, _ := packResponse.Marshal()

			logger.Debug("to client", string(data))
			client.SpeedMonitor.Add(int64(len(data)), -1)

			w.Write(data)
		}

	}else if packRequest.Type == schema.CLIENT_SEND_CMD {
		packResponse := client.cmdHandler(packRequest)
		data, err := packResponse.Marshal()
		if err != nil {
			w.Write(data)
		}

	}else if packRequest.Type == schema.CLIENT_REQUEST_CMD {
		select {
		case packResponse := <- client.CmdToClientChann:
			if data, err := packResponse.Marshal(); err == nil {

				logger.Debug("to client", string(data))
				client.SpeedMonitor.Add(int64(len(data)), -1)

				w.Write(data)
			}
		}
	}
}