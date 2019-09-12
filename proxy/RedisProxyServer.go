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

package main

import (
	"flag"
	"fmt"
	"github.com/pingcap/tidb/store/tikv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"strconv"
	"strings"
	"tedis/proxy/config"
	"tedis/proxy/handler"
	"tedis/proxy/log"
	"tedis/proxy/prometheus"
	"tedis/proxy/redis"
)

var (
	Version        = ""
	GitHash        = ""
	BuildTime      = ""
	versionFlag    = flag.Bool("v", false, "Version")
	serverPort     = flag.Int("port", 6379, "listen port,default: 6379")
	confPath       = flag.String("conf", "", "config file of proxy")
	pdAddr         = flag.String("pd", "127.0.0.1:2379", "pd address,default:localhost:2379")
	logPath        = flag.String("lp", "", "log file path, if empty, default:stdout")
	logLevel       = flag.String("ll", "", "log level:DEBUG|WARN|INFO|ERROR default:DEBUG")
	logMaxKeep     = flag.Uint("kpdays", 7, "keep log days for proxy")
	ignoreTTL      = flag.Bool("it", false, "ignore ttl when read,default false")
	connTimeout    = flag.Int("ct", 60*3, "connect timeout(s),default:180s")
	TimeOutData    = flag.Int("td", 1500, "request time out data")
	PrometheusPort = flag.Int("pp", 8080, "prometheus port,default port 8080")
)

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Println("version: ", Version)
		fmt.Println("build time: ", BuildTime)
		fmt.Println("git version: ", GitHash)
		os.Exit(0)
	}

	proxyConfig := &config.ProxyConfig{}
	config.ParseConf(*confPath, proxyConfig)

	initLog()
	prometheus.TimeOutThresholds = *TimeOutData
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Info("prometheus listen port:", *PrometheusPort)
		http.ListenAndServe(":"+strconv.Itoa(*PrometheusPort), nil)
	}()

	log.Info("serverPort:", *serverPort)
	log.Info("pdAddr:", *pdAddr)
	log.Info("logPath:", *logPath)
	log.Info("logLevel:", *logLevel)
	log.Info("ignoreTTL:", *ignoreTTL)
	log.Info("ConnTimeout:", *connTimeout)

	driver := tikv.Driver{}
	store, err := driver.Open(fmt.Sprintf("tikv://%s?disableGC=true", *pdAddr))
	if err != nil {
		log.Error("creates an TiKV storage with given pdAddr", proxyConfig.PdAddr, err)
		panic(err)
	}
	kvHandler := &handler.TxTikvHandler{Store: store, NameSpace: []byte{0x00}, IgnoreTTL: *ignoreTTL}

	srv, err := redis.NewServer(redis.DefaultConfig().Port(*serverPort).Handler(kvHandler), *connTimeout)
	if err != nil {
		log.Error("creates a redis server error", err)
		panic(err)
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("redis listen port:", proxyConfig.Port, " error:", err)
		panic(err)
	}
}

func initLog() {

	switch strings.ToUpper(*logLevel) {
	case "INFO":
		log.SetLevel(log.LOG_LEVEL_INFO)
	case "ERROR":
		log.SetLevel(log.LOG_LEVEL_ERROR)
	case "WARN":
		log.SetLevel(log.LOG_LEVEL_WARN)
	default:
		log.SetLevel(log.LOG_LEVEL_DEBUG)
	}

	if len(*logPath) > 0 {
		log.SetHighlighting(false)
		log.Debugf("log path: %s, log keep: %u", *logPath, *logMaxKeep)
		err := log.SetOutputByName(*logPath)
		if err != nil {
			log.Fatalf("set log name failed - %s", err)
		}
	}

	log.SetRotateByDay()
	log.SetKeepAge(*logMaxKeep * 24)
	log.RotateDel()
}
