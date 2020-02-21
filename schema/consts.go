package schema

const (
	DIRECTION_FORWARD = "forward"
 	DIRECTION_REVERSE = "reverse"
)

const (
	PROTOCOL_TCP = "tcp"
	PROTOCOL_UDP = "udp"
)

const (
	STATUS_SUCCESS = "success"
	STATUS_FAILED = "failed"
)

const (
	ROLE_READER = "reader"
	ROLE_WRITER = "writer"
	ROLE_QUERY_CONNID = "query conn id"
)

const (
	MSG_TYPE_REGISTER_REQUEST = "registerrequest"
	MSG_TYPE_REGISTER_RESPONSE = "registerresponse"

	MSG_TYPE_OPEN_CONN_REQUEST = "openconnrequest"
	MSG_TYPE_OPEN_CONN_RESPONSE = "openconnresponse"

	MSG_TYPE_CLOSE_CONN_REQUEST = "closeconnrequest"
	MSG_TYPE_CLOSE_CONN_RESPONSE = "closeconnresponse"

	MSG_TYPE_MSG_COMMON_REQUEST = "msgrequest"
	//MSG_TYPE_MSG_RESPONSE = "msgresponse" msg response is some specific response
)