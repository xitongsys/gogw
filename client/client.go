package client

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"gogw/common"
	"gogw/logger"
	"gogw/schema"
)

type Client struct {
	ServerAddr string
	SourceAddr  string
	ToPort int
	Direction string
	Protocol string
	Description string
	ClientId string

	Conns *sync.Map
}

func NewClient(
	serverAddr string,
	sourceAddr string,
	toPort int,
	direction string,
	protocol string,
	description string,
) *Client {
	return &Client {
		ServerAddr: serverAddr,
		SourceAddr: sourceAddr,
		ToPort: toPort,
		Direction: direction,
		Protocol: protocol,
		Description: description,
		ClientId: "",
		Conns: &sync.Map{},
	}
}

func (c *Client) Start() {
	logger.Info(fmt.Sprintf("\nclient start\nServer: %v\nSourceAddr: %v\nToPort: %v\nDirection: %v\nProtocol: %v\nDescription: %v\n", 
	c.ServerAddr, c.SourceAddr, c.ToPort, c.Direction, c.Protocol, c.Description))

	//start heartbeat
	go c.heartbeatLoop()

	for {
		if err := c.register(); err != nil {
			logger.Error(err)
			time.Sleep(2 * time.Second)
			continue
		}

		if c.Direction == schema.DIRECTION_FORWARD {
			if c.Protocol == schema.PROTOCOL_TCP {
				c.startForwardTCPListener()
			}
		}

		c.msgRequestLoop()
	}
}

func (c *Client) heartbeatLoop() {
	for {
		if c.ClientId != "" {
			url := fmt.Sprintf("http://%s/heartbeat?clientid=%s", c.ServerAddr, c.ClientId)
			_, err := http.Get(url)
			if err != nil {
				logger.Error(err)
			}
		}
		time.Sleep(3 * time.Second)
	}
}

func (c *Client) register() error {
	logger.Info("start register")
	url := fmt.Sprintf("http://%v/register", c.ServerAddr)
	msgPack := &schema.MsgPack{
		MsgType: schema.MSG_TYPE_REGISTER_REQUEST,
		Msg: & schema.RegisterRequest {
			SourceAddr: c.SourceAddr,
			ToPort: c.ToPort,
			Direction: c.Direction,
			Protocol: c.Protocol,
			Description: c.Description,
		},
	}

	r, w := io.Pipe()
	go func(){
		schema.WriteMsg(w, msgPack)
		w.Close()
	}()

	response, err := http.Post(url, "", r)
	if err != nil {
		return err
	}

	msgPack, err = schema.ReadMsg(response.Body)
	if err != nil {
		logger.Error(err)
		return err
	}

	msg, _ := msgPack.Msg.(*schema.RegisterResponse)
	if msg.Status == schema.STATUS_SUCCESS {
		c.ClientId = msg.ClientId
	}else{
		err = fmt.Errorf("register failed")
	}

	return err
}

func (c *Client) msgRequestLoop(){
	url := fmt.Sprintf("http://%v/msg?clientid=%v", c.ServerAddr, c.ClientId)
	msgPack := &schema.MsgPack{
		MsgType: schema.MSG_TYPE_MSG_COMMON_REQUEST,
	}

	var buf bytes.Buffer
	schema.WriteMsg(&buf, msgPack)
	data := buf.Bytes()

	for {
		response, err := http.Post(url, "", bytes.NewReader(data))
		if err != nil {
			logger.Error(err)
			return
		}

		msgPackResponse, err := schema.ReadMsg(response.Body)
		if msgPackResponse.MsgType == schema.MSG_TYPE_OPEN_CONN_RESPONSE {
			msg := msgPackResponse.Msg.(*schema.OpenConnResponse)
			
			if c.Direction == schema.DIRECTION_REVERSE {
				c.openReverseConn(msg.ConnId)
			}
		}
	}
}

func (c *Client) openConn(connId string, conn net.Conn) error {
	url := fmt.Sprintf("http://%v/msg?clientid=%v", c.ServerAddr, c.ClientId)
	c.Conns.Store(connId, &common.Conn{
		ConnId: connId,
		Conn: conn,
	})

	//conn -> server
	go func(){
		readerMsgPack := &schema.MsgPack{
			MsgType: schema.MSG_TYPE_OPEN_CONN_REQUEST,
			Msg: &schema.OpenConnRequest {
				ConnId: connId,
				Role: schema.ROLE_READER,
			},
		}

		r, w := io.Pipe()
		go func(){
			schema.WriteMsg(w, readerMsgPack)
			io.Copy(w, conn)
		}()

		http.Post(url, "", r)

		logger.Debug("conn->server done")
	}()

	//server -> conn
	go func(){
		writerMsgPack := &schema.MsgPack{
			MsgType: schema.MSG_TYPE_OPEN_CONN_REQUEST,
			Msg: &schema.OpenConnRequest {
				ConnId: connId,
				Role: schema.ROLE_WRITER,
			},
		}

		r, w := io.Pipe()
		go func(){
			schema.WriteMsg(w, writerMsgPack)
			w.Close()
		}()

		response , err := http.Post(url, "", r)
		if err != nil {
			return
		}

		io.Copy(conn, response.Body)
		logger.Debug("server -> conn done")
	}()

	return nil
}

func (c *Client) openReverseConn(connId string) error {
	var conn net.Conn
	var err error
	conn, err = net.Dial(c.Protocol, c.SourceAddr)
	if err != nil {
		return err
	}

	c.openConn(connId, conn)
	return nil
}

func (c *Client) startForwardTCPListener() error {
	url := fmt.Sprintf("http://%v/msg?clientid=%v", c.ServerAddr, c.ClientId)
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

			msgPack := & schema.MsgPack {
				MsgType: schema.MSG_TYPE_OPEN_CONN_REQUEST,
				Msg: & schema.OpenConnRequest{
					Role: schema.ROLE_QUERY_CONNID,
				},
			}

			r, w := io.Pipe()
			go func(){
				schema.WriteMsg(w, msgPack)
				w.Close()
			}()

			response, err := http.Post(url, "", r)
			if err != nil {
				logger.Error(err)
				continue
			}

			msgPack, err = schema.ReadMsg(response.Body)
			if err != nil || msgPack.MsgType != schema.MSG_TYPE_OPEN_CONN_RESPONSE{
				logger.Error(err)
				continue
			}

			msg, ok := msgPack.Msg.(*schema.OpenConnResponse)
			if !ok || msg.Status != schema.STATUS_SUCCESS {
				logger.Error("msg error")
				continue
			}

			connId := msg.ConnId
			c.openConn(connId, conn)
		}
	}()

	return nil
}

