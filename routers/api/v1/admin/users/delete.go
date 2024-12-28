package users

import (
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/database"

	"github.com/kataras/iris/v12"
)

func deleteHandler(ctx iris.Context) {
	l.SetFunction("deleteHandler")

	id, err := ctx.Params().GetUint("id")
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	if id == 1 {
		middleware.Error(ctx, middleware.CodeAdminUserDeleteError, iris.StatusBadRequest)
		return
	}

	var instances []models.Instances
	result := database.DB.Where("user_id = ?", id).Find(&instances)
	if result.Error != nil {
		l.Error("get user instances error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeAdminUserDeleteError, iris.StatusInternalServerError)
		return
	}

	if len(instances) > 0 {
		middleware.Error(ctx, middleware.CodeAdminUserInstanceNoEmpty, iris.StatusBadRequest)
		return
	}

	user := models.Users{
		ID: id,
	}
	result = database.DB.Delete(&user)
	if result.Error != nil {
		l.Error("delete user error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeAdminUserDeleteError, iris.StatusInternalServerError)
		return
	}

	middleware.Success(ctx)
}
