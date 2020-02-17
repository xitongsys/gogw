package server 

import (
	"time"
	"net/http"
	"io/ioutil"
	"sync"
	

	"gogw/common/schema"
	"gogw/logger"
)

type Client struct {
	Lock sync.Mutex

	ClientId schema.ClientId
	ClientAddr string
	ToPort int
	Direction string
	Protocol string
	SourceAddr string
	Description string

	FromClientChanns map[schema.ConnectionId]chan *schema.PackRequest
	ToClientChanns map[schema.ConnectionId]chan *schema.PackResponse
	CmdToClientChann chan *schema.PackResponse

	SpeedMonitor *SpeedMonitor
	LastHeartbeat time.Time

	CmdHandler func(packRequest *schema.PackRequest) *schema.PackResponse
}


func (client *Client) GetClientId() schema.ClientId {
	return client.ClientId
}

func (client *Client) GetClientAddr() string {
	return client.ClientAddr
}

func (client *Client) GetToPort() int {
	return client.ToPort
}

func (client *Client) GetDirection() string {
	return client.Direction
}

func (client *Client) GetProtocol() string {
	return client.Protocol
}

func (client *Client) GetSourceAddr() string {
	return client.SourceAddr
}

func (client *Client) GetDescription() string {
	return client.Description
}

func (client *Client) GetConnectionNumber() int {
	return len(client.FromClientChanns)
}

func (client *Client) GetSpeedMonitor() *SpeedMonitor {
	return client.SpeedMonitor
}

func (client *Client) GetLastHeartbeat() time.Time {
	return client.LastHeartbeat
}

func (client *Client) SetLastHeartbeat(t time.Time) {
	client.LastHeartbeat = t
}

func (client *Client) RequestHandler(w http.ResponseWriter, req *http.Request) {
	defer func(){
		if err := recover(); err != nil {
			logger.Warn(err)
		}
	}()

	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Debug("from client ", string(bs))
	client.SpeedMonitor.Add(-1, int64(len(bs)))

	packRequest := &schema.PackRequest{}
	if err = packRequest.Unmarshal(bs); err != nil {
		logger.Error(err)
		return
	}

	if packRequest.Type == schema.CLIENT_SEND_PACK {
		client.FromClientChanns[packRequest.ConnId] <- packRequest

	}else if packRequest.Type == schema.CLIENT_REQUEST_PACK {
		if packResponse, ok := <- client.ToClientChanns[packRequest.ConnId]; ok {
			data, _ := packResponse.Marshal()

			//logger.Debug("to client", string(packResponse))
			client.SpeedMonitor.Add(int64(len(data)), -1)

			w.Write(data)
		}

	}else if packRequest.Type == schema.CLIENT_SEND_CMD {
		packResponse := client.CmdHandler(packRequest)
		data, _ := packResponse.Marshal()

		//logger.Debug("to client", string(packResponse))
		client.SpeedMonitor.Add(int64(len(data)), -1)

		w.Write(data)
		
	}else if packRequest.Type == schema.CLIENT_REQUEST_CMD {
		select {
		case packResponse := <- client.CmdToClientChann:
			if data, err := packResponse.Marshal(); err == nil {

				//logger.Debug("to client", string(packResponse))
				client.SpeedMonitor.Add(int64(len(data)), -1)

				w.Write(data)
			}
		}
	}
}