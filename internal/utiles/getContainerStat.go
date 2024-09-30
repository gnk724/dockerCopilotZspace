package utiles

import (
	"encoding/json"
	"github.com/onlyLTY/dockerCopilotZspace/zspace/internal/svc"
	"github.com/onlyLTY/dockerCopilotZspace/zspace/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
)

func GetContainerStat(ctx *svc.ServiceContext, id string) (stat types.ContainerStat, err error) {
	jwtToken, endpointsId, err := GetNewJwt(ctx)
	if err != nil {
		logx.Errorf("GetNewJwt error: %v", err)
		return
	}
	client := NewCustomClient(jwtToken)
	baseURL := domain + "/api/endpoints/" + endpointsId
	url := baseURL + "/docker/containers/" + id + "/stats?stream=false"
	resp, err := client.SendRequest("GET", url, nil)
	if err != nil {
		logx.Errorf("SendRequest error: %v", err)
		return
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			logx.Errorf("Body.Close error: %v", err)
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.Error("io.ReadAll error: %v", err)
		return
	}
	err = json.Unmarshal(body, &stat)
	if err != nil {
		logx.Errorf("json.Unmarshal error: %v", err)
		return
	}
	usedMemory := stat.MemoryStats.Usage - stat.MemoryStats.Stats["cache"]
	availableMemory := stat.MemoryStats.Limit
	if availableMemory == 0 {
		stat.UsageMemoryPercent = 0
	} else {
		stat.UsageMemoryPercent = float64(usedMemory) / float64(availableMemory) * 100
	}
	usageMemoryPercent := float64(usedMemory) / float64(availableMemory) * 100
	logx.Infof("name: %+v", id)
	logx.Infof("stat: %+v", usageMemoryPercent)
	return
}
