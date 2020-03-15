package schema

type MsgPack struct {
	MsgType string
	MsgContent string
	Msg interface{}
}

type RegisterRequest struct {
	SourceAddr string
	ToPort int
	Direction string
	Protocol string
	Description string
	Compress bool
	HttpVersion string
}

type RegisterResponse struct {
	ClientId string
	Status string
}

type OpenConnRequest struct {
	ConnId string

	//ROLE_READER/ROLE_WRITER/ROLE_QUERY_CONNID
	Role string
}

type OpenConnResponse struct {
	ConnId string

	//STATUS_SUCCESS/STATUS_FAILED
	Status string
}