package user

import (
	"megrez/libs/utils"
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/database"

	"github.com/kataras/iris/v12"
)

type resetEmailStruct struct {
	Email string `json:"email"`
}

func resetEmailHandler(ctx iris.Context) {
	l.SetFunction("resetEmailHandler")

	userId, err := ctx.Values().GetInt("userId")
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	var req resetEmailStruct
	if err := ctx.ReadJSON(&req); err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	if req.Email == "" {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	if !utils.EmailFormat(req.Email) {
		middleware.Error(ctx, middleware.CodeEmailFormatError, iris.StatusBadRequest)
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

	result = database.DB.Model(&user).Update("email", req.Email).Update("verify", false)
	if result.Error != nil {
		l.Error("save user error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeInternalPatchError, iris.StatusInternalServerError)
		return
	}

	middleware.Success(ctx)
}
