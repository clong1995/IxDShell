package config

import "flag"

type config struct {
	ServerAddr string
	QiniuAddr  string
	WebAddr    string
	HttpAddr   string
	Platform   string
}

var CONF *config

func init() {
	//addr := "127.0.0.1"
	addr := "quickex.com.cn"
	serverAddr := flag.String("serverAddr", "http://"+addr+":50001", "服务端口")
	qiniuAddr := flag.String("qiniuAddr", "http://storage.quickex.com.cn/", "七牛服务")
	webAddr := flag.String("webAddr", "http://"+addr+":50000", "界面端口")
	httpAddr := flag.String("httpAddr", ":50003", "本地端口")
	platform := flag.String("platform", "", "运行平台mac/windows")
	flag.Parse()
	CONF = new(config)
	CONF.ServerAddr = *serverAddr
	CONF.QiniuAddr = *qiniuAddr
	CONF.WebAddr = *webAddr
	CONF.HttpAddr = *httpAddr
	CONF.Platform = *platform
}
