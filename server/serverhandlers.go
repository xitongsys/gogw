package server

import (
	"net/http"
	"strings"
)

func (server *Server) monitorHandler(w http.ResponseWriter, req *http.Request) {
	if its, ok := req.URL.Query()["key"]; ok && len(its[0]) > 0 {
		key := strings.ToLower(its[0])
		if key == "all" {
		}
	}
}

