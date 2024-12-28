package servers

import (
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/database"
	"megrez/services/redis"
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
	var servers []models.Servers
	totalResult := database.DB.Model(&models.Servers{}).Count(&total)
	if totalResult.Error != nil {
		l.Error("list servers error: %v", totalResult.Error)
		middleware.Error(ctx, middleware.CodeAdminServerListError, iris.StatusInternalServerError)
		return
	}

	result := database.DB.Limit(limit).Offset(offset).Select("id", "name", "ip", "gpu_type", "gpu_num", "gpu_driver_version", "gpu_cuda_version", "cpu_count_per_gpu", "memory_per_gpu", "volume_total", "price", "price_volume", "gpu_used", "volume_used", "created_at").Order("id").Find(&servers)
	if result.Error != nil {
		l.Error("list servers error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeAdminServerListError, iris.StatusInternalServerError)
		return
	}

	for i, server := range servers {
		redis.RawDB.Get(ctx, "remain_gpu:server:"+strconv.Itoa(int(server.ID))).Scan(&servers[i].GpuUsed)
		servers[i].GpuUsed = server.GpuNum - servers[i].GpuUsed

		redis.RawDB.Get(ctx, "remain_volume:server:"+strconv.Itoa(int(server.ID))).Scan(&servers[i].VolumeUsed)
		servers[i].VolumeUsed = server.VolumeTotal - servers[i].VolumeUsed
	}

	middleware.ResultWithTotal(ctx, servers, total)
}
