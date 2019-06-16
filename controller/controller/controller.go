package controller

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"

	"github.com/golang/glog"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	resource "k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	restservicev1 "github.com/Rest-service/controller/pkg/apis/restservicecontroller/v1"
	clientset "github.com/Rest-service/controller/pkg/client/clientset/versioned"
	restservicescheme "github.com/Rest-service/controller/pkg/client/clientset/versioned/scheme"
	informers "github.com/Rest-service/controller/pkg/client/informers/externalversions"
	listers "github.com/Rest-service/controller/pkg/client/listers/restservicecontroller/v1"
)

const controllerAgentName = "controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a Restservice is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Restservice fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Restservice"
	// MessageResourceSynced is the message used for an Event fired when a Restservice
	// is synced successfully
	MessageResourceSynced = "Restservice synced successfully"
)

// Controller is the controller implementation for Restservice resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// sampleclientset is a clientset for our own API group
	sampleclientset clientset.Interface

	deploymentsLister appslisters.DeploymentLister
	deploymentsSynced cache.InformerSynced
	restserviceLister listers.RestServiceLister
	restserviceSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController returns a new sample controller
func NewController(
	kubeclientset kubernetes.Interface,
	sampleclientset clientset.Interface,
	kubeInformerFactory kubeinformers.SharedInformerFactory,
	sampleInformerFactory informers.SharedInformerFactory) *Controller {

	// obtain references to shared index informers for the Deployment and Restservice
	// types.
	deploymentInformer := kubeInformerFactory.Apps().V1().Deployments()
	restserviceInformer := sampleInformerFactory.RestServicecontroller().V1().RestServicees()

	// Create event broadcaster
	// Add restservice-controller types to the default Kubernetes Scheme so Events can be
	// logged for restservice-controller types.
	restservicescheme.AddToScheme(scheme.Scheme)
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:     kubeclientset,
		sampleclientset:   sampleclientset,
		deploymentsLister: deploymentInformer.Lister(),
		deploymentsSynced: deploymentInformer.Informer().HasSynced,
		restserviceLister: restserviceInformer.Lister(),
		restserviceSynced: restserviceInformer.Informer().HasSynced,
		workqueue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Karan"),
		recorder:          recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when Restservice resources change
	restserviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueRestservice,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueRestservice(new)
		},
	})
	// Set up an event handler for when Deployment resources change. This
	// handler will lookup the owner of the given Deployment, and if it is
	// owned by a Restservice resource will enqueue that Restservice resource for
	// processing. This way, we don't need to implement custom logic for
	// handling Deployment resources. More info on this pattern:
	// https://github.com/kubernetes/community/blob/8cafef897a22026d42f5e5bb3f104febe7e29830/contributors/devel/controllers.md
	deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.addDeployment,
		UpdateFunc: controller.updateDeployment,
		DeleteFunc: controller.deleteDeployment,
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting Restservice controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.deploymentsSynced, c.restserviceSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch two workers to process Restservice resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Restservice resource to be synced.
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Restservice resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the Restservice resource with this namespace/name
	Restservice, err := c.restserviceLister.RestServicees(namespace).Get(name)
	if err != nil {
		// The Restservice resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("Restservice '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}
	deploymentName := Restservice.Spec.DeploymentName
	namespaceName := Restservice.Spec.NameSpace
	if deploymentName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		runtime.HandleError(fmt.Errorf("%s: deployment name must be specified", key))
		return nil
	}

	// Get the deployment with the name specified in Restservice.spec
	deployment, err := c.deploymentsLister.Deployments(namespaceName).Get(deploymentName)
	// If the resource doesn't exist, we'll create it
	if err != nil {
		if errors.IsNotFound(err) {
			deployment, err = c.kubeclientset.AppsV1().Deployments(namespaceName).Create(newDeployment(Restservice))

			return err
		}
	}

	// If the Deployment is not controlled by this Restservice resource, we should log
	// a warning to the event recorder and ret
	if !metav1.IsControlledBy(deployment, Restservice) {
		msg := fmt.Sprintf(MessageResourceExists, deployment.Name)
		c.recorder.Event(Restservice, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}
	// Finally, we update the status block of the Restservice resource to reflect the
	// current state of the world
	err = c.updateRestserviceStatus(Restservice, deployment)
	if err != nil {
		return err
	}

	c.recorder.Event(Restservice, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *Controller) updateRestserviceStatus(Restservice *restservicev1.RestService, deployment *appsv1.Deployment) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	RestserviceCopy := Restservice.DeepCopy()
	RestserviceCopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas
	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the Restservice resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := c.sampleclientset.RestServicecontrollerV1().RestServicees(Restservice.Namespace).Update(RestserviceCopy)
	return err
}

// enqueueRestservice takes a Restservice resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Restservice.
func (c *Controller) enqueueRestservice(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the Restservice resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that Restservice resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (c *Controller) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			runtime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			runtime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		glog.V(4).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	glog.V(4).Infof("Processing object: %s", object.GetName())
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		if ownerRef.Kind != "RestService" {
			return
		}

		Restservice, err := c.restserviceLister.RestServicees(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			glog.V(4).Infof("ignoring orphaned object '%s' of Restservice '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		c.enqueueRestservice(Restservice)
		return
	}
}

func (c *Controller) addDeployment(obj interface{}) {
	d := obj.(*appsv1.Deployment)
	klog.V(4).Infof("Adding deployment %s", d.Name)
	c.handleObject(d)

}

func (c *Controller) updateDeployment(old, cur interface{}) {
	oldD := old.(*appsv1.Deployment)
	curD := cur.(*appsv1.Deployment)
	klog.V(4).Infof("Updating deployment %s", oldD.Name)
	c.handleObject(curD)
}

func (c *Controller) deleteDeployment(obj interface{}) {
	d, ok := obj.(*appsv1.Deployment)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			runtime.HandleError(fmt.Errorf("Couldn't get object from tombstone %#v", obj))
			return
		}
		d, ok = tombstone.Obj.(*appsv1.Deployment)
		if !ok {
			runtime.HandleError(fmt.Errorf("Tombstone contained object that is not a Deployment %#v", obj))
			return
		}
	}
	klog.V(4).Infof("Deleting deployment %s", d.Name)
	c.handleObject(d)
}

// newDeployment creates a new Deployment for a Restservice resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Restservice resource that 'owns' it.
func newDeployment(Restservice *restservicev1.RestService) *appsv1.Deployment {
	labels := map[string]string{
		"app":        Restservice.Spec.DeploymentName,
		"controller": Restservice.Name,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Restservice.Spec.DeploymentName,
			Namespace: Restservice.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(Restservice, schema.GroupVersionKind{
					Group:   restservicev1.SchemeGroupVersion.Group,
					Version: restservicev1.SchemeGroupVersion.Version,
					Kind:    "RestService",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: Restservice.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            Restservice.Spec.DeploymentName,
							Image:           Restservice.Spec.Image,
							ImagePullPolicy: "Never",

							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("1000m"),
									corev1.ResourceMemory: resource.MustParse("500Mi"),
								},
							},
						},
					},
				},
			},
		},
	}
}

func int32Ptr(i int32) *int32 { return &i }
