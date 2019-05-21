package fake

import (
	restservicecontroller_v1 "github.com/AIDI/rpa/crd/pkg/apis/restservicecontroller/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeRestServicees implements RestServiceInterface
type FakeRestServicees struct {
	Fake *FakeRestServicecontrollerV1
	ns   string
}

var restserviceesResource = schema.GroupVersionResource{Group: "kloud9", Version: "v1", Resource: "Karan"}

var restserviceesKind = schema.GroupVersionKind{Group: "kloud9", Version: "v1", Kind: "RestService"}

// Get takes name of the restservice, and returns the corresponding restservice object, and an error if there is any.
func (c *FakeRestServicees) Get(name string, options v1.GetOptions) (result *restservicecontroller_v1.RestService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(restserviceesResource, c.ns, name), &restservicecontroller_v1.RestService{})

	if obj == nil {
		return nil, err
	}
	return obj.(*restservicecontroller_v1.RestService), err
}

// List takes label and field selectors, and returns the list of RestServicees that match those selectors.
func (c *FakeRestServicees) List(opts v1.ListOptions) (result *restservicecontroller_v1.RestServiceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(restserviceesResource, restserviceesKind, c.ns, opts), &restservicecontroller_v1.RestServiceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &restservicecontroller_v1.RestServiceList{}
	for _, item := range obj.(*restservicecontroller_v1.RestServiceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested restservicees.
func (c *FakeRestServicees) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(restserviceesResource, c.ns, opts))

}

// Create takes the representation of a restservice and creates it.  Returns the server's representation of the restservice, and an error, if there is any.
func (c *FakeRestServicees) Create(restservice *restservicecontroller_v1.RestService) (result *restservicecontroller_v1.RestService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(restserviceesResource, c.ns, restservice), &restservicecontroller_v1.RestService{})

	if obj == nil {
		return nil, err
	}
	return obj.(*restservicecontroller_v1.RestService), err
}

// Update takes the representation of a restservice and updates it. Returns the server's representation of the restservice, and an error, if there is any.
func (c *FakeRestServicees) Update(restservice *restservicecontroller_v1.RestService) (result *restservicecontroller_v1.RestService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(restserviceesResource, c.ns, restservice), &restservicecontroller_v1.RestService{})

	if obj == nil {
		return nil, err
	}
	return obj.(*restservicecontroller_v1.RestService), err
}

// Delete takes name of the restservice and deletes it. Returns an error if one occurs.
func (c *FakeRestServicees) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(restserviceesResource, c.ns, name), &restservicecontroller_v1.RestService{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRestServicees) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(restserviceesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &restservicecontroller_v1.RestServiceList{})
	return err
}

// Patch applies the patch and returns the patched restservice.
func (c *FakeRestServicees) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *restservicecontroller_v1.RestService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(restserviceesResource, c.ns, name, data, subresources...), &restservicecontroller_v1.RestService{})

	if obj == nil {
		return nil, err
	}
	return obj.(*restservicecontroller_v1.RestService), err
}
