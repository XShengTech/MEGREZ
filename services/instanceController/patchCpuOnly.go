package instanceController

import (
	"bytes"
	"encoding/json"
	"errors"
	"megrez/libs/request"
	"strconv"
)

func patchCpuOnly(ip string, port int, apikey string,
	instanceName string,
	volumeName string, oldVolumeName string,
) (err error) {
	l.SetFunction("patchCpuOnly")

	data := patchReqStruct{
		CpuPatch: &cpuPatchStruct{
			CpuCount: 1,
		},
		GpuPatch: &gpuPatchStruct{
			GpuCount: 0,
		},
		MemoryPatch: &MemoryPatchStruct{
			Memory: "2GB",
		},
	}

	if oldVolumeName != volumeName {
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
		l.Error("patch cpu only error code: %d", c.GetStatusCode())
		return errors.New("patch cpu only request error")
	}

	var res resStruct
	err = json.Unmarshal(c.GetBody(), &res)
	if err != nil {
		l.Error("unmarshal response data error: %v", err)
		return err
	}

	if res.Code != 200 {
		l.Error("patch cpu only code: %d, error: %s", res.Code, res.Msg)
		return errors.New(res.Msg)
	}

	return nil
}
