package common

const (
	// 成功
	SUCC = iota
	// 失败
	ERR
	// 禁止
	FDB_REQ
	// 参数异常
	PARAM_FMT
	// 参数为空
	PARAM_EMPTY
)
const (
	SUCC_STR            = "成功"
	ERR_STR             = "失败"
	PARAM_READ_STR      = "参数读取错误"
	PARAM_FMT_STR       = "参数格式化错误"
	FDB_REQ_STR         = "请求方式错误"
	POST_STR            = "请求错误"
	JSON_STR            = "JSON转化错误"
	FILE_NOT_EXISTS_STR = "文件不存在"
	TOKEN_EMPTY_STR     = "认证为空"
)
