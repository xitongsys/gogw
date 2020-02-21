package server

import (
	"fmt"
	"net/http"
	"strconv"

	"gogw/common"
	"gogw/logger"
)

type Server struct {
	ServerAddr    string
	Clients map[string]*Client
}

func NewServer(serverAddr string) *Server {
	return &Server{
		ServerAddr:    serverAddr,
		Clients:       make(map[string]*Client),
	}
}

//client register
func (s *Server) registerHandler(w http.ResponseWriter, req *http.Request) {
	clientId := common.UUID("clientid")

	defer func(){
		if err := recover(); err != nil {
			logger.Error(err)
		}

		delete(s.Clients, clientId)
	}()
	
	clientAddr := req.RemoteAddr
	toPort := req.URL.Query()["toport"][0]
	direction := req.URL.Query()["direction"][0]
	protocol := req.URL.Query()["protocol"][0]
	sourceAddr := req.URL.Query()["sourceaddr"][0]
	description := req.URL.Query()["description"][0]

	toPortI, err := strconv.Atoi(toPort)
	if err != nil {
		logger.Error(err)
		return
	}

	client := NewClient(
		clientId,
		clientAddr,
		toPortI,
		direction,
		protocol,
		sourceAddr,
		description,
		w,
		req,
	)

	s.Clients[clientId] = client
	err = client.Start()
	logger.Error(err)
}

//reverse client send new connection
func (s *Server) reverseNewConnHandler(w http.ResponseWriter, req *http.Request) {
	defer func(){
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	
	clientId := req.URL.Query()["clientid"][0]
	if client, ok := s.Clients[clientId]; ok {
		client.ReverseNewConnHandler(w, req)
	}
}

func (s *Server) Start() {
	logger.Info(fmt.Sprintf("\nserver start\nAddr: %v\n", s.ServerAddr))

	http.HandleFunc("/register", s.registerHandler)
	http.HandleFunc("/reversenewconn", s.reverseNewConnHandler)
	http.HandleFunc("/monitor", s.monitorHandler)
	http.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("./ui"))))
	http.ListenAndServe(s.ServerAddr, nil)
}
