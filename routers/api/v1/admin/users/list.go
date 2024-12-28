package users

import (
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/database"
	"strconv"

	"github.com/kataras/iris/v12"
)

func listHandler(ctx iris.Context) {
	l.SetFunction("listHandler")
	var err error

	offset := 0
	limit := 20
	offsetStr := ctx.URLParam("offset")
	limitStr := ctx.URLParam("limit")

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
			return
		}
	}
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
			return
		}
	}

	var total int64
	var users []models.Users
	totalResult := database.DB.Model(&models.Users{}).Count(&total)
	if totalResult.Error != nil {
		l.Error("list users error: %v", totalResult.Error)
		middleware.Error(ctx, middleware.CodeAdminUserListError, iris.StatusInternalServerError)
		return
	}

	result := database.DB.Limit(limit).Offset(offset).Select("id", "username", "email", "role", "balance", "created_at").Order("id").Find(&users)
	if result.Error != nil {
		l.Error("list users error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeAdminUserListError, iris.StatusInternalServerError)
		return
	}

	middleware.ResultWithTotal(ctx, users, total)
}
