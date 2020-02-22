package common

import "net"

type Conn struct {
	ConnId string
	Conn net.Conn
}
