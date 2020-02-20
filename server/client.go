package server 

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
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
		LastHeartbeat: time.Now(),
	}
}

func (c *Client) Start() error {
	msg := &schema.Msg{}

	if c.Direction == schema.DIRECTION_REVERSE {
		if err := c.startReverse(); err != nil {
			return err
		}
	}

	for {
		if err := common.ReadObject(c.MsgConn, msg); err != nil {
			return err
		}

		c.msgHandler(msg)
	}

	return nil
}

func (c *Client) msgHandler(msg *schema.Msg) {
}

func (c *Client) readMsg() (*schema.Msg, error) {
	msg := &schema.Msg{}
	return msg, common.ReadObject(msg)
}

func (c *Client) writeMsg(msg *schema.Msg) error {
	return common.WriteObject(c.MsgConn, msg)
}

func (c *Client) startReverse() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", c.ToPort))
	if err != nil {
		return err
	}
	
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				logger.Warn(err)
				return
			}

			c.writeMsg(&schema.Msg{MsgType: schema.MSG_OPEN_CONN})
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