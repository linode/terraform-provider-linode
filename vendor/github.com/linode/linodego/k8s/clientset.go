package k8s

import (
	"encoding/base64"
	"fmt"

	"github.com/linode/linodego"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/transport"
)

// NewClientsetFromBytes builds a Clientset from a given Kubeconfig.
//
// Takes an optional transport.WrapperFunc to add request/response middleware to
// api-server requests.
func BuildClientsetFromConfig(
	lkeKubeconfig *linodego.LKEClusterKubeconfig,
	transportWrapper transport.WrapperFunc,
) (kubernetes.Interface, error) {
	kubeConfigBytes, err := base64.StdEncoding.DecodeString(lkeKubeconfig.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to decode kubeconfig: %s", err)
	}

	restClientConfig, err := clientcmd.RESTConfigFromKubeConfig(kubeConfigBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse LKE cluster kubeconfig: %s", err)
	}

	if transportWrapper != nil {
		restClientConfig.Wrap(transportWrapper)
	}

	clientset, err := kubernetes.NewForConfig(restClientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build k8s client from LKE cluster kubeconfig: %s", err)
	}
	return clientset, nil
}
