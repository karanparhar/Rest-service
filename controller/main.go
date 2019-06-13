package main

import (
	"flag"
	"time"

	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	controller "github.com/Rest-service/controller/controller"
	kube "github.com/Rest-service/controller/kube"
	clientset "github.com/Rest-service/controller/pkg/client/clientset/versioned"
	informers "github.com/Rest-service/controller/pkg/client/informers/externalversions"
	"github.com/Rest-service/controller/pkg/signals"
	kubeinformers "k8s.io/client-go/informers"
)

var (
	kubeconfig string
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
}
func main() {

	var kubeClient kubernetes.Interface
	var cfg *rest.Config

	cfg, err := rest.InClusterConfig()
	if err != nil {
		kubeClient, cfg = kube.GetClientOutOfCluster(kubeconfig)
	} else {
		kubeClient, cfg = kube.GetClient()
	}

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	exampleClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	exampleInformerFactory := informers.NewSharedInformerFactory(exampleClient, time.Second*30)

	controller := controller.NewController(kubeClient, exampleClient, kubeInformerFactory, exampleInformerFactory)

	go kubeInformerFactory.Start(stopCh)
	go exampleInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}
