package challengedeployerservice

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/sdslabs/katana/lib/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	g "github.com/sdslabs/katana/configs"
)


var (
	config       = g.ChallengeDeployerConfig
	KatanaConfig = g.KatanaConfig
	kubeclient   *kubernetes.Clientset
)

// Get kubernetes client
func GetClient(pathToCfg string) error {
	if pathToCfg == "" {
		pathToCfg = filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)
	}
	config, err := clientcmd.BuildConfigFromFlags("", pathToCfg)
	if err != nil {
		return err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	kubeclient = client
	return nil
}

func getPods(lbls map[string]string) ([]v1.Pod, error) {
	set := labels.Set(lbls)
	pods, err := kubeclient.CoreV1().Pods(KatanaConfig.KubeNameSpace).List(context.Background(), metav1.ListOptions{LabelSelector: set.AsSelector().String()})
	if err != nil {
		return nil, err
	}
	return pods.Items, nil
}

// Check if path already exists or not
func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Clone the provided remote repository into repos/<local>
func Clone(remote string, local string, auth *githttp.BasicAuth) error {
	cloneConfig := &git.CloneOptions{
		URL:      remote,
		Progress: os.Stdout,
		Auth:     auth,
	}

	tmpdir, err := ioutil.TempDir("tmp", local)
	if err != nil {
		return err
	}

	if _, err := git.PlainClone(fmt.Sprintf(tmpdir), false, cloneConfig); err != nil {
		return err
	}
	return compressAndMove(tmpdir, fmt.Sprintf("challenges/%s.zip", local))
}

// Compress the src directory into <dst>.zip
func compressAndMove(src string, dst string) error {
	if exists(dst) {
		return fmt.Errorf("File %s: already exists or cannot be accessed", dst)
	}

	outfile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer outfile.Close()

	w := zip.NewWriter(outfile)
	defer w.Close()

	n := len(src)
	if src[n-1] == '/' {
		src = src[:n-1]
		n -= 1
	}

	addToZip := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// TODO: Check for arbitrary read/ write

		f, err := w.Create(path[n:])
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}

	if err = filepath.Walk(src, addToZip); err != nil {
		return err
	}
	return os.RemoveAll(src)
}

// Send file to given URI and include provided parameters in the request
func SendFile(file *os.File, params map[string]string, filename, uri string) error {
	client := &http.Client{}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(config.ArtifactLabel, filename)
	if err != nil {
		return err
	}

	if _, err = io.Copy(part, file); err != nil {
		return err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	if err = writer.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client.Do(req)
	return nil
}

// broadcast sends given file to the broadcast service, to be forwarded to all pods marked with app=cofig.Teamlabel
func Broadcast(file string) error {
	chal, err := os.Open(filepath.Join("challenges", file))
	if err != nil {
		return err
	}
	defer chal.Close()

	lbls := utils.GetTeamPodLabels()
	teamPods, err := getPods(lbls)
	if err != nil {
		return err
	}

	addresses := []string{}
	for _, pod := range teamPods {
		addresses = append(addresses, fmt.Sprintf("%s:%d", pod.Status.PodIP, config.TeamClientPort))
	}

	addressesEncoded, err := json.Marshal(addresses)
	if err != nil {
		return err
	}

	params := make(map[string]string)
	params["targets"] = string(addressesEncoded)

	return SendFile(chal, params, file, fmt.Sprintf("%s:%d", KatanaConfig.KubeHost, config.BroadcastPort))
}
