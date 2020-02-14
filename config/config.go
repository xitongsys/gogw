package config

import (
	"io/ioutil"
	"encoding/json"
)

type ClientConfig struct {
	ServerAddr string
	LocalAddr  string
	RemotePort int
	Description string
}

type ServerConfig struct {
	ServerAddr string
}

type Config struct {
	Role string
	Client ClientConfig
	Server ServerConfig
}

func NewConfig(file string) (*Config, error) {
	cfg := &Config{}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return cfg, err
	}

	err = cfg.Unmarshal(data)
	return cfg, err
}

func (cfg *Config) Marshal() ([]byte, error) {
	return json.Marshal(cfg)
}

func (cfg *Config) Unmarshal(data []byte) error {
	return json.Unmarshal(data, cfg)
}
