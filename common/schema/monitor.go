package schema

type AllInfo struct {
	ServerAddr string
	Clients []*ClientInfo
}

type ClientInfo struct {
	ClientId ClientId
	Port int
	ConnectionNumber int
	UploadSpeed int
	DownloadSpeed int
}