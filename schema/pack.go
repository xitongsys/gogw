package schema

import (
	"encoding/json"
)

type PackRequest struct {
	ClientId ClientId
	ConnId ConnectionId
	Type PackType
}

func (packRequest *PackRequest) Marshal() ([]byte, error) {
	return json.Marshal(packRequest)
}

func (packRequest *PackRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, packRequest)
}

type PackResponse struct {
	ClientId ClientId
	ConnId ConnectionId
	Type PackType
	PackContent string
}

func (packResponse *PackResponse) Marshal() ([]byte, error) {
	return json.Marshal(packResponse)
}

func (packResponse *PackResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, packResponse)
}