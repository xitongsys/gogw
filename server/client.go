package server

import (
	"net"
	"net/http"
	"io/ioutil"
	"fmt"

	"gogw/schema"
	"gogw/logger"
	"gogw/common"
)

type Client struct {
	ClientId schema.ClientId
	Port int
	Conns map[schema.ConnectionId][2]chan schema.PackType
}

func (client *Client) start() {
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
	connId := common.UUID()
	_ = connId
}

func (client *Client) closeConnection(connId schema.ConnectionId) {
	
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

	if packRequest.Type == schema.CLOSE {

	}
}