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
	OPEN
	NORMAL
	CLOSE
)