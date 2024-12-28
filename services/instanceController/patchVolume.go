package instanceController

import (
	"bytes"
	"encoding/json"
	"errors"
	"megrez/libs/request"
	"strconv"
)

type patchVolumeReqStruct struct {
	Size string `json:"size"`
}

type patchVolumeResStruct struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Name string `json:"name"`
		Size string `json:"size"`
	} `json:"data"`
}

func patchVolume(ip string, port int, apikey string,
	volumeName string, newVolumeSize int,
) (newVolumeName string, err error) {
	l.SetFunction("patchVolume")

	data := patchVolumeReqStruct{
		Size: strconv.Itoa(newVolumeSize) + "GB",
	}

	reqBytes, err := json.Marshal(data)
	if err != nil {
		l.Error("marshal request data error: %v", err)
		return "", err
	}

	c := request.NewRequest().Patch().
		SetUrl("http://" + ip + ":" + strconv.Itoa(port) + apiPrefix + volumePrefix + "/" + volumeName + volumeSize).
		SetAuthorization("Bearer " + apikey).
		SetUserAgent("megrez").
		SetBody(bytes.NewBuffer(reqBytes))
	c.Do()

	if c.GetStatusCode() != 200 {
		l.Error("patch volume error code: %d", c.GetStatusCode())
		return "", errors.New("patch volume request error")
	}

	var res patchVolumeResStruct
	err = json.Unmarshal(c.GetBody(), &res)
	if err != nil {
		l.Error("unmarshal response data error: %v", err)
		return "", err
	}

	if res.Code != 200 {
		l.Error("patch volume code: %d, error: %s", res.Code, res.Msg)
		return "", errors.New(res.Msg)
	}

	return res.Data.Name, nil
}
