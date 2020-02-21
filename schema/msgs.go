package schema

type MsgPack struct {
	MsgType string
	MsgContent string
	Msg interface{}
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