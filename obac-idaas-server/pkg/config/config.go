package config

import (
	"encoding/json"
	"entrust.com/iat/golang-server/pkg/auth/oidc"
	"entrust.com/iat/golang-server/pkg/config"
	"io/ioutil"
	"os"
)

const (
	configFile          = "rest_server_config.json"
	defaultConfigPath   = "config"
	defaultResourcePath = "resources"
)

type MainConfig struct {
	ServerConfig config.ServerConfig `json:"server"`
	AuthConfig   oidc.AuthConfig     `json:"authentication"`
	App          UriConfig           `json:"app"`
	Idp          UriConfig           `json:"idp"`
	BcProvider   ProviderConfig      `json:"bcProviders"`
}

type UriConfig struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Path     string `json:"path"`
}

type ProviderConfig struct {
	GoerliApiUrl          string `json:"goerliApiUrl"`
	MumbaiApiUrl          string `json:"mumbaiApiUrl"`
	BinanceApiUrl         string `json:"binanceApiUrl"`
	GoerliApiWebSocketUrl string `json:"goerliApiWebSocketUrl"`
}

func LoadConfig(configPath string) (MainConfig, error) {
	myConfig := new(MainConfig)
	file, err := os.Open(configPath + "/" + configFile)
	if err != nil {
		return MainConfig{}, err
	}

	input, err := ioutil.ReadAll(file)
	if err != nil {
		return MainConfig{}, err
	}

	err = json.Unmarshal(input, &myConfig)
	if err != nil {
		return MainConfig{}, err
	}
	return *myConfig, nil
}
