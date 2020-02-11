package server

import (
	"fmt"
	"net/http"

	"gogw/common/schema"
	"gogw/common"
	"gogw/logger"
)

type Server struct {
	ServerAddr string

	Clients map[schema.ClientId]*Client
}

func NewServer(serverAddr string) *Server {
	return & Server{
		ServerAddr: serverAddr,
	}
}

func (server *Server) registerHandler(w http.ResponseWriter, req *http.Request) {
	if ps, ok := req.URL.Query()["port"]; ok && len(ps[0])>0{
		clientId := schema.ClientId(common.UUID())
		var port int
		fmt.Sscanf(ps[0], "%d", &port)

		client := NewClient(clientId, port)
		server.Clients[clientId] = client

		if err := client.Start(); err != nil {
			logger.Error(err)
			delete(server.Clients, clientId)
		}
	}
}

func (server *Server) packHandler(w http.ResponseWriter, req *http.Request) {
	if cs, ok := req.URL.Query()["clientid"]; ok && len(cs[0])>0 {
		clientId := schema.ClientId(cs[0])
		if client, ok := server.Clients[clientId]; ok && client != nil{
			client.requestHandler(w, req)
		}
	}
}

func (server *Server) testHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("gogw"))
}

func (server *Server) Start() {
	logger.Info("server start")

	http.HandleFunc("/register", server.registerHandler)
	http.HandleFunc("/pack", server.packHandler)
	http.HandleFunc("/test", server.testHandler)
	http.ListenAndServe(server.ServerAddr, nil)
}
