package user

import (
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/database"
	"megrez/services/redis"

	"github.com/kataras/iris/v12"
)

type forgerPasswordStruct struct {
	Code       string `json:"code"`
	Password   string `json:"password"`
	RePassword string `json:"repassword"`
}

func forgetPasswordHandler(ctx iris.Context) {
	l.SetFunction("forgetPasswordHandler")

	var req forgerPasswordStruct
	if err := ctx.ReadJSON(&req); err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	if req.Code == "" || req.Password == "" || req.RePassword == "" {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	if req.Password != req.RePassword {
		middleware.Error(ctx, middleware.CodePasswordNotMatch, iris.StatusBadRequest)
		return
	}

	rdb := redis.RawDB
	v := rdb.Get(ctx, forgetPasswordRedisKeyPrefix+req.Code)

	if v.Err() != nil {
		middleware.Error(ctx, middleware.CodeUserVerifyInvalid, iris.StatusBadRequest)
		return
	}

	r := rdb.Del(ctx, forgetPasswordRedisKeyPrefix+req.Code)
	if r.Err() != nil {
		middleware.Error(ctx, middleware.CodeServeBusy, iris.StatusInternalServerError)
		l.Error("delete redis verify code error: %v", r.Err())
		return
	}

	id, err := v.Int()
	if err != nil {
		middleware.Error(ctx, middleware.CodeUserVerifyInvalid, iris.StatusBadRequest)
		return
	}

	user := models.Users{
		ID: uint(id),
	}
	result := database.DB.First(&user)
	if result.Error != nil {
		l.Error("get user error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeUserVerifyInvalid, iris.StatusBadRequest)
		return
	}

	user.Password = user.PasswordHash(req.Password)

	result = database.DB.Model(&user).Update("password", user.Password)
	if result.Error != nil {
		middleware.Error(ctx, middleware.CodeServeBusy, iris.StatusInternalServerError)
		l.Error("update user password Error: %v", result.Error)
		return
	}

	middleware.Success(ctx)
}
