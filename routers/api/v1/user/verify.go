package user

import (
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/database"
	"megrez/services/redis"

	"github.com/kataras/iris/v12"
)

func verifyHandler(ctx iris.Context) {
	l.SetFunction("verifyHandler")

	code := ctx.Params().GetString("code")

	rdb := redis.RawDB
	v := rdb.Get(ctx, verifyRedisKeyPrefix+code)

	if v.Err() != nil {
		middleware.Error(ctx, middleware.CodeUserVerifyInvalid, iris.StatusBadRequest)
		return
	}

	r := rdb.Del(ctx, verifyRedisKeyPrefix+code)
	if r.Err() != nil {
		middleware.Error(ctx, middleware.CodeServeBusy, iris.StatusInternalServerError)
		l.Error("delete redis verify code error: %v", r.Err())
		return
	}

	email := v.Val()
	l.Debug("verify email: %s", email)
	user := models.Users{
		Email: email,
	}
	result := database.DB.Where(&user).First(&user)
	if result.Error != nil {
		l.Error("get user error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeUserVerifyInvalid, iris.StatusBadRequest)
		return
	}

	result = database.DB.Model(&user).Update("verify", true)
	if result.Error != nil {
		middleware.Error(ctx, middleware.CodeServeBusy, iris.StatusInternalServerError)
		l.Error("update user verify status Error: %v", result.Error)
		return
	}

	middleware.Success(ctx)
}
