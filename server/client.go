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
	ClientId    string
	ClientAddr  string
	ToPort      int
	Direction   string
	Protocol    string
	SourceAddr  string
	Description string

	Conns    *sync.Map
	MsgChann chan *schema.MsgPack

	TCPListener net.Listener
	UDPListener *net.UDPConn
	UDPAddrToConnId map[string]string

	LastHeartbeatTime time.Time
	SpeedMonitor      *SpeedMonitor
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

	return &Client{
		ClientId:    clientId,
		ClientAddr:  clientAddr,
		ToPort:      toPort,
		Direction:   direction,
		Protocol:    protocol,
		SourceAddr:  sourceAddr,
		Description: description,

		Conns:    &sync.Map{},
		MsgChann: make(chan *schema.MsgPack),

		UDPAddrToConnId: make(map[string]string),

		LastHeartbeatTime: time.Now(),
		SpeedMonitor:      NewSpeedMonitor(),
	}
}

func (c *Client) Start() (err error) {
	if c.Direction == schema.DIRECTION_REVERSE {
		if c.Protocol == schema.PROTOCOL_TCP {
			return c.startReverseTCPListener()
		}

		if c.Protocol == schema.PROTOCOL_UDP {
			return c.startReverseUDPListener()
		}
	}

	return nil
}

func (c *Client) Stop() {
	close(c.MsgChann)
	c.Conns.Range(func(k, v interface{}) bool {
		conn, _ := v.(*common.Conn)
		conn.Conn.Close()
		return true
	})

	if c.Protocol == schema.PROTOCOL_TCP {
		c.TCPListener.Close()
	}

	if c.Protocol == schema.PROTOCOL_UDP {
		c.UDPListener.Close()
	}
}

func (c *Client) addConn(connId string, conn net.Conn) {
	c.Conns.Store(connId,
		&common.Conn{
			ConnId: connId,
			Conn:   conn,
		})
}

func (c *Client) deleteConn(connId string) {
	c.Conns.Delete(connId)
	if c.UDPAddrToConnId != nil {
		delete(c.UDPAddrToConnId, connId)
	}
}

func (c *Client) startReverseTCPListener() (err error) {
	c.TCPListener, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", c.ToPort))
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := c.TCPListener.Accept()
			if err != nil {
				logger.Error(err)
				return
			}

			connId := common.UUID("connid")
			c.addConn(connId, conn)

			msgPack := &schema.MsgPack{
				MsgType: schema.MSG_TYPE_OPEN_CONN_RESPONSE,
				Msg: &schema.OpenConnResponse{
					ConnId: connId,
					Status: schema.STATUS_SUCCESS,
				},
			}

			c.MsgChann <- msgPack
		}
	}()

	return nil
}

func (c *Client) startReverseUDPListener() (err error) {
	c.UDPListener, err = net.ListenUDP(c.Protocol, &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: c.ToPort})
	if err != nil {
		return err
	}

	go func() {
		bs := make([]byte, PACKSIZE)
		for {
			_, remoteAddr, err := c.UDPListener.ReadFromUDP(bs)
			if err != nil {
				logger.Error(err)
				return
			}

			if connId, ok := c.UDPAddrToConnId[remoteAddr.String()]; !ok {
				connId = common.UUID("connid")
				conn := common.NewUDPConn(remoteAddr, c.UDPListener)
				c.addConn(connId, conn)
				c.UDPAddrToConnId[remoteAddr.String()] = connId

				msgPack := &schema.MsgPack{
					MsgType: schema.MSG_TYPE_OPEN_CONN_RESPONSE,
					Msg: &schema.OpenConnResponse{
						ConnId: connId,
						Status: schema.STATUS_SUCCESS,
					},
				}

				c.MsgChann <- msgPack
			}

			if value, ok := c.Conns.Load(c.UDPAddrToConnId[remoteAddr.String()]); ok {
				conn, _ := value.(*common.Conn)
				conn.Conn.Write(bs)
			}
		}
	}()

	return nil
}
