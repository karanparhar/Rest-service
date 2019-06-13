package versioned

import (
	restservicecontrollerv1 "github.com/Rest-service/controller/pkg/client/clientset/versioned/typed/restservicecontroller/v1"
	glog "github.com/golang/glog"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	RestServicecontrollerV1() restservicecontrollerv1.RestServicecontrollerV1Interface
	// Deprecated: please explicitly pick a version if possible.
	RestServicecontroller() restservicecontrollerv1.RestServicecontrollerV1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	restservicecontrollerV1 *restservicecontrollerv1.RestServicecontrollerV1Client
}

// RestServicecontrollerV1 retrieves the RestServicecontrollerV1Client
func (c *Clientset) RestServicecontrollerV1() restservicecontrollerv1.RestServicecontrollerV1Interface {
	return c.restservicecontrollerV1
}

// Deprecated: RestServicecontroller retrieves the default version of RestServicecontrollerClient.
// Please explicitly pick a version.
func (c *Clientset) RestServicecontroller() restservicecontrollerv1.RestServicecontrollerV1Interface {
	return c.restservicecontrollerV1
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.restservicecontrollerV1, err = restservicecontrollerv1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		glog.Errorf("failed to create the DiscoveryClient: %v", err)
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.restservicecontrollerV1 = restservicecontrollerv1.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.restservicecontrollerV1 = restservicecontrollerv1.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
