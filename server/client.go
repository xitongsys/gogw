package server

import (
	"net"
	"net/http"
	"io"
	"io/ioutil"
	"fmt"

	"gogw/common/schema"
	"gogw/logger"
	"gogw/common"
)

const (
	PACKSIZE = 10240
	BUFFSIZE = 100
)

type Client struct {
	ClientId schema.ClientId
	Port int
	Listener net.Listener
	FromClientChanns map[schema.ConnectionId]chan *schema.PackRequest
	ToClientChanns map[schema.ConnectionId]chan *schema.PackResponse
	CmdToClientChann chan *schema.PackResponse
}

func NewClient(clientId schema.ClientId, port int) *Client {
	return & Client {
		ClientId: clientId,
		Port: port,
		FromClientChanns: make(map[schema.ConnectionId]chan *schema.PackRequest),
		ToClientChanns: make(map[schema.ConnectionId]chan *schema.PackResponse),
		CmdToClientChann: make(chan *schema.PackResponse),
	}
}

func (client *Client) Start() (err error) {
	client.Listener, err = net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%d", client.Port))
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

func (client *Client) Stop() error {
	if client.Listener == nil {
		return nil
	}
	
	return client.Listener.Close()
}

func (client *Client) openConnection(conn net.Conn) {
	connId := schema.ConnectionId(common.UUID())
	toChann, fromChann := make(chan *schema.PackResponse, BUFFSIZE), make(chan *schema.PackRequest, BUFFSIZE)
	client.ToClientChanns[connId] = toChann
	client.FromClientChanns[connId] = fromChann

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
			client.closeConnection(connId, conn)
		}()

		bs := make([]byte, PACKSIZE)
		for {
			if n, err := conn.Read(bs); err == nil && n > 0 {
				pack := & schema.PackResponse {
					ClientId: client.ClientId,
					ConnId: connId,
					Type: schema.NORMAL,
					Content: string(bs[:n]),
				}

				logger.Debug("to client: ", *pack)

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
			client.closeConnection(connId, conn)
		}()

		for {
			pack, ok := <- fromChann
			logger.Debug("from client", *pack)

			if ok {
				_, err := io.WriteString(conn, pack.Content)
				if err != nil {
					return
				}
			}
		}
	}()
}

func (client *Client) closeConnection(connId schema.ConnectionId, conn net.Conn) {
	conn.Close()
	close(client.ToClientChanns[connId])
	close(client.FromClientChanns[connId])
	delete(client.ToClientChanns, connId)
	delete(client.FromClientChanns, connId)
}

func (client *Client) cmdHandler(packRequest *schema.PackRequest) *schema.PackResponse {
	packResponse := & schema.PackResponse{}
	return packResponse
}

func (client *Client) requestHandler(w http.ResponseWriter, req *http.Request) {
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error(err)
		return
	}

	packRequest := &schema.PackRequest{}
	if err = packRequest.Unmarshal(bs); err != nil {
		logger.Error(err)
		return
	}

	if packRequest.Type == schema.CLIENT_SEND_PACK {
		client.FromClientChanns[packRequest.ConnId] <- packRequest

	}else if packRequest.Type == schema.CLIENT_REQUEST_PACK {
		packResponse := <- client.ToClientChanns[packRequest.ConnId]
		data, err := packResponse.Marshal()
		if err != nil {
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
				w.Write(data)
			}
		}
	}
}