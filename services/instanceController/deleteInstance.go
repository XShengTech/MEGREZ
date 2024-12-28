package instanceController

import (
	"encoding/json"
	"errors"
	"megrez/libs/request"
	"strconv"
)

func deleteInstance(ip string, port int, apikey string,
	instanceName string,
) (err error) {
	l.SetFunction("deleteInstance")

	c := request.NewRequest().Delete().
		SetUrl("http://" + ip + ":" + strconv.Itoa(port) + apiPrefix + instancePrefix + "/" + instanceName).
		SetAuthorization("Bearer " + apikey).
		SetUserAgent("megrez")
	c.Do()

	if c.GetStatusCode() != 200 {
		l.Error("delete instance error code: %d", c.GetStatusCode())
		return errors.New("delete instance request error")
	}

	var res resStruct
	err = json.Unmarshal(c.GetBody(), &res)
	if err != nil {
		l.Error("unmarshal response data error: %v", err)
		return
	}

	if res.Code != 200 {
		l.Error("delete instance code: %d, error: %s", res.Code, res.Msg)
		return errors.New(res.Msg)
	}

	return
}
