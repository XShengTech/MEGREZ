package dispatcher

import (
	"context"
	"encoding/json"
	"megrez/services/redis"
	"strconv"
)

func run(serverID uint, recursion bool) {
	lc := l.Clone()
	lc.SetFunction("run:" + strconv.Itoa(int(serverID)))

	if !recursion && status[int(serverID)] {
		return
	}

	mx.Lock()
	status[int(serverID)] = true
	mx.Unlock()

	ctx := context.Background()

	llen, err := redis.RawDB.LLen(ctx, "dispatcher:server:"+strconv.Itoa(int(serverID))).Result()
	if err != nil {
		lc.Error("Failed to get length of list, Server ID: %d, Error: %v", serverID, err)
		mx.Lock()
		status[int(serverID)] = false
		mx.Unlock()
		return
	}

	if llen == 0 {
		mx.Lock()
		status[int(serverID)] = false
		mx.Unlock()
		l.Debug("No data in redis, Server ID: %d", serverID)
		return
	}

	dataStr, err := redis.RawDB.LPop(ctx, "dispatcher:server:"+strconv.Itoa(int(serverID))).Result()
	if err != nil {
		lc.Error("Failed to pop data from redis, Server ID: %d, Error: %v", serverID, err)
		mx.Lock()
		status[int(serverID)] = false
		mx.Unlock()

		go run(serverID, true)
		return
	}

	var data Data
	err = json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		lc.Error("Failed to unmarshal data, Server ID: %d, Data: %s, Error: %v", serverID, dataStr, err)
		mx.Lock()
		status[int(serverID)] = false
		mx.Unlock()

		run(serverID, true)
		return
	}

	switch data.Type {
	case Add:
		err = add(serverID, data)
		if err != nil {
			lc.Error("Failed to add instance, Server ID: %d, Data: %+v, Error: %v", serverID, data, err)
		}
	case Control:
		err = control(serverID, data)
		if err != nil {
			lc.Error("Failed to control instance, Server ID: %d, Data: %+v, Error: %v", serverID, data, err)
		}
	case Modify:
		err = modify(serverID, data)
		if err != nil {
			lc.Error("Failed to modify instance, Server ID: %d, Data: %+v, Error: %v", serverID, data, err)
		}
	case Delete:
		err = delete(serverID, data)
		if err != nil {
			lc.Error("Failed to delete instance, Server ID: %d, Data: %+v, Error: %v", serverID, data, err)
		}
	default:
	}

	go run(serverID, true)

	return
}
