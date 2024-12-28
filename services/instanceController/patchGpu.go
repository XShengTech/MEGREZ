package instanceController

import (
	"bytes"
	"encoding/json"
	"errors"
	"megrez/libs/request"
	"strconv"
)

func patchGpu(ip string, port int, apikey string,
	instanceName string,
	cpuCountPerGpu int, memoryPerGpu int,
	volumeName string, oldVolumeName string,
	gpuCount int, oldGpuCount int,
) (err error) {
	l.SetFunction("patchGpu")

	var data patchReqStruct

	if gpuCount == oldGpuCount && volumeName == oldVolumeName {
		return nil
	}

	if gpuCount != oldGpuCount {
		data.GpuPatch = &gpuPatchStruct{
			GpuCount: gpuCount,
		}
		data.CpuPatch = &cpuPatchStruct{
			CpuCount: cpuCountPerGpu * gpuCount,
		}
		data.MemoryPatch = &MemoryPatchStruct{
			Memory: strconv.Itoa(memoryPerGpu*gpuCount) + "GB",
		}
	}

	if volumeName != oldVolumeName {
		data.VolumePatch = &volumePatchStruct{
			OldBind: bindStruct{
				Src:  oldVolumeName,
				Dest: "/root/megrez-tmp",
			},
			NewBind: bindStruct{
				Src:  volumeName,
				Dest: "/root/megrez-tmp",
			},
		}
	}

	reqBytes, err := json.Marshal(data)
	if err != nil {
		l.Error("marshal request data error: %v", err)
		return err
	}

	c := request.NewRequest().Patch().
		SetUrl("http://" + ip + ":" + strconv.Itoa(port) + apiPrefix + instancePrefix + "/" + instanceName).
		SetAuthorization("Bearer " + apikey).
		SetUserAgent("megrez").
		SetBody(bytes.NewBuffer(reqBytes))
	c.Do()

	if c.GetStatusCode() != 200 {
		l.Error("patch gpu error code: %d", c.GetStatusCode())
		return errors.New("patch gpu request error")
	}

	var res resStruct
	err = json.Unmarshal(c.GetBody(), &res)
	if err != nil {
		l.Error("unmarshal response data error: %v", err)
		return err
	}

	if res.Code != 200 {
		l.Error("patch gpu code: %d, error: %s", res.Code, res.Msg)
		return errors.New(res.Msg)
	}

	return nil
}
