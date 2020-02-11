package schema

type ClientId string

type ConnectionId string

type ErrorCode int
const (
	_ ErrorCode = iota
	SUCCESS
	FAILED
)

type PackType int 
const (
	_ PackType = iota
	CLIENT_REQUEST_PACK
	CLIENT_SEND_PACK
	CLIENT_REQUEST_CMD
	CLIENT_SEND_CMD

	SERVER_PACK
	SERVER_CMD

	OPEN
	NORMAL
	CLOSE
)

const (
	CMD_OPEN_CONN = "open_conn"
	CMD_CLOSE_CONN = "close_conn"
)