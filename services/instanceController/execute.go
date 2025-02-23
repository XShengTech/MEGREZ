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

func SetJupyterPassword(ip string, port int, apikey string,
	containerName, password string) (err error) {
	l.SetFunction("SetJupyterPassword")

	// Set Jupyter Password
	data := executeReq{
		Cmd: []string{
			"sed",
			"-i",
			"/^c.ServerApp.token = ''/c\\c.ServerApp.token = '" + password + "'",
			"/root/.jupyter/jupyter_notebook_config.py",
		},
	}

	reqBytes, err := json.Marshal(data)
	if err != nil {
		l.Error("marshal request data error: %v", err)
		return
	}

	c := request.NewRequest().Post().
		SetUrl("http://" + ip + ":" + strconv.Itoa(port) + apiPrefix + instancePrefix + "/" + containerName + instanceExecute).
		SetAuthorization("Bearer " + apikey).
		SetUserAgent("megrez").
		SetBody(bytes.NewBuffer(reqBytes))
	c.Do()

	if c.GetStatusCode() != 200 {
		l.Error("set jupyter password error: %d", c.GetStatusCode())
		return errors.New("set jupyter password request error")
	}

	var res resStruct
	err = json.Unmarshal(c.GetBody(), &res)
	if err != nil {
		l.Error("unmarshal response data error: %v", err)
		return err
	}

	if res.Code != 200 {
		l.Error("set jupyter password code: %d, error: %s", res.Code, res.Msg)
		return errors.New(res.Msg)
	}

	// Restart Jupyter
	data = executeReq{
		Cmd: []string{
			"service",
			"jupyter",
			"restart",
		},
	}

	reqBytes, err = json.Marshal(data)
	if err != nil {
		l.Error("marshal request data error: %v", err)
		return
	}

	c = request.NewRequest().Post().
		SetUrl("http://" + ip + ":" + strconv.Itoa(port) + apiPrefix + instancePrefix + "/" + containerName + instanceExecute).
		SetAuthorization("Bearer " + apikey).
		SetUserAgent("megrez").
		SetBody(bytes.NewBuffer(reqBytes))
	c.Do()

	if c.GetStatusCode() != 200 {
		l.Error("restart jupyter error: %d", c.GetStatusCode())
		return errors.New("restart jupyter request error")
	}

	err = json.Unmarshal(c.GetBody(), &res)
	if err != nil {
		l.Error("unmarshal response data error: %v", err)
		return err
	}

	if res.Code != 200 {
		l.Error("restart jupyter code: %d, error: %s", res.Code, res.Msg)
		return errors.New(res.Msg)
	}

	return nil
}

func SetCodeServerPassword(ip string, port int, apikey string,
	containerName, password string) (err error) {
	l.SetFunction("SetCodeServerPassword")

	// Set Code-Server Password
	data := executeReq{
		Cmd: []string{
			"sed",
			"-i",
			"/^password: CHANGETHIS/c\\password: " + password,
			"/root/.config/code-server/config.yaml",
		},
	}

	reqBytes, err := json.Marshal(data)
	if err != nil {
		l.Error("marshal request data error: %v", err)
		return
	}

	c := request.NewRequest().Post().
		SetUrl("http://" + ip + ":" + strconv.Itoa(port) + apiPrefix + instancePrefix + "/" + containerName + instanceExecute).
		SetAuthorization("Bearer " + apikey).
		SetUserAgent("megrez").
		SetBody(bytes.NewBuffer(reqBytes))
	c.Do()

	if c.GetStatusCode() != 200 {
		l.Error("set code-server password error: %d", c.GetStatusCode())
		return errors.New("set code-server password request error")
	}

	var res resStruct
	err = json.Unmarshal(c.GetBody(), &res)
	if err != nil {
		l.Error("unmarshal response data error: %v", err)
		return err
	}

	if res.Code != 200 {
		l.Error("set code-server password code: %d, error: %s", res.Code, res.Msg)
		return errors.New(res.Msg)
	}

	// Restart Code-Server
	data = executeReq{
		Cmd: []string{
			"service",
			"code-server",
			"restart",
		},
	}

	reqBytes, err = json.Marshal(data)
	if err != nil {
		l.Error("marshal request data error: %v", err)
		return
	}

	c = request.NewRequest().Post().
		SetUrl("http://" + ip + ":" + strconv.Itoa(port) + apiPrefix + instancePrefix + "/" + containerName + instanceExecute).
		SetAuthorization("Bearer " + apikey).
		SetUserAgent("megrez").
		SetBody(bytes.NewBuffer(reqBytes))
	c.Do()

	if c.GetStatusCode() != 200 {
		l.Error("restart code-server error: %d", c.GetStatusCode())
		return errors.New("restart code-server request error")
	}

	err = json.Unmarshal(c.GetBody(), &res)
	if err != nil {
		l.Error("unmarshal response data error: %v", err)
		return err
	}

	if res.Code != 200 {
		l.Error("restart code-server code: %d, error: %s", res.Code, res.Msg)
		return errors.New(res.Msg)
	}

	return nil
}
