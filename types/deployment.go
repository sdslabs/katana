package types

import (
	"context"

	"k8s.io/client-go/kubernetes"
)

// ManifestConfig contains the total data to be injected in manifest templates
type ManifestConfig struct {
	TeamCount         uint
	TeamLabel         string
	KubeNameSpace     string
	TeamPodName       string
	ContainerName     string
	ChallengDir       string
	TempDir           string
	InitFile          string
	DaemonPort        uint
	MongoUsername     string
	MongoPassword     string
	MySQLPassword     string
	SSHPassword       string
	HarborKey         string
	HarborCrt         string
	HarborCaCrt       string
	HarborIP          string
	WireguardIP       string
	NodeAffinityValue string
}

type ResourceStatus struct {
	Name          string
	TotalReplicas int32
	ReadyReplicas int32
	Ready         bool
}

type ResourcePinger func(context.Context, *kubernetes.Clientset, map[string]string) ([]*ResourceStatus, bool, error)

type Repo struct {
	FullName string `json:"full_name"`
}

type GogsRequest struct {
	Ref        string `json:"ref"`
	Before     string `json:"before"`
	Repository Repo   `json:"repository"`
}
