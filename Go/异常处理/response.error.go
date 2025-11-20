package i18n

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

const HeaderAioErr = "aio-response-error"

type HTTPResponse struct {
	StatusCode string `json:"code"`
	Message    string `json:"message"`
	Level      string `json:"level"`
	DisplayUS  string `json:"display_us"`
	DisplayCN  string `json:"display_cn"`
}

// NewHttpResponse error 有可能被层层包裹，很多方法都会将error包裹在另一个类型的error里，
// 比如： Wrap() WithMessage() WithCode()等等
// 因此如果遇到了非Error的类型，可以先检查该方法是否实现了 Cause()error 方法，
// 如果实现了该方法，则一层层深入，直到找到Error类型
// 需要注意的是 Cause()error 接口方法并非笔者定义，而是官方errors包中定义的，用于获取根因，
// 因此大多数error类型都会实现该方法
func NewHttpResponse(ctx echo.Context, err error, i18n Ii18nOpr) HTTPResponse {
	for {
		switch e := err.(type) {
		case Error:
			code := e.Code()
			args := e.Args()
			level := e.Level()
			argsString := make([]string, len(args))
			for i, a := range args {
				tmp, ok := a.(string)
				if !ok {
					tmp = fmt.Sprintf("%v", a)
				}
				argsString[i] = tmp
			}
			ctx.Response().Header().Set(HeaderAioErr, err.Error())
			return HTTPResponse{
				StatusCode: code,
				Level:      ErrLevelToString[level],
				DisplayUS:  i18n.Translate("en-us", code, argsString...),
				DisplayCN:  i18n.Translate("zh-cn", code, argsString...),
			}
		case interface{ Cause() error }:
			err = e.Cause()
			continue
		default:
			ctx.Response().Header().Set(HeaderAioErr, e.Error())
			return HTTPResponse{
				StatusCode: "0",
				Level:      ErrLevelToString[ERROR],
				DisplayUS:  fmt.Sprintf("unknown error, the reason: %s", e.Error()),
				DisplayCN:  fmt.Sprintf("未知错误，原因：%s", e.Error()),
			}
		}
	}
}
