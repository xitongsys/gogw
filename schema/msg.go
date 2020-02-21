package schema

const (
	_ = iota
	MSG_OPEN_CONN
	MSG_CLOSE_CONN
	MSG_SET_CLIENT_ID
)

type Msg struct {
	MsgType int
	MsgContent string
}