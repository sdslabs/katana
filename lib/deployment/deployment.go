package deployment

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sdslabs/katana/logging"
	"io"
	"log"
	"path/filepath"
	"text/template"

	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/utils"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

var logger = &logging.GlobalLogger

// ApplyManifest applies a given manifest to the cluster
func ApplyManifest(kubeconfig *rest.Config, kubeclientset *kubernetes.Clientset, manifest []byte, namespace string) error {
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
		if err != nil {
			return err
		}
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
				unstructuredObj.SetNamespace(namespace)
			}
			dri = dd.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = dd.Resource(mapping.Resource)
		}

		if _, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{}); err != nil {
			if _, err := dri.Update(context.Background(), unstructuredObj, metav1.UpdateOptions{}); err != nil {
				if unstructuredObj.GetObjectKind().GroupVersionKind().Kind == "PersistentVolumeClaim" {
					// Skip PVCs
					continue
					// TODO: Handle PVCs, currently on deletion of PVCs, the cluster is stuck in a loop
				}
				_ = dri.Delete(context.Background(), unstructuredObj.GetName(), metav1.DeleteOptions{})
				watcher, err := dri.Watch(context.Background(), metav1.ListOptions{
					FieldSelector: fmt.Sprintf("metadata.name=%s", unstructuredObj.GetName()),
				})
				if err != nil {
					return err
				}
				defer watcher.Stop()
				for event := range watcher.ResultChan() {
					if event.Type == watch.Deleted {
						_, err = dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{})
						if err != nil {
							return err
						}
						break
					}
				}
			}
		}
	}

	if err != io.EOF {
		return err
	}
	return nil
}

// DeployCluster retrieves and applies all the manifest templates specified in config
// after injecting the necessary values
func DeployCluster(kubeconfig *rest.Config, kubeclientset *kubernetes.Clientset) error {
	clusterConfig := g.ClusterConfig

	deploymentConfig := utils.DeploymentConfig()

	clientset, _ := utils.GetKubeClient()

	nodes, err := utils.GetNodes(clientset)
	if err != nil {
		log.Println(err)
	}
	deploymentConfig.NodeAffinityValue = nodes[0].Name

	for _, m := range clusterConfig.TemplatedManifests {
		manifest := &bytes.Buffer{}
		logger.Info().Msgf("Applying: %s\n", m)
		tmpl, err := template.ParseFiles(filepath.Join(clusterConfig.TemplatedManifestDir, m))
		if err != nil {
			return err
		}
		if err = tmpl.Execute(manifest, deploymentConfig); err != nil {
			return err
		}
		//TODO add a receiver handler to avoid throttling or apply all manifest concurrently
		if err = ApplyManifest(kubeconfig, kubeclientset, manifest.Bytes(), g.KatanaConfig.KubeNameSpace); err != nil {
			return err
		}

	}

	return nil
}

func DeployChallengeToCluster(challengeName, teamName string, firstPatch bool, replicas int32) error {

	teamNamespace := teamName + "-ns"
	kubeclient, _ := utils.GetKubeClient()

	//to-do verify image exist in harbor, wait for sometime if not return error
	deploymentsClient := kubeclient.AppsV1().Deployments(teamNamespace)
	imageName := "harbor.katana.local/katana/" + challengeName + ":latest"
	if firstPatch {
		/// Retrieve the existing deployment
		existingDeployment, err := deploymentsClient.Get(context.TODO(), challengeName, metav1.GetOptions{})
		if err != nil {
			logger.Error().Msgf("Error in retrieving existing deployment.")
			logger.Error().Err(err)
			return err
		}

		existingDeployment.Spec.Template.Spec.Containers[0].Image = "harbor.katana.local/katana/" + teamName + "-" + challengeName + ":latest"

		_, err = deploymentsClient.Update(context.TODO(), existingDeployment, metav1.UpdateOptions{})
		if err != nil {
			logger.Error().Msgf("Error in updating deployment.")
			logger.Error().Err(err)
			return err
		}

		logger.Info().Msgf("Updated deployment with new image.")
		return nil
	}

	manifest := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: teamNamespace,
			Name:      challengeName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": challengeName,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": challengeName,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            challengeName + "-" + teamName,
							Image:           imageName,
							ImagePullPolicy: v1.PullPolicy("IfNotPresent"),
							Ports: []v1.ContainerPort{
								{
									Name:          "challenge-port",
									ContainerPort: 3000,
								},
							},
						},
					},
				},
			},
		},
	}
	log.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), manifest, metav1.CreateOptions{})

	if err != nil {
		log.Println("Unable to create deployement")
		panic(err)
	}

	log.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName()+" in namespace "+teamNamespace)
	return nil
}

func DeployChallengeCheckerToCluster(challengeCheckerName, namespace string, replicas int32) error {

	kubeclient, _ := utils.GetKubeClient()

	deploymentsClient := kubeclient.AppsV1().Deployments(namespace)
	imageName := "harbor.katana.local/katana/" + challengeCheckerName

	manifest := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      challengeCheckerName + "-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": challengeCheckerName,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": challengeCheckerName,
					},
				},
				Spec: v1.PodSpec{
					ServiceAccountName: "cluster-admin",
					Containers: []v1.Container{
						{
							Name:            challengeCheckerName,
							Image:           imageName,
							ImagePullPolicy: v1.PullPolicy("Always"),
							Ports: []v1.ContainerPort{
								{
									Name:          "checker-port",
									ContainerPort: 8080,
								},
							},
						},
					},
				},
			},
		},
	}

	logger.Info().Msgf("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), manifest, metav1.CreateOptions{})

	if err != nil {
		logger.Error().Msgf("Unable to create deployment")
		panic(err)
	}

	logger.Info().Msgf("Created deployment %q.\n", result.GetObjectMeta().GetName()+" in namespace "+namespace)
	return nil
}
