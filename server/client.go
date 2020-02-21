package server 

import (
	"fmt"
	"net"
	"sync"
	"time"

	"gogw/common"
	"gogw/logger"
	"gogw/schema"
)

type Client struct {
	ClientId string
	ClientAddr string
	ToPort int
	Direction string
	Protocol string
	SourceAddr string
	Description string

	Conns *sync.Map
	MsgChann chan *schema.MsgPack
	
	LastHeartbeatTime time.Time
	SpeedMonitor *SpeedMonitor
}

func NewClient(
	clientId string, 
	clientAddr string,
	toPort int,
	direction string,
	protocol string,
	sourceAddr string,
	description string,
	) *Client {

	return & Client {
		ClientId: clientId,
		ClientAddr: clientAddr,
		ToPort: toPort,
		Direction: direction,
		Protocol: protocol,
		SourceAddr: sourceAddr,
		Description: description,

		Conns: &sync.Map{},
		MsgChann: make(chan *schema.MsgPack),

		LastHeartbeatTime: time.Now(),
		SpeedMonitor: NewSpeedMonitor(),
	}
}

func (c *Client) Start() (err error) {
	if c.Direction == schema.DIRECTION_REVERSE && 
		c.Protocol == schema.PROTOCOL_TCP {

		return c.startReverseTCPListener()
	}

	return nil
}

func (c *Client) Stop() {
	close(c.MsgChann)
	c.Conns.Range(func (k,v interface{}) bool {
		conn, _ := v.(*common.Conn)
		conn.Conn.Close()
		return true
	})
}

func (c *Client) addConn(connId string, conn net.Conn) {
	c.Conns.Store(connId, 
		&common.Conn{
		ConnId: connId,
		Conn: conn,
	})
}

func (c *Client) deleteConn(connId string) {
	c.Conns.Delete(connId)
}

func (c *Client) startReverseTCPListener() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", c.ToPort))
	if err != nil {
		return err
	}
	
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				logger.Error(err)
				return
			}

			connId := common.UUID("connid")
			c.addConn(connId, conn)

			msgPack := & schema.MsgPack {
				MsgType: schema.MSG_TYPE_OPEN_CONN_RESPONSE,
				Msg: & schema.OpenConnResponse{
					ConnId: connId,
					Status: schema.STATUS_SUCCESS,
				},
			}

			//TODO: add cleaner to avoid block here
			c.MsgChann <- msgPack
		}
	}()

	return nil
}
