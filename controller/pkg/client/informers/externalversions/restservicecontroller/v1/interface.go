package v1

import (
	internalinterfaces "github.com/Rest-service/controller/pkg/client/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// RestServicees returns a RestServiceInformer.
	RestServicees() RestServiceInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// RestServicees returns a RestServiceInformer.
func (v *version) RestServicees() RestServiceInformer {
	return &restserviceInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
