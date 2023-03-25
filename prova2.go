package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var objectType, objectName, objectNamespace string
	flag.StringVar(&objectType, "type", "", "Kubernetes object type")
	flag.StringVar(&objectName, "name", "", "Kubernetes object name")
	flag.StringVar(&objectNamespace, "namespace", "default", "Kubernetes object namespace")
	flag.Parse()

	if objectType == "" || objectName == "" {
		fmt.Println("Please provide both the object type and name")
		return
	}

	// Create a Kubernetes client
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Get the specified Kubernetes object
	object, err := clientset.CoreV1().RESTClient().Get().
		Namespace(objectNamespace).
		Resource(objectType + "s").
		Name(objectName).
		Do(context.Background()).
		Raw()
	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Printf("Kubernetes %s %s not found in namespace %s\n", objectType, objectName, objectNamespace)
			return
		}
		panic(err.Error())
	}

	// Output the JSON definition of the Kubernetes object
	fmt.Printf("%s\n", string(object))
}
