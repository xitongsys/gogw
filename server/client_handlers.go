package server

import (
	"io"
	"net/http"

	"gogw/logger"
	"gogw/schema"
	"gogw/common"
)

func (c *Client) HttpHandler(w http.ResponseWriter, req *http.Request) {
	msgPack, err := schema.ReadMsg(req.Body)
	if err != nil {
		logger.Error(err)
		return
	}

	if msgPack.MsgType == schema.MSG_TYPE_OPEN_CONN_REQUEST {
		msg, _ := msgPack.Msg.(*schema.OpenConnRequest)
		c.openConnHandler(msg, w, req)

	}else if msgPack.MsgType == schema.MSG_TYPE_MSG_COMMON_REQUEST {
		msg := <- c.MsgChann
		schema.WriteMsg(w, msg)
	}
}

func (c *Client) openConnHandler(msg *schema.OpenConnRequest, w http.ResponseWriter, req *http.Request) {
	if msg.Role == schema.ROLE_QUERY_CONNID {
		//TODO: forward client

	}else if msg.Role == schema.ROLE_READER {
		if conni, ok := c.Conns.Load(msg.ConnId); ok {
			conn, _ := conni.(*common.Conn)
			_, err := io.Copy(conn.Conn, req.Body)
			logger.Error(err)
		}	

	}else if msg.Role == schema.ROLE_WRITER {
		if conni, ok := c.Conns.Load(msg.ConnId); ok {
			conn, _ := conni.(*common.Conn)
			_, err := io.Copy(w, conn.Conn)
			logger.Error(err)
		}

	}else {
		logger.Error("Unknown role", msg.Role)
	}
}