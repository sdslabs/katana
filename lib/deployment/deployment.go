package deployment

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"text/template"

	g "github.com/sdslabs/katana/configs"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

// ApplyManifest applies a given manifest to the cluster
func ApplyManifest(kubeconfig *rest.Config, kubeclientset *kubernetes.Clientset, manifest []byte) error {
	dd, err := dynamic.NewForConfig(kubeconfig)
	if err != nil {
		return err
	}

	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(manifest), 100)
	for {
		var rawObj runtime.RawExtension
		if err = decoder.Decode(&rawObj); err != nil {
			break
		}

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return err
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		gr, err := restmapper.GetAPIGroupResources(kubeclientset.Discovery())
		if err != nil {
			return err
		}

		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return err
		}

		var dri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace(g.KatanaConfig.KubeNameSpace)
			}
			dri = dd.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = dd.Resource(mapping.Resource)
		}

		if _, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{}); err != nil {
			return err
		}
	}

	if err != io.EOF {
		return err
	} else {
		return nil
	}
}

// DeployCluster retrieves and applies all the manifest templates specified in config
// after injecting the necessary values
func DeployCluster(kubeconfig *rest.Config, kubeclientset *kubernetes.Clientset) error {
	clusterConfig := g.ClusterConfig

	deploymentConfig := DeploymentConfig{
		FluentHost:            fmt.Sprintf("\"elasticsearch.%s.svc.cluster.local\"", g.KatanaConfig.KubeNameSpace),
		KubeNameSpace:         g.KatanaConfig.KubeNameSpace,
		TeamCount:             clusterConfig.TeamCount,
		TeamLabel:             clusterConfig.TeamLabel,
		BroadcastCount:        clusterConfig.BroadcastCount,
		BroadcastLabel:        clusterConfig.BroadcastLabel,
		BroadcastPort:         g.ServicesConfig.ChallengeDeployer.BroadcastPort,
		ChallengDir:           g.TeamVmConfig.ChallengeDir,
		TempDir:               g.TeamVmConfig.TempDir,
		InitFile:              g.TeamVmConfig.InitFile,
		DaemonPort:            g.TeamVmConfig.DaemonPort,
		ChallengeDeployerHost: g.ServicesConfig.ChallengeDeployer.Host,
		ChallengeArtifact:     g.ServicesConfig.ChallengeDeployer.ArtifactLabel,
	}

	for _, m := range clusterConfig.Manifests {
		manifest := &bytes.Buffer{}
		fmt.Printf("Applying: %s\n", m)
		tmpl, err := template.ParseFiles(filepath.Join(clusterConfig.ManifestDir, m))
		if err != nil {
			return err
		}

		if err = tmpl.Execute(manifest, deploymentConfig); err != nil {
			return err
		}

		if err = ApplyManifest(kubeconfig, kubeclientset, manifest.Bytes()); err != nil {
			return err
		}
	}

	return nil
}

func PollDeployments(kubeclientset *kubernetes.Clientset, activePods chan<- string) (string, error) {
	client := kubeclientset.CoreV1()
	selector := make(map[string]string)
	selector["deployment"] = g.ClusterConfig.DeploymentLabel
	opts := metav1.ListOptions{LabelSelector: labels.Set(selector).AsSelector().String()}

	statusCount := make(map[corev1.PodPhase]uint)
	for {
		pods, err := client.Pods(g.KatanaConfig.KubeNameSpace).List(context.Background(), opts)
		if err != nil {
			return "", err
		}

		statusCount[corev1.PodFailed] = 0
		statusCount[corev1.PodPending] = 0
		statusCount[corev1.PodRunning] = 0
		statusCount[corev1.PodUnknown] = 0

		for _, pod := range pods.Items {
			statusCount[pod.Status.Phase] += 1
		}

		fmt.Print("\r \r")
		fmt.Printf("Running: %d\tPending: %d\tFailed: %d\tUnknown: %d", statusCount[corev1.PodRunning], statusCount[corev1.PodPending], statusCount[corev1.PodFailed], statusCount[corev1.PodFailed])
	}
}
