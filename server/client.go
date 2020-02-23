package server

import (
	"fmt"
	"net"
	"sync"
	"time"

	"gogw/common"
	"gogw/logger"
	"gogw/schema"
	"gogw/monitor"
)

type Client struct {
	ClientId    string
	ClientAddr  string
	ToPort      int
	Direction   string
	Protocol    string
	SourceAddr  string
	Description string
	Compress bool

	Conns    *sync.Map
	ConnNumber int
	MsgChann chan *schema.MsgPack

	TCPListener net.Listener
	UDPListener *net.UDPConn
	UDPAddrToConnId map[string]string

	LastHeartbeatTime time.Time
	SpeedMonitor      *monitor.SpeedMonitor
}

func NewClient(
	clientId string,
	clientAddr string,
	toPort int,
	direction string,
	protocol string,
	sourceAddr string,
	description string,
	compress bool,
) *Client {

	return &Client{
		ClientId:    clientId,
		ClientAddr:  clientAddr,
		ToPort:      toPort,
		Direction:   direction,
		Protocol:    protocol,
		SourceAddr:  sourceAddr,
		Description: description,
		Compress: compress,

		Conns:    &sync.Map{},
		ConnNumber: 0,
		MsgChann: make(chan *schema.MsgPack),

		UDPAddrToConnId: make(map[string]string),

		LastHeartbeatTime: time.Now(),
		SpeedMonitor:      monitor.NewSpeedMonitor(),
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

	if c.TCPListener != nil {
		c.TCPListener.Close()
	}

	if c.UDPListener != nil {
		c.UDPListener.Close()
	}
}

func (c *Client) addConn(connId string, conn net.Conn) {
	//just approximate 
	if _, ok := c.Conns.Load(connId); !ok {
		c.ConnNumber++
	}

	c.Conns.Store(connId,
		&common.Conn{
			ConnId: connId,
			Conn:   conn,
		})
}

func (c *Client) deleteConn(connId string) {
	//just approximate
	if _, ok := c.Conns.Load(connId); ok {
		c.ConnNumber--
	}

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

			logger.Info(fmt.Sprintf("New Connection\nClientId: %v\nSourceAddr: %v\nRemoteAddr: %v\n", 
			c.ClientId, c.SourceAddr, conn.RemoteAddr))
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
			n, remoteAddr, err := c.UDPListener.ReadFromUDP(bs)
			if err != nil {
				logger.Error(err)
				return
			}

			addr := remoteAddr.String()
			if connId, ok := c.UDPAddrToConnId[addr]; !ok {
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

				logger.Info(fmt.Sprintf("New Connection\nClientId: %v\nSourceAddr: %v\nRemoteAddr: %v\n", 
				c.ClientId, c.SourceAddr, conn.RemoteAddr))
			}

			connId := c.UDPAddrToConnId[addr]
			if value, ok := c.Conns.Load(connId); ok {
				conn, _ := value.(*common.Conn)
				if udpConn, ok := conn.Conn.(*common.UDPConn); ok {
					udpConn.PipeWriter.Write(bs[:n])
				}
			}
		}
	}()

	return nil
}
