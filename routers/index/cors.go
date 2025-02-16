package index

import "github.com/kataras/iris/v12"

func cors(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "GET")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	ctx.Header("Access-Control-Max-Age", "86400")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Next()
}
