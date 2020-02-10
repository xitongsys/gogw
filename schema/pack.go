package schema

import (
	"encoding/json"
)

type PackRequest struct {
	ClientId ClientId
	ConnId ConnectionId
	Type PackType
}

func (packRequest *PackRequest) marshal() ([]byte, error) {
	return json.Marshal(packRequest)
}

func (packRequest *PackRequest) unmarshal(data []byte) error {
	return json.Unmarshal(data, packRequest)
}

type PackResponse struct {
	ClientId ClientId
	ConnId ConnectionId
	Type PackType
	PackContent string
}

func (packResponse *PackResponse) marshal() ([]byte, error) {
	return json.Marshal(packResponse)
}

func (packResponse *PackResponse) unmarshal(data []byte) error {
	return json.Unmarshal(data, packResponse)
}