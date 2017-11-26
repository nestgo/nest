package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v1"
)

var (
	cfg *Config
)

const (
	//FilePath 配置文件路径
	FilePath = "./config/config.yml"
)

//Config config
type Config struct {
	Address string `yaml:"address"`
}

//GetConfig 获取配置
func GetConfig() *Config {
	if cfg != nil {
		return cfg
	}
	data, err := ioutil.ReadFile(FilePath)
	if err != nil {
		panic(err)
	}
	c := new(Config)
	err = yaml.Unmarshal(data, c)
	if err != nil {
		panic(err)
	}
	cfg = c
	return cfg
}
