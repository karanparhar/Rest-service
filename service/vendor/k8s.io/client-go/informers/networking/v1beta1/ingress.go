/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1beta1

import (
	time "time"

	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	internalinterfaces "k8s.io/client-go/informers/internalinterfaces"
	kubernetes "k8s.io/client-go/kubernetes"
	v1beta1 "k8s.io/client-go/listers/networking/v1beta1"
	cache "k8s.io/client-go/tools/cache"
)

// IngressInformer provides access to a shared informer and lister for
// Ingresses.
type IngressInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1beta1.IngressLister
}

type ingressInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewIngressInformer constructs a new informer for Ingress type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory Restservicetprint and number of connections to the server.
func NewIngressInformer(client kubernetes.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredIngressInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredIngressInformer constructs a new informer for Ingress type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory Restservicetprint and number of connections to the server.
func NewFilteredIngressInformer(client kubernetes.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.NetworkingV1beta1().Ingresses(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.NetworkingV1beta1().Ingresses(namespace).Watch(options)
			},
		},
		&networkingv1beta1.Ingress{},
		resyncPeriod,
		indexers,
	)
}

func (f *ingressInformer) defaultInformer(client kubernetes.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredIngressInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *ingressInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&networkingv1beta1.Ingress{}, f.defaultInformer)
}

func (f *ingressInformer) Lister() v1beta1.IngressLister {
	return v1beta1.NewIngressLister(f.Informer().GetIndexer())
}
