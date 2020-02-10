package schema

type RegisterRequest struct {
	Port int
}

type RegisterResponse struct {
	ClientId ClientId
	Code ErrorCode
}

type UnregisterRequest struct {
	ClientId ClientId
}

type UnregisterResponse struct {
	ClientId ClientId
	Code ErrorCode
}
