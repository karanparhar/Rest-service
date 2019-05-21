package v1

import (
	v1 "github.com/AIDI/rpa/crd/pkg/apis/restservicecontroller/v1"
	"github.com/AIDI/rpa/crd/pkg/client/clientset/versioned/scheme"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	rest "k8s.io/client-go/rest"
)

type RestServicecontrollerV1Interface interface {
	RESTClient() rest.Interface
	RestServiceesGetter
}

// RestServicecontrollerV1Client is used to interact with features provided by the restservicecontroller.kubeplus group.
type RestServicecontrollerV1Client struct {
	restClient rest.Interface
}

func (c *RestServicecontrollerV1Client) RestServicees(namespace string) RestServiceInterface {
	return newRestServicees(c, namespace)
}

// NewForConfig creates a new RestServicecontrollerV1Client for the given config.
func NewForConfig(c *rest.Config) (*RestServicecontrollerV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &RestServicecontrollerV1Client{client}, nil
}

// NewForConfigOrDie creates a new RestServicecontrollerV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *RestServicecontrollerV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new RestServicecontrollerV1Client for the given RESTClient.
func New(c rest.Interface) *RestServicecontrollerV1Client {
	return &RestServicecontrollerV1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *RestServicecontrollerV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
