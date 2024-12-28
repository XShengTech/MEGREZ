package instances

import (
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/database"

	"github.com/kataras/iris/v12"
)

type labelReqStruct struct {
	Label string `json:"label"`
}

func labelHandler(ctx iris.Context) {
	l.SetFunction("labelHandler")

	userId, err := ctx.Values().GetInt("userId")
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	id, err := ctx.Params().GetUint("id")
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	var req labelReqStruct
	err = ctx.ReadJSON(&req)
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	instance := models.Instances{
		ID: id,
	}
	result := database.DB.Where(&instance).Where("user_id = ?", userId).First(&instance)
	if result.Error != nil {
		l.Error("query instance error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeInstanceQueryError, iris.StatusInternalServerError)
		return
	}

	result = database.DB.Model(&instance).Update("label", req.Label)
	if result.Error != nil {
		l.Error("save instance error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeInstanceSaveError, iris.StatusInternalServerError)
		return
	}

	middleware.Success(ctx)
}
