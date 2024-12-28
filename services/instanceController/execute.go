package instanceController

import (
	"bytes"
	"encoding/json"
	"errors"
	"megrez/libs/request"
	"strconv"
)

type executeReq struct {
	Cmd []string `json:"cmd"`
}

func SetRootPassword(ip string, port int, apikey string,
	containerName, password string) (err error) {
	l.SetFunction("SetRootPassword")

	data := executeReq{
		Cmd: []string{
			"sh",
			"-c",
			"echo root:" + password + " | chpasswd",
		},
	}

	reqBytes, err := json.Marshal(data)
	if err != nil {
		l.Error("marshal request data error: %v", err)
		return err
	}

	c := request.NewRequest().Post().
		SetUrl("http://" + ip + ":" + strconv.Itoa(port) + apiPrefix + instancePrefix + "/" + containerName + instanceExecute).
		SetAuthorization("Bearer " + apikey).
		SetUserAgent("megrez").
		SetBody(bytes.NewBuffer(reqBytes))
	c.Do()

	if c.GetStatusCode() != 200 {
		l.Error("set root password error: %d", c.GetStatusCode())
		return errors.New("set root password request error")
	}

	var res resStruct
	err = json.Unmarshal(c.GetBody(), &res)
	if err != nil {
		l.Error("unmarshal response data error: %v", err)
		return err
	}

	if res.Code != 200 {
		l.Error("set root password code: %d, error: %s", res.Code, res.Msg)
		return errors.New(res.Msg)
	}

	return nil
}
