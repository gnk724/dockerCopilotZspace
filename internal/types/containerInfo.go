package types

import (
	docker "github.com/docker/docker/api/types"
)

type Container struct {
	docker.Container
	Update bool `json:"Update"`
}

type ContainerStat struct {
	docker.Stats
	UsageMemoryPercent float64 `json:"UsageMemoryPercent"`
}
