package toml

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type TomlConfig struct {
	Pillx         PillConfig
	Etcd          DBConfig
	Redis0        DBConfig
	Mysql         []DBConfig
	Mongo         []DBConfig
	Elasticsearch DBConfig
}

var GlobalTomlConfig TomlConfig

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBname   string
}

type PillConfig struct {
	GatewayOuterHost string
	GatewayOuterPort int
	GatewayInnerHost string
	GatewayInnerPort int
	WorkerInnerHost  string
	WorkerInnerPort  int
	GatewayName      string
	WorkerName       string
}

func LoadTomlConfig(filename string) (TomlConfig, error) {

	var tc TomlConfig
	tomlData, err1 := ioutil.ReadFile(filename)
	if err1 != nil {
		fmt.Println("Read failed", err1)
		return tc, err1
	}

	if _, err := toml.Decode(string(tomlData), &tc); err != nil {
		fmt.Println("ReadToml failed", err)
		return tc, err
	}

	GlobalTomlConfig = tc
	return tc, nil
}
