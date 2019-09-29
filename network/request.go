package network

import (
	. "IxDShell/common"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

//json字符串请求转结构体
func GetReqJson(r *http.Request, i interface{}) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return fmt.Errorf(PARAM_READ_STR)
	} else {
		err := json.Unmarshal(data, i)
		if err != nil {
			log.Println(err)
			return fmt.Errorf(PARAM_FMT_STR)
		}
	}
	return nil
}

//TODO 非空检测，后期支持多层级校验
func CheckEmptyReqParam(itf interface{}) error {
	t := reflect.TypeOf(itf).Elem()
	v := reflect.ValueOf(itf).Elem()
	var f reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		f = t.Field(i)
		//必填项校验
		if f.Tag.Get("required") == "true" {
			//string类型
			if f.Type.Name() == "string" && v.Field(i).String() == "" {
				errStr := fmt.Sprintf("%s 的参数列表中，%s 不得为空", t.Name(), f.Name)
				log.Println(errStr)
				return fmt.Errorf(errStr)
			}
			//TODO 其他类型
		}
		//TODO 其他校验
	}
	return nil
}

//跨域
func Origin(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin,X-Requested-With,Content-Type,Accept,Authorization")
}
