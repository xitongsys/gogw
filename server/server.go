package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"gogw/common"
	"gogw/logger"
	"gogw/schema"
)

type Server struct {
	ServerAddr    string
	Clients *sync.Map

	TimeoutSecond time.Duration
}

func NewServer(serverAddr string, timeoutSecond int) *Server {
	return &Server{
		ServerAddr:    serverAddr,
		Clients:       &sync.Map{},
		TimeoutSecond: time.Second * time.Duration(timeoutSecond),
	}
}

//client register
func (s *Server) registerHandler(w http.ResponseWriter, req *http.Request) {
	defer func(){
		req.Body.Close()
	}()

	msgPack, err := schema.ReadMsg(req.Body)
	if err != nil {
		logger.Error(err)
		return
	}

	msg, ok := msgPack.Msg.(*schema.RegisterRequest)
	if ! ok {
		return
	}

	clientId := common.UUID("clientid")

	client := NewClient(
		clientId,
		req.RemoteAddr,
		msg.ToPort,
		msg.Direction,
		msg.Protocol,
		msg.SourceAddr,
		msg.Description,
		msg.Compress,
		msg.HttpVersion,
	)

	s.Clients.Store(clientId, client)
	defer func(){
		if err != nil {
			s.Clients.Delete(clientId)
		}
	}()

	if err = client.Start(); err != nil {
		return
	}

	msgPack = & schema.MsgPack {
		MsgType: schema.MSG_TYPE_REGISTER_RESPONSE,
		Msg: & schema.RegisterResponse {
			ClientId: clientId,
			Status: schema.STATUS_SUCCESS,
		},
	}

	err = schema.WriteMsg(w, msgPack)
}

//msg to client 
func (s *Server) msgHandler(w http.ResponseWriter, req *http.Request) {
	defer func(){
		req.Body.Close()
	}()

	if its, ok := req.URL.Query()["clientid"]; ok && len(its[0]) > 0 {
		clientId := its[0]
		if ci, ok := s.Clients.Load(clientId); ok {
			client, _ := ci.(*Client)
			client.HttpHandler(w, req)
		}
	}
}

func (s *Server) cleanerLoop() {
	for {
		t := time.Now()
		shouldDelete := []string{}
		s.Clients.Range(func (k, v interface{}) bool {
			client, _ := v.(*Client)
			if t.Sub(client.LastHeartbeatTime).Milliseconds() > s.TimeoutSecond.Milliseconds() {
				shouldDelete = append(shouldDelete, client.ClientId)
				client.Stop()
			}
			return true
		})

		for _, clientId := range shouldDelete {
			s.Clients.Delete(clientId)
		}

		time.Sleep(time.Second * 10)
	}
}

func (s *Server) getAllInfo() *schema.AllInfo {
	allInfo := &schema.AllInfo {
		ServerAddr: s.ServerAddr,
		Clients: []*schema.ClientInfo{},
	}

	s.Clients.Range(func (k,v interface{}) bool {
		client := v.(*Client)
		cinfo := & schema.ClientInfo {
			ClientId: client.ClientId,
			ClientAddr: client.ClientAddr,
			Port: client.ToPort,
			Protocol: client.Protocol,
			SourceAddr: client.SourceAddr,
			Direction: client.Direction,
			Description: client.Description,
			Compress: client.Compress,
			HttpVersion: client.HttpVersion,
			ConnectionNumber: common.Max(client.ConnNumber, 0),
			UploadSpeed: client.UploadSpeedMonitor.GetSpeed(),
			DownloadSpeed: client.DownloadSpeedMonitor.GetSpeed(),
		}

		allInfo.Clients = append(allInfo.Clients, cinfo)
		return true
	})
	
	return allInfo
}

func (s *Server) Start() {
	logger.Info(fmt.Sprintf("\nserver start\nAddr: %v\n", s.ServerAddr))

	//start client cleaner
	go s.cleanerLoop()

	http.HandleFunc("/register", s.registerHandler)
	http.HandleFunc("/msg", s.msgHandler)
	http.HandleFunc("/heartbeat", s.heartbeatHandler)
	http.HandleFunc("/monitor", s.monitorHandler)
	http.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("./ui"))))
	http.ListenAndServe(s.ServerAddr, nil)
}
