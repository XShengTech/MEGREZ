package servers

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

	var instances []models.Instances
	result := database.DB.Where("server_id = ?", id).Find(&instances)
	if result.Error != nil {
		l.Error("get instances error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeAdminServerDeleteError, iris.StatusInternalServerError)
		return
	}
	if len(instances) > 0 {
		middleware.Error(ctx, middleware.CodeAdminServerInstanceError, iris.StatusBadRequest)
		return
	}

	server := models.Servers{
		ID: id,
	}
	result = database.DB.Delete(&server)
	if result.Error != nil {
		l.Error("delete server error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeAdminServerDeleteError, iris.StatusInternalServerError)
		return
	}

	middleware.Success(ctx)
}
