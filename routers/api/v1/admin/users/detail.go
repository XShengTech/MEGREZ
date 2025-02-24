package users

import (
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/database"

	"github.com/kataras/iris/v12"
)

func detailHandler(ctx iris.Context) {
	l.SetFunction("detailHandler")

	id, err := ctx.Params().GetUint("id")
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	user := models.Users{
		ID: id,
	}
	result := database.DB.Select("id", "username", "email", "role", "verify", "balance", "created_at").First(&user)
	if result.Error != nil {
		l.Error("detail user error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeAdminUserDetailError, iris.StatusInternalServerError)
		return
	}

	middleware.Result(ctx, user)
}
