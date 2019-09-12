// Copyright (c) 2019 ELEME, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"tedis/proxy/log"
)

var proxyConfig ProxyConfig

type ProxyConfig struct {
	Port       int    `yaml:"port"`
	TurnOnAuth bool   `yaml:turnOnAuth`
	Password   string `yaml:"password"`
	PdAddr     string `yaml:"pdAddr"`
	LogPath    string `yaml:"logPath"`
	LogLevel   string `yaml:"logLevel"`
	SsAddr     string `yaml:"statsAddr"`
	SsName     string `yaml:"stats"`
}

func ParseConf(path string, config *ProxyConfig) {
	if path == "" {
		log.Info("config path is blank use default proxy config")
		return
	}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("read config %s file error: %s", path, err)
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.Errorf("parse config %s file error: %s", path, err)
	}
}

func GetProxyConfig() ProxyConfig {
	return proxyConfig
}
