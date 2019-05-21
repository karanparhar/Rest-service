package fake

import (
	v1 "github.com/AIDI/rpa/crd/pkg/client/clientset/versioned/typed/restservicecontroller/v1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeRestServicecontrollerV1 struct {
	*testing.Fake
}

func (c *FakeRestServicecontrollerV1) RestServicees(namespace string) v1.RestServiceInterface {
	return &FakeRestServicees{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeRestServicecontrollerV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
