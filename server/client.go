package server

import (
	"net"
	"net/http"
	"io/ioutil"

	"github.com/xitongsys/gogw/schema"
	"github.com/xitongsys/gogw/logger"
	"github.com/xitongys/gogw/common"
)

type Client struct {
	ClientId schema.ClientId
	Port int
	Conns map[schema.ConnectionId][2]chan schema.PackType
}

func (client *Cient) start() {
	l, err := net.Listen("tcp4", client.Port)
	if err != nil {
		logger.Error(err)
		return
	}
	defer l.close()

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