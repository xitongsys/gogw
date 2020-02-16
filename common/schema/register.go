package schema

import (
	"encoding/json"
)

type RegisterRequest struct {
	SourceAddr string
	ToPort int
	Direction string
	Protocol string
	Description string
}

func (registerRequest *RegisterRequest) Marshal() ([]byte, error){
	return json.Marshal(registerRequest)
}

func (registerRequest *RegisterRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, registerRequest)
}

type RegisterResponse struct {
	ClientId ClientId
	ToPort int
	Code ErrorCode
}

func (registerResponse *RegisterResponse) Marshal() ([]byte, error) {
	return json.Marshal(registerResponse)
}

func (registerResponse *RegisterResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, registerResponse)
}

type UnregisterRequest struct {
	ClientId ClientId
}

type UnregisterResponse struct {
	ClientId ClientId
	Code ErrorCode
}
