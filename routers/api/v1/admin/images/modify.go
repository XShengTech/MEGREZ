package images

import (
	"encoding/json"
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/database"

	"github.com/kataras/iris/v12"
)

func modifyHandler(ctx iris.Context) {
	l.SetFunction("modifyHandler")

	req := make(map[string]string)
	err := ctx.ReadJSON(&req)
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	system := models.System{
		Key: imagesKey,
	}
	result := database.DB.FirstOrCreate(&system)
	if result.Error != nil {
		l.Error("get system config error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeServeBusy, iris.StatusInternalServerError)
		return
	}

	valueBytes, err := json.Marshal(req)
	if err != nil {
		l.Error("marshal system error: %v", err)
		middleware.Error(ctx, middleware.CodeServeBusy, iris.StatusInternalServerError)
		return
	}

	system.Value = string(valueBytes)
	result = database.DB.Save(&system)
	if result.Error != nil {
		l.Error("save system error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeServeBusy, iris.StatusInternalServerError)
		return
	}

	middleware.Success(ctx)
}
