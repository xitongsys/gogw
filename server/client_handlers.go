package server

import (
	"net"
	"net/http"

	"gogw/common"
	"gogw/logger"
	"gogw/schema"
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

	}else{
		logger.Error("Unknown MsgType = ", msgPack.MsgType)
	}
}

func (c *Client) openConnHandler(msg *schema.OpenConnRequest, w http.ResponseWriter, req *http.Request) {
	if msg.Role == schema.ROLE_QUERY_CONNID {
		//Forward client: open a new conn
		msgPack := & schema.MsgPack {
			MsgType: schema.MSG_TYPE_OPEN_CONN_RESPONSE,
			Msg: & schema.OpenConnResponse {
				ConnId: "",
				Status: schema.STATUS_FAILED,
			},
		}

		var conn net.Conn
		var err error
		if conn, err = net.Dial(c.Protocol, c.SourceAddr); err == nil {
			connId := common.UUID("connid")
			msgPack.Msg = & schema.OpenConnResponse {
				ConnId: connId,
				Status: schema.STATUS_SUCCESS,
			}

			c.addConn(connId, conn)
		}

		schema.WriteMsg(w, msgPack)

	}else if msg.Role == schema.ROLE_READER {
		if conni, ok := c.Conns.Load(msg.ConnId); ok {
			conn, _ := conni.(*common.Conn)

			common.Copy(conn.Conn, req.Body, false, c.Compress, c.UploadSpeedMonitor)
			c.deleteConn(msg.ConnId)
		}	

	}else if msg.Role == schema.ROLE_WRITER {
		if conni, ok := c.Conns.Load(msg.ConnId); ok {
			conn, _ := conni.(*common.Conn)

			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")

			common.Copy(w, conn.Conn, c.Compress, false, c.DownloadSpeedMonitor)
			c.deleteConn(msg.ConnId)
		}

	}else {
		logger.Error("Unknown role", msg.Role)
	}
}