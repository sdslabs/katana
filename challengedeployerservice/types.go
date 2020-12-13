package challengedeployerservice

type ChallengeDeployerConfig struct {
	TeamLabel     string
	BroadcastPort int
	KubeHost      string
	KubeNameSpace string
	KubeConfig    string
}
