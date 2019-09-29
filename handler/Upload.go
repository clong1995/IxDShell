package handler

import (
	. "IxDShell/common"
	"IxDShell/network"
	"IxDShell/param/upload"
	"IxDShell/service"
	"net/http"
)

//上传一个
func UploadOne(w http.ResponseWriter, r *http.Request) {
	//跨域
	network.Origin(w)
	if r.Method == http.MethodGet {
		network.FbdReq(w)
	} else if r.Method == http.MethodPost {
		//参数
		p := new(upload.One)
		err := p.Format(w, r)
		if err != nil {
			return
		}
		//token
		token := r.Header.Get("Authorization")
		if token == "" {
			network.ErrStr(w, TOKEN_EMPTY_STR)
			return
		}
		//业务
		err = service.UploadOne(p, token)
		if err != nil {
			network.ErrStr(w, err.Error())
			return
		}
		network.Succ(w, "")
	}
}

//重启任务列表
func UploadRestartTask(w http.ResponseWriter, r *http.Request) {
	//跨域
	network.Origin(w)
	if r.Method == http.MethodGet {
		network.FbdReq(w)
	} else if r.Method == http.MethodPost {
		//参数
		//token
		token := r.Header.Get("Authorization")
		if token == "" {
			network.ErrStr(w, TOKEN_EMPTY_STR)
			return
		}
		//业务
		row, err := service.UploadRestartTask(token)
		if err != nil {
			network.ErrStr(w, err.Error())
			return
		}
		network.Succ(w, row)
	}
}
