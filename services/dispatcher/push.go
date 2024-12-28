package dispatcher

import (
	"context"
	"encoding/json"
	"megrez/services/redis"
	"strconv"
)

func Push(serverID uint, data Data) (err error) {
	l.SetFunction("Push")

	ctx := context.Background()

	dataBytes, err := json.Marshal(data)
	if err != nil {
		l.Error("Failed to marshal data, Server ID: %d, Data: %+v, Error: %v", serverID, data, err)
		return
	}

	_, err = redis.RawDB.RPush(ctx, "dispatcher:server:"+strconv.Itoa(int(serverID)), string(dataBytes)).Result()

	if err != nil {
		l.Error("Failed to push data to redis, Server ID: %d, Data: %+v, Error: %v", serverID, data, err)
		return
	}

	l.Info("Push data to redis success, Server ID: %d, Data: %+v", serverID, data)

	go run(serverID, false)

	return
}
