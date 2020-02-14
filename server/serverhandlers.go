package server

import (
	"net/http"
	"strings"
	"time"

	"gogw/common/schema"
	"gogw/logger"
)

func (server *Server) heartbeatHandler(w http.ResponseWriter, req *http.Request) {
	if cs, ok := req.URL.Query()["clientid"]; ok && len(cs[0]) > 0 {
		clientId := schema.ClientId(cs[0])
		if client, ok := server.Clients[clientId]; ok {
			client.LastHeartbeat = time.Now()
			w.Write([]byte("gogw ok"))
		}
	}
}

func (server *Server) monitorHandler(w http.ResponseWriter, req *http.Request) {
	if its, ok := req.URL.Query()["key"]; ok && len(its[0]) > 0 {
		key := strings.ToLower(its[0])
		if key == "all" {
			if data, err := server.getAllInfo().Marshal(); err == nil {
				w.Write(data)
			}else{
				logger.Error(err)
			}
		}
	}
}

func (server *Server) getAllInfo() *schema.AllInfo {
	allInfo := &schema.AllInfo {
		ServerAddr: server.ServerAddr,
		Clients: []*schema.ClientInfo{},
	}

	for _, client := range server.Clients {
		cinfo := & schema.ClientInfo {
			ClientId: client.ClientId,
			ClientAddr: client.ClientAddr,
			Port: client.PortTo,
			SourceAddr: client.SourceAddr,
			Description: client.Description,
			ConnectionNumber: len(client.Conns),
			UploadSpeed: client.SpeedMonitor.GetUploadSpeed(),
			DownloadSpeed: client.SpeedMonitor.GetDownloadSpeed(),
		}

		allInfo.Clients = append(allInfo.Clients, cinfo)
	}

	return allInfo
}