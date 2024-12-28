package instances

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

	userId, err := ctx.Values().GetInt("userId")
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

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
	var instances []models.Instances
	totalResult := database.DB.Model(&models.Instances{}).Where("user_id = ?", userId).Count(&total)
	if totalResult.Error != nil {
		l.Error("list instances error: %v", totalResult.Error)
		middleware.Error(ctx, middleware.CodeInstanceListError, iris.StatusInternalServerError)
		return
	}

	result := database.DB.Where("user_id = ?", userId).Limit(limit).Offset(offset).Select("id", "server_id", "cpu_only", "gpu_count", "volume_size", "ssh_address", "ssh_passwd", "jupyter_address", "tensor_board_address", "grafana_address", "status", "image_name", "label", "created_at").Order("id").Find(&instances)
	if result.Error != nil {
		l.Error("list instances error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeInstanceListError, iris.StatusInternalServerError)
		return
	}

	res := make([]instanceStruct, len(instances))
	for i, instance := range instances {
		res[i] = instanceStruct{
			Instances: instance,
		}
		server := models.Servers{
			ID: instance.ServerID,
		}
		result := database.DB.Select("name", "gpu_type", "gpu_num", "cpu_count_per_gpu", "memory_per_gpu").First(&server)
		if result.Error != nil {
			l.Error("query server %d error: %v", instance.ServerID, result.Error)
			continue
		}

		redis.RawDB.Get(ctx, "remain_gpu:server:"+strconv.Itoa(int(instance.ServerID))).Scan(&server.GpuUsed)
		server.GpuUsed = server.GpuNum - server.GpuUsed

		res[i].ServerName = server.Name
		res[i].GpuType = server.GpuType
		res[i].GpuNum = server.GpuNum
		res[i].GpuUsed = server.GpuUsed
		res[i].CpuCountPerGpu = server.CpuCountPerGpu
		res[i].MemoryPerGpu = server.MemoryPerGpu
	}

	middleware.ResultWithTotal(ctx, res, total)
}
