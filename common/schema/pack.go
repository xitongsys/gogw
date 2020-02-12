package schema

import (
	"encoding/json"
	"encoding/base64"
)

type PackRequest struct {
	ClientId ClientId
	ConnId ConnectionId
	Type PackType
	Content string
}

func (packRequest *PackRequest) Marshal() ([]byte, error) {
	packRequest.Content = base64.StdEncoding.EncodeToString([]byte(packRequest.Content))
	return json.Marshal(packRequest)
}

func (packRequest *PackRequest) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, packRequest)
	if err != nil {
		return err
	}
	data, err = base64.StdEncoding.DecodeString(packRequest.Content)
	packRequest.Content = string(data)
	return err
}

type PackResponse struct {
	ClientId ClientId
	ConnId ConnectionId
	Type PackType
	Content string
}

func (packResponse *PackResponse) Marshal() ([]byte, error) {
	packResponse.Content = base64.StdEncoding.EncodeToString([]byte(packResponse.Content))
	return json.Marshal(packResponse)
}

func (packResponse *PackResponse) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, packResponse)
	if err != nil {
		return err
	}
	data, err = base64.StdEncoding.DecodeString(packResponse.Content)
	packResponse.Content = string(data)
	return err
}