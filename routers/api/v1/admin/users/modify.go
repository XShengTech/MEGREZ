package users

import (
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/database"

	"github.com/kataras/iris/v12"
)

type modifyReqStruct struct {
	Password *string `json:"password"`
	Role     *int    `json:"role"`
}

func modifyHandler(ctx iris.Context) {
	l.SetFunction("modifyHandler")

	userId, err := ctx.Params().GetUint("id")
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	var req modifyReqStruct
	err = ctx.ReadJSON(&req)
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	if userId == 1 {
		if req.Role != nil {
			middleware.Error(ctx, middleware.CodeAdminUserModifyError, iris.StatusBadRequest)
			return
		}
	}

	user := models.Users{
		ID: userId,
	}
	result := database.DB.Where(&user).First(&user)
	if result.Error != nil {
		middleware.Error(ctx, middleware.CodeServeBusy, iris.StatusInternalServerError)
		return
	}

	if req.Password != nil {
		if *req.Password != "" {
			user.Password = user.PasswordHash(*req.Password)
		}
	}

	if req.Role != nil {
		user.Role = *req.Role
	}

	result = database.DB.Save(&user)
	if result.Error != nil {
		middleware.Error(ctx, middleware.CodeServeBusy, iris.StatusInternalServerError)
		return
	}

	middleware.Success(ctx)
}
