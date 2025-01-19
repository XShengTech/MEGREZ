package instanceController

import (
	"bytes"
	"encoding/json"
	"errors"
	"megrez/libs/crypto"
	"megrez/libs/request"
	"strconv"
)

type createInstanceReqStruct struct {
	ImageName      string       `json:"imageName"`
	ReplicaSetName string       `json:"replicaSetName"`
	CpuCount       int          `json:"cpuCount"`
	GpuCount       int          `json:"gpuCount"`
	Memory         string       `json:"memory"`
	Binds          []bindStruct `json:"binds"`
	Env            []string     `json:"env"`
	ContainerPorts []string     `json:"containerPorts"`
}

func createInstance(ip string, port int, apikey string,
	gpuCount, cpuCount, memorySize int, volumeName string, imageName string,
) (containerName string, err error) {
	l.SetFunction("createInstance")

	containerName, err = crypto.GenerateUUIDWithoutHyphen()
	if err != nil {
		l.Error("generate instance name error: %v", err)
		return "", err
	}

	data := createInstanceReqStruct{
		ImageName:      imageName,
		ReplicaSetName: containerName,
		CpuCount:       cpuCount,
		GpuCount:       gpuCount,
		Memory:         strconv.Itoa(memorySize) + "GB",
		ContainerPorts: []string{
			"22",    // SSH
			"6007",  // TensorBoard
			"8888",  // Jupyter Notebook
			"3000",  // Grafana
			"8080",  // Code-Server
			"34567", // Custom Port
		},
		Env: []string{
			"NVIDIA_DRIVER_CAPABILITIES=video,compute,utility",
		},
	}

	if volumeName != "" {
		data.Binds = []bindStruct{
			{
				Src:  volumeName,
				Dest: "/root/megrez-tmp",
			},
			{
				Src:  "/data/pub",
				Dest: "/root/megrez-pub",
			},
		}
	}

	reqBytes, err := json.Marshal(data)
	if err != nil {
		l.Error("marshal request data error: %v", err)
		return "", err
	}

	c := request.NewRequest().Post().
		SetUrl("http://" + ip + ":" + strconv.Itoa(port) + apiPrefix + instancePrefix).
		SetAuthorization("Bearer " + apikey).
		SetUserAgent("megrez").
		SetBody(bytes.NewBuffer(reqBytes))
	c.Do()

	if c.GetStatusCode() != 200 {
		l.Error("create instance error: %d", c.GetStatusCode())
		return "", errors.New("create instance request error")
	}

	var res resStruct
	err = json.Unmarshal(c.GetBody(), &res)
	if err != nil {
		l.Error("unmarshal response data error: %v", err)
		return "", err
	}

	if res.Code != 200 {
		l.Error("create instance code: %d, error: %s", res.Code, res.Msg)
		return "", errors.New(res.Msg)
	}

	return
}
