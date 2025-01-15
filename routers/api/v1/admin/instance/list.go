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
	totalResult := database.DB.Model(&models.Instances{}).Count(&total)
	if totalResult.Error != nil {
		l.Error("list instances error: %v", totalResult.Error)
		middleware.Error(ctx, middleware.CodeInstanceListError, iris.StatusInternalServerError)
		return
	}

	result := database.DB.Limit(limit).Offset(offset).Select("id", "user_id", "server_id", "cpu_only", "gpu_count", "volume_size", "ssh_address", "ssh_passwd", "jupyter_address", "tensor_board_address", "grafana_address", "code_server_address", "status", "image_name", "label", "created_at").Order("id").Find(&instances)
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
		user := models.Users{
			ID: instance.UserID,
		}
		result := database.DB.Select("username").First(&user)
		if result.Error == nil {
			res[i].Username = user.Username
		} else {
			l.Error("query user %d error: %v", instance.UserID, result.Error)
		}
		server := models.Servers{
			ID: instance.ServerID,
		}
		result = database.DB.Select("name", "gpu_type", "gpu_num", "cpu_count_per_gpu", "memory_per_gpu").First(&server)
		if result.Error == nil {

			redis.RawDB.Get(ctx, "remain_gpu:server:"+strconv.Itoa(int(instance.ServerID))).Scan(&server.GpuUsed)
			server.GpuUsed = server.GpuNum - server.GpuUsed

			res[i].ServerName = server.Name
			res[i].GpuType = server.GpuType
			res[i].GpuNum = server.GpuNum
			res[i].GpuUsed = server.GpuUsed
			res[i].CpuCountPerGpu = server.CpuCountPerGpu
			res[i].MemoryPerGpu = server.MemoryPerGpu
		} else {
			l.Error("query server %d error: %v", instance.ServerID, result.Error)
		}
	}

	middleware.ResultWithTotal(ctx, res, total)
}
