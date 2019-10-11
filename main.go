package main

import (
	"IxDShell/server/aria2"
	"IxDShell/shell"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func main() {
	//单独协程运行启动aria2
	go func() {
		err := aria2.StartAria2()
		if err != nil {
			os.Exit(2)
		}
	}()

	//mac客户端
	//shell.StartMac()
	//windows客户端
	shell.StartWindows()
}
