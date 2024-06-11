package controller

import (
	"sync"
	"time"

	typeV1 "github.com/istio-ecosystem/admiral-api/pkg/apis/admiral/v1"
	clientset "github.com/istio-ecosystem/admiral-api/pkg/client/clientset/versioned"
	"github.com/istio-ecosystem/admiral-api/pkg/client/informers/externalversions/admiral/v1"

	log "github.com/sirupsen/logrus"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type ShardController struct {
	KubernetesClientset kubernetes.Interface
	Informer            cache.SharedIndexInformer
	Mutex               sync.Mutex
	ShardCache          map[string]*typeV1.Shard
}

func NewShardController(stopCh <-chan struct{}, crdClientSet clientset.Interface, kubernetesClientset kubernetes.Interface, resyncPeriod time.Duration) (*ShardController, error) {

	informerFactory := informers.NewSharedInformerFactoryWithOptions(kubernetesClientset, resyncPeriod)
	informerFactory.Start(stopCh)

	shardInformer := v1.NewShardInformer(crdClientSet,
		metaV1.NamespaceAll,
		resyncPeriod,
		cache.Indexers{})

	shardInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			shard, ok := obj.(*typeV1.Shard)
			if !ok {
				log.Warn("shard type mismatch")
			}
			log.Infof(shard.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			shard, ok := newObj.(*typeV1.Shard)
			if !ok {
				log.Warn("shard type mismatch")
			}
			log.Infof(shard.Name)
		},
		DeleteFunc: func(obj interface{}) {
			shard, ok := obj.(*typeV1.Shard)
			if !ok {
				log.Warn("shard type mismatch")
			}
			log.Infof(shard.Name)
		},
	})

	return &ShardController{
		Informer:            shardInformer,
		KubernetesClientset: kubernetesClientset,
		Mutex:               sync.Mutex{},
		ShardCache:          make(map[string]*typeV1.Shard),
	}, nil
}
