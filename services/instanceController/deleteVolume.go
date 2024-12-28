package instanceController

import (
	"encoding/json"
	"errors"
	"megrez/libs/request"
	"strconv"
	"strings"
)

func deleteVolume(ip string, port int, apikey string,
	volumeName string, noAll bool,
) (err error) {
	l.SetFunction("deleteVolume")

	url := "http://" + ip + ":" + strconv.Itoa(port) + apiPrefix + volumePrefix + "/" + strings.Split(volumeName, "-")[0]
	if noAll {
		url = "http://" + ip + ":" + strconv.Itoa(port) + apiPrefix + volumePrefix + "/" + volumeName + "?noall=true"
	}

	c := request.NewRequest().Delete().
		SetUrl(url).
		SetAuthorization("Bearer " + apikey).
		SetUserAgent("megrez")
	c.Do()

	if c.GetStatusCode() != 200 {
		l.Error("delete volume error code: %d", c.GetStatusCode())
		return errors.New("delete volume request error")
	}

	var res resStruct
	err = json.Unmarshal(c.GetBody(), &res)
	if err != nil {
		l.Error("unmarshal response data error: %v", err)
		return
	}

	if res.Code != 200 {
		l.Error("delete volume code: %d, error: %s", res.Code, res.Msg)
		return errors.New(res.Msg)
	}

	return nil
}
