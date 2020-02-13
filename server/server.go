package server

import (
	"fmt"
	"net"
	"net/http"

	"gogw/common"
	"gogw/common/schema"
	"gogw/logger"
)

type Server struct {
	ServerAddr string

	Clients map[schema.ClientId]*Client
}

func NewServer(serverAddr string) *Server {
	return &Server{
		ServerAddr: serverAddr,
		Clients:    make(map[schema.ClientId]*Client),
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
	if ps, ok := req.URL.Query()["port"]; ok && len(ps[0]) > 0 {
		clientId := schema.ClientId(common.UUID())
		var port int
		fmt.Sscanf(ps[0], "%d", &port)

		err := server.checkPort(port)

		registerResponse := &schema.RegisterResponse{
			ClientId: clientId,
			Code:     schema.SUCCESS,
		}

		if err == nil {
			client := NewClient(clientId, req.RemoteAddr, port)
			server.Clients[clientId] = client
	
			if err = client.Start(); err == nil {
				var data []byte
				if data, err = registerResponse.Marshal(); err == nil {
					_, err = w.Write(data)
				}
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
}

func (server *Server) packHandler(w http.ResponseWriter, req *http.Request) {
	if cs, ok := req.URL.Query()["clientid"]; ok && len(cs[0]) > 0 {
		clientId := schema.ClientId(cs[0])
		if client, ok := server.Clients[clientId]; ok && client != nil {
			client.requestHandler(w, req)
		}
	}
}

func (server *Server) Start() {
	logger.Info("server start:", server.ServerAddr)

	http.HandleFunc("/register", server.registerHandler)
	http.HandleFunc("/pack", server.packHandler)
	http.HandleFunc("/test", server.testHandler)
	http.HandleFunc("/monitor", server.monitorHandler)
	http.ListenAndServe(server.ServerAddr, nil)
}
