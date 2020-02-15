package server

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"gogw/common"
	"gogw/common/schema"
	"gogw/logger"
)

type Server struct {
	ServerAddr    string
	TimeoutSecond time.Duration

	Lock    sync.Mutex
	Clients map[schema.ClientId]Client
}

func NewServer(serverAddr string, timeoutSecond int) *Server {
	return &Server{
		ServerAddr:    serverAddr,
		Clients:       make(map[schema.ClientId]Client),
		TimeoutSecond: time.Duration(timeoutSecond) * time.Second,
	}
}

func (server *Server) cleaner() {
	server.Lock.Lock()
	defer server.Lock.Unlock()

	t := time.Now()
	shouldDelete := []schema.ClientId{}
	for clientId, client := range server.Clients {
		if t.Sub(client.GetLastHeartbeat()).Milliseconds() > server.TimeoutSecond.Milliseconds() {
			shouldDelete = append(shouldDelete, clientId)
			client.Stop()
		}
	}

	for _, clientId := range shouldDelete {
		delete(server.Clients, clientId)
	}
}

func (server *Server) checkPort(port int) error {
	l, err := net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		return err
	}
	l.Close()
	return nil
}

func (server *Server) registerHandler(w http.ResponseWriter, req *http.Request) {
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error(err)
		return
	}

	registerRequest := &schema.RegisterRequest{}
	if err = registerRequest.Unmarshal(bs); err != nil {
		logger.Error(err)
		return
	}

	if registerRequest.ToPort <= 0 {
		for registerRequest.ToPort = 1000; registerRequest.ToPort < 65535; registerRequest.ToPort++ {
			if server.checkPort(registerRequest.ToPort) == nil {
				break
			}
		}
	}

	if err = server.checkPort(registerRequest.ToPort); err != nil {
		logger.Error(err)
		return
	}

	clientId := schema.ClientId(common.UUID("clientid"))

	registerResponse := &schema.RegisterResponse{
		ClientId: clientId,
		ToPort:   registerRequest.ToPort,
		Code:     schema.SUCCESS,
	}

	client := NewClientTCP(clientId, req.RemoteAddr, registerRequest.ToPort, registerRequest.SourceAddr, registerRequest.Description)

	server.Lock.Lock()
	defer server.Lock.Unlock()
	server.Clients[clientId] = client

	if err = client.Start(); err == nil {
		var data []byte
		if data, err = registerResponse.Marshal(); err == nil {
			_, err = w.Write(data)
		}
	}

	if err != nil {
		delete(server.Clients, clientId)
		registerResponse.Code = schema.FAILED
		data, _ := registerResponse.Marshal()
		w.Write(data)

		logger.Error(err)
	}

}

func (server *Server) packHandler(w http.ResponseWriter, req *http.Request) {
	if cs, ok := req.URL.Query()["clientid"]; ok && len(cs[0]) > 0 {
		clientId := schema.ClientId(cs[0])
		if client, ok := server.Clients[clientId]; ok && client != nil {
			client.RequestHandler(w, req)
		}
	}
}

func (server *Server) Start() {
	logger.Info(fmt.Sprintf("\nserver start\nAddr: %v\nTimeoutSecond: %v\n", server.ServerAddr, int(server.TimeoutSecond.Seconds())))

	http.HandleFunc("/register", server.registerHandler)
	http.HandleFunc("/pack", server.packHandler)
	http.HandleFunc("/heartbeat", server.heartbeatHandler)
	http.HandleFunc("/monitor", server.monitorHandler)
	http.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("./ui"))))

	//cleaner
	go func() {
		for {
			server.cleaner()
			time.Sleep(time.Second * 30)
		}
	}()

	http.ListenAndServe(server.ServerAddr, nil)
}
