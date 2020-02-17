package client

import (
	"gogw/common/schema"
	"gogw/logger"
)

type IClient interface {
	Start()
}

func NewClient(
	serverAddr string, 
	sourceAddr string, 
	toPort int, 
	direction string, 
	protocol string, 
	description string, 
	timeoutSecond int) IClient {

	if direction == schema.DIRECTION_FORWARD {
		return NewClientForward(serverAddr, sourceAddr, toPort, protocol, description, timeoutSecond)
		
	}else if direction == schema.DIRECTION_REVERSE {
		return NewClientReverse(serverAddr, sourceAddr, toPort, protocol, description, timeoutSecond)
	}

	logger.Error("unsupported direction: ", direction)
	return nil
}