package cloud

var PathToCloudPackage string = "/lib/cloud"
var RelativePathToTfexe string = "/tfinstall"
var ExecPath string
var PathToAzureTf string = "/azure"

type Cloud interface {
	CreateCluster() error
	DestroyCluster() error
	ObtainKubeConfig() error
}
