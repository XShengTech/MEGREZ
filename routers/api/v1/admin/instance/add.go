package instances

import (
	"megrez/libs/crypto"
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/database"
	"megrez/services/dispatcher"
	"megrez/services/redis"
	"strconv"

	"github.com/kataras/iris/v12"
)

type addReqStruct struct {
	UserID   uint `json:"user_id"`
	ServerID uint `json:"server_id"`

	ImageName  string `json:"image_name"`
	GpuCount   int    `json:"gpu_count"`
	VolumeSize int    `json:"volume_size"`
}

func addHandler(ctx iris.Context) {
	l.SetFunction("addHandler")

	var req addReqStruct
	err := ctx.ReadJSON(&req)
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	if req.UserID == 0 || req.ServerID == 0 || req.ImageName == "" || req.GpuCount <= 0 || req.VolumeSize < 50 {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
	}

	user := models.Users{
		ID: req.UserID,
	}
	result := database.DB.First(&user)
	if result.Error != nil {
		l.Error("query user error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeAdminUserQueryError, iris.StatusInternalServerError)
		return
	}

	server := models.Servers{
		ID: req.ServerID,
	}
	result = database.DB.First(&server)
	if result.Error != nil {
		l.Error("query server error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeServerQueryError, iris.StatusInternalServerError)
		return
	}

	remainGpu, err := redis.RawDB.IncrBy(ctx, "remain_gpu:server:"+strconv.Itoa(int(server.ID)), int64(-req.GpuCount)).Result()
	if err != nil {
		l.Error("incrby gpu num error: %v", err)
		middleware.Error(ctx, middleware.CodeServerQueryError, iris.StatusInternalServerError)
		return
	}

	remainVolume, err := redis.RawDB.IncrBy(ctx, "remain_volume:server:"+strconv.Itoa(int(server.ID)), int64(-req.VolumeSize-30)).Result()
	if err != nil {
		l.Error("incrby volume size error: %v", err)
		middleware.Error(ctx, middleware.CodeServerQueryError, iris.StatusInternalServerError)
		return
	}

	if remainGpu < 0 || remainVolume < 0 {
		redis.RawDB.IncrBy(ctx, "remain_gpu:server:"+strconv.Itoa(int(server.ID)), int64(req.GpuCount))
		redis.RawDB.IncrBy(ctx, "remain_volume:server:"+strconv.Itoa(int(server.ID)), int64(req.VolumeSize+30))
		middleware.Error(ctx, middleware.CodeResourceInsufficient, iris.StatusBadRequest)
		return
	}

	instance := models.Instances{
		UserID:     req.UserID,
		ServerID:   req.ServerID,
		ImageName:  req.ImageName,
		GpuCount:   req.GpuCount,
		VolumeSize: req.VolumeSize,

		SshPasswd: crypto.Hex(16),

		Status: models.InstanceStatusReady,
	}
	result = database.DB.Create(&instance)
	if result.Error != nil {
		l.Error("create instance error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeInternalCreateError, iris.StatusInternalServerError)
		return
	}

	dispatcherData := dispatcher.Data{
		Type:       dispatcher.Add,
		InstanceID: instance.ID,
	}
	dispatcher.Push(instance.ServerID, dispatcherData)

	middleware.Success(ctx)
}
