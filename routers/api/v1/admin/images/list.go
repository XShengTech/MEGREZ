package images

import (
	"encoding/json"
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/database"

	"github.com/kataras/iris/v12"
)

func listHandler(ctx iris.Context) {
	l.SetFunction("listHandler")

	res := make(map[string]string)

	system := models.System{
		Key: imagesKey,
	}
	result := database.DB.FirstOrCreate(&system)
	if result.Error != nil {
		l.Error("get system error: %v", result.Error)
		middleware.Result(ctx, res)
		return
	}

	err := json.Unmarshal([]byte(system.Value), &res)
	if err != nil {
		l.Error("unmarshal system error: %v", err)
		middleware.Error(ctx, middleware.CodeServeBusy, iris.StatusInternalServerError)
		return
	}

	middleware.Result(ctx, res)
}
