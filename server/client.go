package server

import (
	"gogw/common/schema"
	"net/http"
	"time"
)

type Client interface {
	Start() error 
	Stop() error
	RequestHandler(w http.ResponseWriter, req *http.Request)

	GetClientId() schema.ClientId
	GetClientAddr() string
	GetPortTo() int
	GetProtocol() string
	GetSourceAddr() string
	GetDescription() string
	GetConnectionNumber() int
	GetSpeedMonitor() *SpeedMonitor
	GetLastHeartbeat() time.Time

	SetLastHeartbeat(t time.Time) 
}