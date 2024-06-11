package manager

import (
	"fmt"
	admiralv1 "github.com/istio-ecosystem/admiral-api/pkg/client/clientset/versioned/typed/admiral/v1"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type LoadKubeClient interface {
	LoadAdmiralApiClientFromPath(path string) (admiralv1.AdmiralV1Interface, error)
	LoadAdmiralApiClientFromConfig(config *rest.Config) (admiralv1.AdmiralV1Interface, error)
}

type KubeClient struct{}

func (loader *KubeClient) LoadAdmiralApiClientFromPath(kubeConfigPath string) (admiralv1.AdmiralV1Interface, error) {
	config, err := getConfig(kubeConfigPath)
	if err != nil || config == nil {
		return nil, err
	}

	return loader.LoadAdmiralApiClientFromConfig(config)
}

func (loader *KubeClient) LoadAdmiralApiClientFromConfig(config *rest.Config) (admiralv1.AdmiralV1Interface, error) {
	return admiralv1.NewForConfig(config)
}

func getConfig(kubeConfigPath string) (*rest.Config, error) {
	logrus.Infof("getting kubeconfig from: %#v", kubeConfigPath)
	// create the config from the path
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)

	if err != nil || config == nil {
		return nil, fmt.Errorf("could not retrieve kubeconfig: %v", err)
	}
	return config, err
}
