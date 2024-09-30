package utiles

import (
	"encoding/json"
	"github.com/onlyLTY/dockerCopilotZspace/zspace/internal/svc"
	MyType "github.com/onlyLTY/dockerCopilotZspace/zspace/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/url"
)

func GetContainerList(ctx *svc.ServiceContext, filterRunning bool) ([]MyType.Container, error) {
	jwtToken, endpointsId, err := GetNewJwt(ctx)
	if err != nil {
		logx.Errorf("GetNewJwt error: %v", err)
		logx.Errorf("请检查环境的account变量以及是否为hosts模式")
		return nil, err
	}
	client := NewCustomClient(jwtToken)
	baseURL := domain + "/api/endpoints/" + endpointsId
	URL := baseURL + "/docker/containers/json"
	queryParams := url.Values{}
	if filterRunning {
		queryParams.Add("all", "false")
	} else {
		queryParams.Add("all", "true")
	}
	URL += "?" + queryParams.Encode()
	resp, err := client.SendRequest("GET", URL, nil)
	if err != nil {
		logx.Errorf("SendRequest error: %v", err)
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.Error("io.ReadAll error: %v", err)
		return nil, err
	}
	var dockerContainerList []MyType.Container
	err = json.Unmarshal(body, &dockerContainerList)
	if err != nil {
		logx.Errorf("json.Unmarshal error: %v", err)
		return nil, err
	}
	return dockerContainerList, nil
}

func CheckImageUpdate(ctx *svc.ServiceContext, containerListData []MyType.Container) []MyType.Container {
	for i, v := range containerListData {
		if _, ok := ctx.HubImageInfo.Data[v.ImageID]; ok {
			if ctx.HubImageInfo.Data[v.ImageID].NeedUpdate {
				containerListData[i].Update = true
			}
		}
	}
	return containerListData
}
