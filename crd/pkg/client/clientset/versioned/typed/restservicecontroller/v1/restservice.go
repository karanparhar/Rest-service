package v1

import (
	v1 "github.com/Rest-service/crd/pkg/apis/restservicecontroller/v1"
	scheme "github.com/Rest-service/crd/pkg/client/clientset/versioned/scheme"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// RestServiceesGetter has a method to return a RestServiceInterface.
// A group's client should implement this interface.
type RestServiceesGetter interface {
	RestServicees(namespace string) RestServiceInterface
}

// RestServiceInterface has methods to work with RestService resources.
type RestServiceInterface interface {
	Create(*v1.RestService) (*v1.RestService, error)
	Update(*v1.RestService) (*v1.RestService, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.RestService, error)
	List(opts meta_v1.ListOptions) (*v1.RestServiceList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.RestService, err error)
	RestServiceExpansion
}

// restservicees implements RestServiceInterface
type restservicees struct {
	client rest.Interface
	ns     string
}

// newRestServicees returns a RestServicees
func newRestServicees(c *RestServicecontrollerV1Client, namespace string) *restservicees {
	return &restservicees{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the restservice, and returns the corresponding restservice object, and an error if there is any.
func (c *restservicees) Get(name string, options meta_v1.GetOptions) (result *v1.RestService, err error) {
	result = &v1.RestService{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("Karan").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of RestServicees that match those selectors.
func (c *restservicees) List(opts meta_v1.ListOptions) (result *v1.RestServiceList, err error) {
	result = &v1.RestServiceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("Karan").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested restservicees.
func (c *restservicees) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("Karan").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a restservice and creates it.  Returns the server's representation of the restservice, and an error, if there is any.
func (c *restservicees) Create(restservice *v1.RestService) (result *v1.RestService, err error) {
	result = &v1.RestService{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("Karan").
		Body(restservice).
		Do().
		Into(result)
	return
}

// Update takes the representation of a restservice and updates it. Returns the server's representation of the restservice, and an error, if there is any.
func (c *restservicees) Update(restservice *v1.RestService) (result *v1.RestService, err error) {
	result = &v1.RestService{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("Karan").
		Name(restservice.Name).
		Body(restservice).
		Do().
		Into(result)
	return
}

// Delete takes name of the restservice and deletes it. Returns an error if one occurs.
func (c *restservicees) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("Karan").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *restservicees) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("Karan").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched restservice.
func (c *restservicees) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.RestService, err error) {
	result = &v1.RestService{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("Karan").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
