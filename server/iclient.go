package server

import (
	"gogw/common/schema"
	"net/http"
	"time"
)

const (
	PACKSIZE = 1024*1024
	BUFFSIZE = 100
)

type IClient interface {
	Start() error 
	Stop() error
	RequestHandler(w http.ResponseWriter, req *http.Request)

	GetClientId() schema.ClientId
	GetClientAddr() string
	GetToPort() int
	GetDirection() string
	GetProtocol() string
	GetSourceAddr() string
	GetDescription() string
	GetConnectionNumber() int
	GetSpeedMonitor() *SpeedMonitor
	GetLastHeartbeat() time.Time

	SetLastHeartbeat(t time.Time) 
}
