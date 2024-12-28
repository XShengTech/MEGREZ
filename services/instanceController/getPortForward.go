package instanceController

import (
	"encoding/json"
	"errors"
	"megrez/libs/request"
	"strconv"
	"strings"
)

type infoResStruct struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Info struct {
			HostConfig struct {
				PortBindings map[string][]struct {
					HostIP   string `json:"HostIp"`
					HostPort string `json:"HostPort"`
				}
			} `json:"HostConfig"`
		} `json:"info"`
	} `json:"data"`
}

func GetPortForward(ip string, port int, apikey string, containerName string) (portBindings map[string]string, err error) {
	l.SetFunction("GetPortForward")

	portBindings = make(map[string]string)

	c := request.NewRequest().Get().
		SetUrl("http://" + ip + ":" + strconv.Itoa(port) + apiPrefix + instancePrefix + "/" + containerName).
		SetAuthorization("Bearer " + apikey).
		SetUserAgent("megrez")
	c.Do()

	if c.GetStatusCode() != 200 {
		l.Error("get port forward error: %d", c.GetStatusCode())
		return nil, errors.New("get port forward request error")
	}

	var res infoResStruct
	err = json.Unmarshal(c.GetBody(), &res)
	if err != nil {
		l.Error("unmarshal response data error: %v", err)
		return
	}

	if res.Code != 200 {
		l.Error("get port forward code: %d, error: %s", res.Code, res.Msg)
		return nil, errors.New(res.Msg)
	}

	for k, v := range res.Data.Info.HostConfig.PortBindings {
		portBindings[strings.Split(k, "/")[0]] = v[0].HostPort
	}

	return portBindings, nil
}
