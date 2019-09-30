package main

import (
	"IxDShell/server/http"
	"IxDShell/shell"
	"log"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func main() {
	//启动http
	go http.StartHttp("")
	//mac客户端
	//shell.StartMac()
	//windows客户端
	shell.StartWindows()
}
