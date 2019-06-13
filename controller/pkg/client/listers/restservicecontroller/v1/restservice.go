package v1

import (
	v1 "github.com/Rest-service/controller/pkg/apis/restservicecontroller/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// RestServiceLister helps list RestServicees.
type RestServiceLister interface {
	// List lists all RestServicees in the indexer.
	List(selector labels.Selector) (ret []*v1.RestService, err error)
	// RestServicees returns an object that can list and get RestServicees.
	RestServicees(namespace string) RestServiceNamespaceLister
	RestServiceListerExpansion
}

// restserviceLister implements the RestServiceLister interface.
type restserviceLister struct {
	indexer cache.Indexer
}

// NewRestServiceLister returns a new RestServiceLister.
func NewRestServiceLister(indexer cache.Indexer) RestServiceLister {
	return &restserviceLister{indexer: indexer}
}

// List lists all RestServicees in the indexer.
func (s *restserviceLister) List(selector labels.Selector) (ret []*v1.RestService, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.RestService))
	})
	return ret, err
}

// RestServicees returns an object that can list and get RestServicees.
func (s *restserviceLister) RestServicees(namespace string) RestServiceNamespaceLister {
	return restserviceNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// RestServiceNamespaceLister helps list and get RestServicees.
type RestServiceNamespaceLister interface {
	// List lists all RestServicees in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.RestService, err error)
	// Get retrieves the RestService from the indexer for a given namespace and name.
	Get(name string) (*v1.RestService, error)
	RestServiceNamespaceListerExpansion
}

// restserviceNamespaceLister implements the RestServiceNamespaceLister
// interface.
type restserviceNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all RestServicees in the indexer for a given namespace.
func (s restserviceNamespaceLister) List(selector labels.Selector) (ret []*v1.RestService, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.RestService))
	})
	return ret, err
}

// Get retrieves the RestService from the indexer for a given namespace and name.
func (s restserviceNamespaceLister) Get(name string) (*v1.RestService, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("test.crd"), name)
	}
	return obj.(*v1.RestService), nil
}
