package upload

import (
	. "IxDShell/common"
	"IxDShell/network"
	"net/http"
)

//参数
type One struct {
	//本地要上传的文件
	LocalPath string `json:"localPath" required:"true"`
	Pid       string `json:"pid" required:"true"`
}

func (p *One) Format(w http.ResponseWriter, r *http.Request) error {
	//json转结构体
	err := network.GetReqJson(r, p)
	if err != nil {
		network.ErrStrCode(w, err.Error(), PARAM_FMT)
		return err
	}

	//非空检验
	err = network.CheckEmptyReqParam(p)
	if err != nil {
		network.ErrStrCode(w, err.Error(), PARAM_EMPTY)
		return err
	}

	//其他业务校验

	return nil
}
