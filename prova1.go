package main

// import (
// 	"context"
// 	"flag"
// 	"fmt"
// 	"os"

// 	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
// 	"k8s.io/apimachinery/pkg/fields"
// 	"k8s.io/apimachinery/pkg/runtime/schema"
// 	"k8s.io/client-go/discovery"
// 	"k8s.io/client-go/dynamic"
// 	"k8s.io/client-go/rest"
// 	"k8s.io/client-go/tools/clientcmd"

// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// )

// // func listGroupVersions() {
// // 	// create a Kubernetes REST client config
// // 	config, err := rest.InClusterConfig()
// // 	if err != nil {
// // 		fmt.Fprintf(os.Stderr, "Error creating REST client config: %v", err)
// // 		return
// // 	}

// // 	// create a new discovery client
// // 	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
// // 	if err != nil {
// // 		fmt.Fprintf(os.Stderr, "Error creating discovery client: %v", err)
// // 		return
// // 	}

// // 	// get a list of all the server-supported API groups
// // 	groups, err := discoveryClient.ServerGroups()
// // 	if err != nil {
// // 		fmt.Fprintf(os.Stderr, "Error getting server groups: %v", err)
// // 		return
// // 	}

// // 	// iterate over the groups and print their versions
// // 	for _, group := range groups.Groups {
// // 		fmt.Printf("Group: %s\n", group.Name)
// // 		for _, version := range group.Versions {
// // 			fmt.Printf("  Version: %s\n", version.Version)
// // 		}
// // 	}
// // }

// func guessObject(dynamicClient *dynamic.DynamicClient, namespace string, objName string) (string, string, error) {
// 	// create a Kubernetes REST client config
// 	config, err := rest.InClusterConfig()
// 	if err != nil {
// 		return "", "", fmt.Errorf("Error creating REST client config: %v", err)
// 	}

// 	// create a new discovery client
// 	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
// 	if err != nil {
// 		return "", "", fmt.Errorf("Error creating discovery client: %v", err)
// 	}

// 	apiResources, err := discoveryClient.ServerGroups()
// 	if err != nil {
// 		return "", "", fmt.Errorf("Error getting server groups: %v", err)
// 	}

// 	// Loop through possible group versions and try to list objects with the given name
// 	for _, apiResource := range apiResources.Groups {
// 		for _, gv := range apiResource.Versions {
// 			objList, err := dynamicClient.Resource(
// 				schema.GroupVersionResource{
// 					Group:    apiResource.Name,
// 					Version:  gv.Version,
// 					Resource: objName + "s",
// 				},
// 			).Namespace(namespace).List(
// 				context.Background(),
// 				metav1.ListOptions{Limit: 1},
// 			)
// 			if err == nil && len(objList.Items) > 0 {
// 				return apiResource.Name, gv.Version, nil // found in group
// 			}
// 		}
// 	}
// 	return "", "", fmt.Errorf("object '%s' not found in any group/version", objName)
// }

// func main() {
// 	// Define flags
// 	// objectPtr := flag.String("object", "", "Kubernetes object to query")
// 	// fieldPtr := flag.String("field", "", "Object field to check")
// 	// expectedPtr := flag.String("expected", "", "Expected value for object field")
// 	countPtr := flag.Bool("count", false, "Count objects with the same value")
// 	notPtr := flag.Bool("not", false, "Inverse expression")
// 	namespacePtr := flag.String("namespace", "", "Namespace to use")

// 	// Parse command-line arguments
// 	flag.Parse()

// 	args := flag.Args()

// 	// Check if required arguments are provided
// 	if *objectPtr == "" || *fieldPtr == "" {
// 		fmt.Println("Object and field arguments are required")
// 		os.Exit(1)
// 	}

// 	// Set up Kubernetes client
// 	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
// 		clientcmd.NewDefaultClientConfigLoadingRules(),
// 		&clientcmd.ConfigOverrides{},
// 	)
// 	config, err := kubeconfig.ClientConfig()
// 	if err != nil {
// 		fmt.Printf("Error loading kubeconfig: %s", err)
// 		os.Exit(1)
// 	}
// 	dynamicClient, err := dynamic.NewForConfig(config)
// 	if err != nil {
// 		fmt.Printf("Error creating Kubernetes client: %s", err)
// 		os.Exit(1)
// 	}

// 	// Build field selector for objects
// 	fieldSelector := fields.OneTermEqualSelector(*fieldPtr, *expectedPtr).String()

// 	guessed_group, guessed_version, err := guessObject(dynamicClient, *namespacePtr, *objectPtr)
// 	if err != nil {
// 		fmt.Printf("Error guessing obj's group: %s", err)
// 		os.Exit(1)
// 	}

// 	// Query Kubernetes API
// 	obj := &unstructured.Unstructured{}
// 	obj.SetGroupVersionKind(schema.GroupVersionKind{
// 		Group:   guessed_group,
// 		Version: guessed_version,
// 		Kind:    *objectPtr,
// 	})

// 	objList, err := dynamicClient.Resource(
// 		schema.GroupVersionResource{
// 			Group:    guessed_group,
// 			Version:  guessed_version,
// 			Resource: *objectPtr + "s",
// 		},
// 	).Namespace(*namespacePtr).List(
// 		context.Background(),
// 		metav1.ListOptions{FieldSelector: fieldSelector},
// 	)

// 	if err != nil {
// 		fmt.Printf("Error querying Kubernetes API: %s", err)
// 		os.Exit(1)
// 	}

// 	// Parse object list and check field values
// 	objs := objList.Items
// 	if *countPtr {
// 		// Count objects with the same value
// 		count := len(objs)
// 		if *expectedPtr != "" {
// 			count = 0
// 			for _, obj := range objs {
// 				val, found, err := unstructured.NestedString(obj.Object, *fieldPtr)
// 				if err != nil || !found || val != *expectedPtr {
// 					continue
// 				}
// 				count++
// 			}
// 		}
// 		if (*notPtr && count != 0) || (!*notPtr && count != len(objs)) {
// 			fmt.Printf("%d objects with %s=%s\n", count, *fieldPtr, *expectedPtr)
// 			os.Exit(1)
// 		}
// 	} else {
// 		// Check if all objects have the expected value
// 		expectedValue := *expectedPtr
// 		if expectedValue == "" {
// 			val, found, err := unstructured.NestedString(objs[0].Object, *fieldPtr)
// 			if err != nil || !found {
// 				fmt.Printf("Error getting field value for object 0: %s", err)
// 				os.Exit(1)
// 			}
// 			expectedValue = val
// 		}
// 		for i, obj := range objs {
// 			val, found, err := unstructured.NestedString(obj.Object, *fieldPtr)
// 			if err != nil || !found || val != expectedValue {
// 				if (*notPtr && (err != nil || found || val == expectedValue)) || (!*notPtr && (err == nil && found && val == expectedValue)) {
// 					fmt.Printf("%s.%s is not %s for object %d: %s\n", *objectPtr, *fieldPtr, expectedValue, i, err)
// 					os.Exit(1)
// 				}
// 			}
// 		}
// 	}
// }
