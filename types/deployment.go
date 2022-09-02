package types

import (
	"context"

	"k8s.io/client-go/kubernetes"
)

// ManifestConfig contains the total data to be injected in manifest templates
type ManifestConfig struct {
	TeamCount             uint
	TeamLabel             string
	BroadcastCount        uint
	BroadcastLabel        string
	BroadcastPort         uint
	KubeNameSpace         string
	FluentHost            string
	ChallengeDir           string
	TempDir               string
	InitFile              string
	DaemonPort            uint
	ChallengeDeployerHost string
	ChallengeArtifact     string
}

type ResourceStatus struct {
	Name          string
	TotalReplicas int32
	ReadyReplicas int32
	Ready         bool
}

type ResourcePinger func(context.Context, *kubernetes.Clientset, map[string]string) ([]*ResourceStatus, bool, error)
