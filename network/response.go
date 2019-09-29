package network

import (
	. "IxDShell/common"
	"encoding/json"
	"log"
	"net/http"
)

// RespMsg : http响应数据的通用结构
type respMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Succ(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	resp := new(respMsg)
	resp.Code = SUCC
	resp.Msg = SUCC_STR
	resp.Data = data

	r, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(r)
	if err != nil {
		log.Println(err)
	}
}

func ErrStrCode(w http.ResponseWriter, msg string, code int) {
	resp := new(respMsg)
	resp.Code = code
	resp.Msg = msg

	r, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(r)
	if err != nil {
		log.Println(err)
	}
}

func ErrStr(w http.ResponseWriter, msg string) {
	ErrStrCode(w, msg, ERR)
}

func FbdReq(w http.ResponseWriter) {
	ErrStrCode(w, FDB_REQ_STR, FDB_REQ)
}
