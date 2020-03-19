package config

import (
	"io/ioutil"
	"encoding/json"
)

type ClientConfig struct {
	ServerAddr string
	SourceAddr  string
	ToPort int
	Protocol string
	Direction string
	Description string
	Compress bool
	HttpVersion string
}

type ServerConfig struct {
	ServerAddr string
	TimeoutSecond int
}

type Config struct {
	Clients []ClientConfig
	Server ServerConfig
}

func NewConfig(data []byte) (*Config, error) {
	cfg := &Config{}
	err := cfg.Unmarshal(data)
	
	return cfg, err
}

func NewConfigFromFile(file string) (*Config, error){
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return NewConfig(data)
}

func (cfg *Config) Marshal() ([]byte, error) {
	return json.Marshal(cfg)
}

func (cfg *Config) Unmarshal(data []byte) error {
	return json.Unmarshal(data, cfg)
}
