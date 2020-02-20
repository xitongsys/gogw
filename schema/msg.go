package schema

const (
	MSG_OPEN_CONN
	MSG_CLOSE_CONN
)

type Msg struct {
	MsgType string
}