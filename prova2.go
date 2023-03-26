package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/tidwall/sjson"
)

func main() {
	var objectType, objectName, objectNamespace, selector, value string
	flag.StringVar(&objectType, "type", "", "Kubernetes object type")
	flag.StringVar(&objectName, "name", "", "Kubernetes object name")
	flag.StringVar(&objectNamespace, "namespace", "default", "Kubernetes object namespace")
	flag.StringVar(&selector, "selector", "", "Kubernetes field selector")
	flag.StringVar(&value, "value", "", "Expected value for the selected field")
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
	options := v1.GetOptions{}
	if selector != "" {
		options = v1.GetOptions{
			FieldSelector: selector,
		}
	}
	object, err := clientset.CoreV1().RESTClient().Get().
		Namespace(objectNamespace).
		Resource(objectType + "s").
		Name(objectName).
		VersionedParams(&options, v1.ParameterCodec).
		Do(context.Background()).
		Raw()
	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Printf("Kubernetes %s %s not found in namespace %s\n", objectType, objectName, objectNamespace)
			return
		}
		panic(err.Error())
	}

	// Use jq to filter the output based on the selector and the expected value
	if value != "" {
jsonString := string(object)

// Use sjson to get the value of the selector
value, err := sjson.Get(jsonString, selector)
if err != nil {
    fmt.Printf("Error getting value for selector %s: %s\n", selector, err.Error())
    return
}

// Check if the value matches the expected value
if value != expectedValue {
    fmt.Printf("Kubernetes %s %s in namespace %s does not have the expected value %s for the field selector %s\n", objectType, objectName, objectNamespace, expectedValue, selector)
    return
}
	}

	// Output the JSON definition of the Kubernetes object
	fmt.Printf("%s\n", string(object))
}
