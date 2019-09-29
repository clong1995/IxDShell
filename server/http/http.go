package http

import (
	. "IxDShell/config"
	_ "IxDShell/router"
	"log"
	"net/http"
)

func StartHttp(addr string) {
	if addr == "" {
		addr = CONF.HttpAddr
	}
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
