package schema

import (
	"github.com/vmihailenco/msgpack/v4"
)

type PackRequest struct {
	ClientId ClientId
	ConnId ConnectionId
	Type PackType
	Content string
}

func (packRequest *PackRequest) Marshal() ([]byte, error) {
	return msgpack.Marshal(packRequest)
}

func (packRequest *PackRequest) Unmarshal(data []byte) error {
	return msgpack.Unmarshal(data, packRequest)
}

type PackResponse struct {
	ClientId ClientId
	ConnId ConnectionId
	Type PackType
	Content string
}

func (packResponse *PackResponse) Marshal() ([]byte, error) {
	return msgpack.Marshal(packResponse)
}

func (packResponse *PackResponse) Unmarshal(data []byte) error {
	return msgpack.Unmarshal(data, packResponse)
}