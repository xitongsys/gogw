package server

import (
	"net/http"
	"strings"
	"time"
)

func (server *Server) monitorHandler(w http.ResponseWriter, req *http.Request) {
	if its, ok := req.URL.Query()["key"]; ok && len(its[0]) > 0 {
		key := strings.ToLower(its[0])
		if key == "all" {
			if data, err := server.getAllInfo().Marshal(); err == nil {
				w.Write(data)
			}
		}
	}
}

func (server *Server) heartbeatHandler(w http.ResponseWriter, req *http.Request) {
	if cs, ok := req.URL.Query()["clientid"]; ok && len(cs[0]) > 0 {
		clientId := cs[0]
		if value, ok := server.Clients.Load(clientId); ok {
			client, _ := value.(*Client)
			client.LastHeartbeatTime = time.Now()
		}
	}
}
