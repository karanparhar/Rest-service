package v1

import (
	time "time"

	restservicecontroller_v1 "github.com/Rest-service/crd/pkg/apis/restservicecontroller/v1"
	versioned "github.com/Rest-service/crd/pkg/client/clientset/versioned"
	internalinterfaces "github.com/Rest-service/crd/pkg/client/informers/externalversions/internalinterfaces"
	v1 "github.com/Rest-service/crd/pkg/client/listers/restservicecontroller/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// RestServiceInformer provides access to a shared informer and lister for
// RestServicees.
type RestServiceInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.RestServiceLister
}

type restserviceInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewRestServiceInformer constructs a new informer for RestService type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewRestServiceInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredRestServiceInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredRestServiceInformer constructs a new informer for RestService type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredRestServiceInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.RestServicecontrollerV1().RestServicees(namespace).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.RestServicecontrollerV1().RestServicees(namespace).Watch(options)
			},
		},
		&restservicecontroller_v1.RestService{},
		resyncPeriod,
		indexers,
	)
}

func (f *restserviceInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredRestServiceInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *restserviceInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&restservicecontroller_v1.RestService{}, f.defaultInformer)
}

func (f *restserviceInformer) Lister() v1.RestServiceLister {
	return v1.NewRestServiceLister(f.Informer().GetIndexer())
}
