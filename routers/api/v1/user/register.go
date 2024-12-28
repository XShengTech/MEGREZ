package user

import (
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/config"
	"megrez/services/database"

	"github.com/kataras/iris/v12"
)

type registerStruct struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func registerHandler(ctx iris.Context) {
	l.SetFunction("registerHandler")

	var userReq registerStruct
	if err := ctx.ReadJSON(&userReq); err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	if userReq.Username == "" || userReq.Email == "" || userReq.Password == "" {
		middleware.Error(ctx, middleware.CodeRegisterRequestError, iris.StatusBadRequest)
		return
	}

	user := models.Users{
		Username: userReq.Username,
		Email:    userReq.Email,
		Role:     0,
		Balance:  0,
	}
	user.Password = user.PasswordHash(userReq.Password)

	if !config.GetSystemVerify() {
		user.Role = 1
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		l.Error("create user error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeRegisterError, iris.StatusInternalServerError)
		return
	}

	middleware.Success(ctx)
}
