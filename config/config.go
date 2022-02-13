package config

import (
	"encoding/json"
	"io/ioutil"
	"react-demo-server/util"
)

type CONFIG struct {
	Port           string   `json:"port"`
	Env            string   `json:"env"`
	TrustedProxies []string `json:"trusted_proxies"`
	DSN            string   `json:"dsn"`
}

var Config CONFIG

func Start() {

	hasConfigFile := util.FileExists("config/config.json")
	if !hasConfigFile {
		panic("config/config.json is missing")
	}

	fileContent, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(fileContent, &Config)
	if err != nil {
		panic(err)
	}
}
