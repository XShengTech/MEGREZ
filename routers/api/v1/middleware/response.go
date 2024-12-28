package middleware

import "github.com/kataras/iris/v12"

func Success(ctx iris.Context) {
	ctx.JSON(Response{
		Code: CodeSuccess,
		Msg:  codeMsgMap[CodeSuccess],
		Data: nil,
	})
}

func Result(ctx iris.Context, data any) {
	ctx.JSON(Response{
		Code: CodeSuccess,
		Msg:  codeMsgMap[CodeSuccess],
		Data: &Data{
			Result: data,
		},
	})
}

func ResultWithTotal(ctx iris.Context, data any, total int64) {
	ctx.JSON(Response{
		Code: CodeSuccess,
		Msg:  codeMsgMap[CodeSuccess],
		Data: &Data{
			Result: data,
			Total:  total,
		},
	})
}

func Error(ctx iris.Context, code ResCode, statusCode int) {
	ctx.StatusCode(statusCode)
	ctx.JSON(Response{
		Code: code,
		Msg:  codeMsgMap[code],
		Data: nil,
	})
}
