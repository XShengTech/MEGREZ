package user

import (
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/database"

	"github.com/kataras/iris/v12"
)

type resetPasswordStruct struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	RePassword  string `json:"re_password"`
}

func resetPasswordHandler(ctx iris.Context) {
	l.SetFunction("resetPasswordHandler")

	userId, err := ctx.Values().GetInt("userId")
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	var req resetPasswordStruct
	if err := ctx.ReadJSON(&req); err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	if req.NewPassword != req.RePassword {
		middleware.Error(ctx, middleware.CodePasswordNotMatch, iris.StatusBadRequest)
		return
	}

	user := models.Users{
		ID: uint(userId),
	}
	result := database.DB.First(&user)
	if result.Error != nil {
		l.Error("get user error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeUserNotExist, iris.StatusInternalServerError)
		return
	}

	if !user.CheckPassword(req.OldPassword) {
		middleware.Error(ctx, middleware.CodePasswordError, iris.StatusBadRequest)
		return
	}

	user.Password = user.PasswordHash(req.NewPassword)
	result = database.DB.Save(&user)
	if result.Error != nil {
		l.Error("save user error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeInternalPatchError, iris.StatusInternalServerError)
		return
	}

	middleware.Success(ctx)
}
