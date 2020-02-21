package server 

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"

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

	Conns chan net.Conn
	MsgConn *common.HttpConn

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
	w http.ResponseWriter,
	r *http.Request,
	) *Client {

	return & Client {
		ClientId: clientId,
		ClientAddr: clientAddr,
		ToPort: toPort,
		Direction: direction,
		Protocol: protocol,
		SourceAddr: sourceAddr,
		Description: description,

		Conns: make(chan net.Conn),
		MsgConn: common.NewHttpConn(r, w),

		SpeedMonitor: NewSpeedMonitor(),
	}
}

func (c *Client) Start() (err error) {
	msg := &schema.Msg{
		MsgType: schema.MSG_SET_CLIENT_ID,
		MsgContent: c.ClientId,
	}

	//send client id to client
	if err = common.WriteMsg(c.MsgConn, msg); err != nil {
		return err
	}

	if c.Direction == schema.DIRECTION_REVERSE {
		if err := c.startReverseListener(); err != nil {
			return err
		}
	}

	//read msg from client
	for {
		if msg, err = common.ReadMsg(c.MsgConn); err != nil {
			return err
		}

		c.msgHandler(msg)
	}

	return nil
}

func (c *Client) msgHandler(msg *schema.Msg) {
}

func (c *Client) startReverseListener() error {
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

			if err = common.WriteMsg(c.MsgConn, &schema.Msg{MsgType: schema.MSG_OPEN_CONN}); err != nil {
				logger.Error(err)
				return
			}
			c.Conns <- conn
		}
	}()

	return nil
}

func (c *Client) ReverseNewConnHandler(w http.ResponseWriter, r *http.Request){
	conn := <- c.Conns
	var wg sync.WaitGroup
	
	wg.Add(1)
	go func(){
		defer func(){
			if err := recover(); err != nil {
				logger.Error(err)
			}
			wg.Done()
		}()

		io.Copy(conn, r.Body)
	}()

	wg.Add(1)
	go func(){
		defer func(){
			if err := recover(); err != nil {
				logger.Error(err)
			}
			wg.Done()
		}()

		io.Copy(w, conn)
	}()

	wg.Wait()
}