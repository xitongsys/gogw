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
	fromClientChanns map[schema.ConnectionId]chan *schema.PackRequest
	toClientChanns map[schema.ConnectionId]chan *schema.PackResponse
}

func (client *Client) Start() {
	l, err := net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%d", client.Port))
	if err != nil {
		logger.Error(err)
		return
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Warn(err)
			continue
		}

		client.openConnection(conn)
	}
}

func (client *Client) openConnection(conn net.Conn) {
	connId := schema.ConnectionId(common.UUID())
	toChann, fromChann := make(chan *schema.PackResponse, BUFFSIZE), make(chan *schema.PackRequest, BUFFSIZE)
	client.toClientChanns[connId] = toChann
	client.fromClientChanns[connId] = fromChann

	//read from conn, send to client
	go func(){
		defer func(){
			client.closeConnection(connId, conn)
		}()

		openPack := & schema.PackResponse {
			ClientId: client.ClientId,
			ConnId: connId,
			Type: schema.OPEN,
		}

		toChann <- openPack

		bs := make([]byte, PACKSIZE)
		for {
			if n, err := conn.Read(bs); err == nil && n > 0 {
				pack := & schema.PackResponse {
					ClientId: client.ClientId,
					ConnId: connId,
					Type: schema.NORMAL,
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
			client.closeConnection(connId, conn)
		}()

		for {
			pack, ok := <- fromChann
			if ok {
				n, err := io.WriteString(conn, pack.Content)
				if err != nil {
					return
				}
			}
		}
	}()
}

func (client *Client) closeConnection(connId schema.ConnectionId, conn net.Conn) {
	conn.Close()
	close(client.toClientChanns[connId])
	close(client.fromClientChanns[connId])
	delete(client.toClientChanns, connId)
	delete(client.fromClientChanns, connId)
}

func (client *Client) requestHandler(w http.ResponseWriter, req *http.Request) {
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error(err)
		return
	}

	packRequest := schema.PackRequest{}
	if err = packRequest.Unmarshal(bs); err != nil {
		logger.Error(err)
		return
	}

	if packRequest.Type == schema.CLIENTPACK {
		client.fromClientChanns[packRequest.ConnId] <- &packRequest

	}else if packRequest.Type == schema.CLIENTREQUEST {
		packResponse := <- client.toClientChanns[packRequest.ConnId]
		data, err := packResponse.Marshal()
		if err != nil {
			w.Write(data)
		}
	}
}