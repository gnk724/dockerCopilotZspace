package module

import (
	"github.com/onlyLTY/dockerCopilotZspace/zspace/internal/svc"
	"github.com/onlyLTY/dockerCopilotZspace/zspace/internal/utiles"
	"github.com/zeromicro/go-zero/core/logx"
)

const autoRestartContainerPerception = 85.0

func AutoRestartContainer(ctx *svc.ServiceContext) {
	containerList, ErrGetContainerList := utiles.GetContainerList(ctx, true)
	if ErrGetContainerList != nil {
		logx.Errorf("panic获取镜像列表出错: %v", ErrGetContainerList)
		return
	}
	for _, v := range containerList {
		containerStat, ErrGetContainerStat := utiles.GetContainerStat(ctx, v.ID)
		if ErrGetContainerStat != nil {
			logx.Errorf("获取容器状态出错: %v", ErrGetContainerStat)
			continue
		}
		if containerStat.UsageMemoryPercent > autoRestartContainerPerception {
			ErrRestartContainer := utiles.RestartContainer(ctx, v.ID)
			if ErrRestartContainer != nil {
				logx.Errorf("重启容器出错: %v", ErrRestartContainer)
			}
		}
	}
}
