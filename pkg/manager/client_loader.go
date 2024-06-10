package manager

import (
	admiralv1 "github.com/istio-ecosystem/admiral-api/pkg/client/clientset/versioned/typed/admiral/v1"
	"k8s.io/client-go/rest"
)

type LoadKubeClient interface {
	LoadAdmiralApiClientFromPath(path string) (admiralv1.AdmiralV1Interface, error)
	LoadAdmiralApiClientFromConfig(config *rest.Config) (admiralv1.AdmiralV1Interface, error)
}

type KubeClient struct{}

func (loader *KubeClient) LoadAdmiralApiClientFromPath(kubeConfigPath string) (admiralv1.AdmiralV1Interface, error) {
	return nil, nil
}

func (loader *KubeClient) LoadAdmiralApiClientFromConfig(config *rest.Config) (admiralv1.AdmiralV1Interface, error) {
	return nil, nil
}
