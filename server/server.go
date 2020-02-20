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
	Clients map[string]*Client
}

func NewServer(serverAddr string, timeoutSecond int) *Server {
	return &Server{
		ServerAddr:    serverAddr,
		Clients:       make(map[schema.ClientId]IClient),
	}
}

func (server *Server) registerHandler(w http.ResponseWriter, req *http.Request) {
	clientId := common.UUID("clientid")

	defer func(){
		if err := recover(); err != nil {
			logger.Error(err)
		}

		delete(server.Clients, clientid)
	}()
	
	clientAddr := req.RemoteAddr
	toPort := req.URL.Query()["toport"][0]
	direction := req.URL.Query()["direction"][0]
	protocol := req.URL.Query()["protocol"][0]
	sourceAddr := req.URL.Query()["sourceaddr"][0]
	description := req.URL.Query()["description"][0]

	client := NewClient(
		clientId,
		clientAddr,
		toPort,
		direction,
		protocol,
		sourceAddr,
		description,
		w,
		r,
	)

	server.Clients[clientId] = client
	err := client.Start()

	logger.Error(err)
}

func (server *Server) reverseNewConnHandler(w http.ResponseWriter, r *http.Request) {
	defer func(){
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	
	clientId := req.URL.Query["clientid"][0]
	if client, ok := server.Clients[clientId]; ok {
		client.ReverseNewConnHandler(w, r)
	}
}

func (server *Server) Start() {
	logger.Info(fmt.Sprintf("\nserver start\nAddr: %v\nTimeoutSecond: %v\n", server.ServerAddr, int(server.TimeoutSecond.Seconds())))

	http.HandleFunc("/register", server.registerHandler)
	http.HandleFunc("/reversenewconn", server.packHandler)
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
