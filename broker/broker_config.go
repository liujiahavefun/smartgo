package main

import (
	"encoding/json"
	"github.com/golang/glog"
	"os"
)

type BrokerConfig struct {
	configFile        string
	TransportProtocol string
	ListenOn          string
	LogFile           string
	RouterServerList  []string
	RouterServerNum   int
}

func NewBrokerConfig(configFile string) *BrokerConfig {
	return &BrokerConfig{
		configFile: configFile,
	}
}

func (self *BrokerConfig) LoadConfig() error {
	file, err := os.Open(self.configFile)
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	defer file.Close()

	//liujia: file内部存的是JSON，将这个JSON转为对应的结构，注意字段名字，嵌套关系，和类型，要一致！
	dec := json.NewDecoder(file)
	err = dec.Decode(&self)
	if err != nil {
		return err
	}

	return nil
}