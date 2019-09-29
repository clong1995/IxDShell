package config

import "flag"

type config struct {
	ServerAddr string
	WebAddr    string
	HttpAddr   string
	Platform   string
}

var CONF *config

func init() {
	serverAddr := flag.String("serverAddr", "http://127.0.0.1:50001", "服务端口")
	webAddr := flag.String("webAddr", "http://127.0.0.1:50000", "界面端口")
	httpAddr := flag.String("httpAddr", ":50003", "本地端口")
	platform := flag.String("platform", "", "运行平台mac/windows")
	flag.Parse()
	CONF = new(config)
	CONF.ServerAddr = *serverAddr
	CONF.WebAddr = *webAddr
	CONF.HttpAddr = *httpAddr
	CONF.Platform = *platform
}
