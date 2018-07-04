package common

import (
	"errors"
)

type LoginReq struct {
	User string `json:"user"`
	Pwd  string `json:"pwd"`
}

// token 解压里面的内容
type TokenContent struct {
	User string `json:"user"`
	Role string `json:"role"`
}

type LoginToken struct {
	Token string `json:"token"`
}

const (
	Success = "success"

	OperateSuccessCode             = 0 //请求成功
	OperateInvalidParamCode        = 1 //违法的参数
	OperateTokenInvalid            = 2 //token 不对
	OperateTokenExpire             = 3 //token 过期
	OperateRPCConnUnSuccessCode    = 4
	OperateFileIsTooBig            = 5
	OperateOpenFileFailed          = 6
	OperateParseJpgFileFailed      = 7
	OperateParsePngFileFailed      = 8
	OperateFileTypeNoSupportFailed = 9
	OperateFileSaveFailed          = 10
	OperateRPCConnUnSuccessDes     = "RPC service 连接失败"
	OperateSuccessDesc             = "success"
	OperateInvalidParamDesc        = "invalid params"
	CreateFailedMsg                = "failed"
	OperateFailedMsg               = "请求失效..."
	HttpRequestSalt                = "freelancer!@#"
	HttpMobileMsg                  = "手机号码格式不对"

	RpcResponseSucceseCode      = 0 //0表示成功 非零查看具体错误码
	RpcResponseInvalidParamCode = 1 //违法的参数
	RPCResponseFailedCode       = 2 //请求RPC内部服务错误
	RpcResponseInvalidParamDesc = "invalid params"
	RpcResponseMsgSucc          = "request success"
	RpcResponseMsgFailed        = "request failed"

	EmptyValue  = ""
	ExprireTime = 60
)

var (
	TokenExprire = errors.New("token exprite")

	TokenNotExist = errors.New("token not exist")

	// 鉴权开关默认关闭
	AuthSwitch = true

	// 签名开关, 默认关闭
	SigSwitch = true

	// 路由过滤
	RouterFilterMap = map[string]string{
		"/leisure/gateway/v1/login":  "1",
		"/leisure/gateway/v1/health": "1",
	}
)
