package instanceController

import (
	"bytes"
	"encoding/json"
	"errors"
	"megrez/libs/crypto"
	"megrez/libs/request"
	"strconv"
)

type createVolumeReqStruct struct {
	Name string `json:"name"`
	Size string `json:"size"`
}

func createVolume(ip string, port int, apikey string,
	volumeSize int,
) (volumeName string, err error) {
	l.SetFunction("createVolume")

	volumeName, err = crypto.GenerateUUIDWithoutHyphen()
	if err != nil {
		l.Error("generate volume name error: %v", err)
		return "", err
	}

	data := createVolumeReqStruct{
		Name: volumeName,
		Size: strconv.Itoa(volumeSize) + "GB",
	}

	reqBytes, err := json.Marshal(data)
	if err != nil {
		l.Error("marshal request data error: %v", err)
		return "", err
	}

	c := request.NewRequest().Post().
		SetUrl("http://" + ip + ":" + strconv.Itoa(port) + apiPrefix + volumePrefix).
		SetAuthorization("Bearer " + apikey).
		SetUserAgent("megrez").
		SetBody(bytes.NewBuffer(reqBytes))
	c.Do()

	if c.GetStatusCode() != 200 {
		l.Error("create volume error code: %d", c.GetStatusCode())
		return "", errors.New("create volume request error")
	}

	var res resStruct
	err = json.Unmarshal(c.GetBody(), &res)
	if err != nil {
		l.Error("unmarshal response data error: %v", err)
		return "", err
	}

	if res.Code != 200 {
		l.Error("create volume code: %d, error: %s", res.Code, res.Msg)
		return "", errors.New(res.Msg)
	}

	volumeName = volumeName + "-1"

	return
}
