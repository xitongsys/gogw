package server

import (
	"net"
	"net/http"
	"io/ioutil"
	"fmt"

	"schema"
	"logger"
	"common"
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
	connId, err := util.uuid()
	if err != nil {
		logger.Warn(err)
		return
	}
}

func (client *Client) closeConnection(connId ConnectionId) {
	
}

func (client *Client) requestHandler(w http.ResponseWriter, req *http.Request) {
	bs, err := ioutil.ReadAll(req)
	if err != nil {
		logger.Error(err)
		return
	}

	pack := PackType{}
	if err = pack.Unmarshal(bs); err != nil {
		logger.Error(err)
		return
	}

	if pack.PackType == CLOSE {

	}
}