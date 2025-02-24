package user

import (
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/database"

	"github.com/kataras/iris/v12"
)

type profile struct {
	ID       uint    `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Role     int     `json:"role"`
	Balance  float64 `json:"balance"`
	Verify   bool    `json:"verify"`
}

func profileHandler(ctx iris.Context) {
	l.SetFunction("profileHandler")
	userId, err := ctx.Values().GetInt("userId")
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
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

	profile := profile{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Balance:  user.Balance,
		Verify:   user.Verify,
	}

	middleware.Result(ctx, profile)
}
