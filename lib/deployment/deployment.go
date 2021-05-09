package deployment

import (
	"bytes"
	"context"
	"io"
	"text/template"

	g "github.com/sdslabs/katana/configs"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

func ApplyManifest(kubeconfig *rest.Config, kubeclient *kubernetes.Clientset, manifest []byte) error {
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

		gr, err := restmapper.GetAPIGroupResources(kubeclient.Discovery())
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
				unstructuredObj.SetNamespace("default")
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

func DeployCluster(kubeconfig *rest.Config, kubeclient *kubernetes.Clientset) error {
	clusterConfig := g.ClusterConfig

	teamConfig := TeamConfig{
		TeamCount: clusterConfig.TeamCount,
		TeamLabel: clusterConfig.TeamLabel,
	}

	broadcastConfig := BroadcastConfig{
		BroadcastCount: clusterConfig.BroadcastCount,
		BroadcastLabel: clusterConfig.BroadcastLabel,
		BroadcastPort:  g.ServicesConfig.ChallengeDeployer.BroadcastPort,
	}

	broadcastServiceConfig := &BroadcastServiceConfig{
		BroadcastLabel: clusterConfig.BroadcastLabel,
		BroadcastPort:  g.ServicesConfig.ChallengeDeployer.BroadcastPort,
	}

	teamtmpl, err := template.ParseFiles("templates/manifests/teams.yml")
	if err != nil {
		return err
	}

	broadcasttmpl, err := template.ParseFiles("templates/manifests/broadcast.yml")
	if err != nil {
		return err
	}

	broadcastsvctmpl, err := template.ParseFiles("templates/manifests/broadcast-service.yml")
	if err != nil {
		return err
	}

	teamManifest := &bytes.Buffer{}
	broadcastManifest := &bytes.Buffer{}
	broadcastServiceManifest := &bytes.Buffer{}

	if err = teamtmpl.Execute(teamManifest, teamConfig); err != nil {
		return err
	}

	if err = broadcasttmpl.Execute(broadcastManifest, broadcastConfig); err != nil {
		return err
	}

	if err = broadcastsvctmpl.Execute(broadcastServiceManifest, broadcastServiceConfig); err != nil {
		return err
	}

	if err = ApplyManifest(kubeconfig, kubeclient, teamManifest.Bytes()); err != nil {
		return err
	}

	if err = ApplyManifest(kubeconfig, kubeclient, broadcastManifest.Bytes()); err != nil {
		return err
	}

	return ApplyManifest(kubeconfig, kubeclient, broadcastServiceManifest.Bytes())
}
