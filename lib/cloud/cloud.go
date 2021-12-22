package cloud

var RelativePathToTfexe string = "tfinstall"

type Cloud interface {
	CreateCluster() error
	DestroyCluster() error
	ObtainKubeConfig() error
}
