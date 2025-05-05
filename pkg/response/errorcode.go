package response

// 错误码规则:
// 1. 成功返回    successCode = 0
// 2. 错误返回均为负数
// 3. 参考 http code 定义
// 			1. 用户错误  	40XXX
// 			2. 系统机错误 	50XXX

var (
	Success          = 0
	ParamError       = 40001
	PermissionDenied = 40002
	ServerError      = 50000
	TooManyRequests  = 50001
)

var ErrorMessage = map[int]string{
	0:     "success",
	40000: "user error",
	40001: "参数有误",
	40002: "无访问权限",
	50000: "server error",
	50001: "请求过多",
}

func message(code int) (msg string) {
	//
	msg, ok := ErrorMessage[code]
	if ok {
		return
	}

	//
	if code < 50000 {
		msg = ErrorMessage[40000]
	} else {
		msg = ErrorMessage[50000]
	}
	return
}
