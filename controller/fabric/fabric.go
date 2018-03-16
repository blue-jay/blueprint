package fabric

import (
	"net/http"

	"github.com/blue-jay/blueprint/lib/flight"
	//"github.com/blue-jay/blueprint/middleware/acl"
	"github.com/blue-jay/core/router"
	"io/ioutil"
	"encoding/json"
	"fmt"


	"flag"
	"os"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// namespace
	v1 "k8s.io/api/core/v1"
	"github.com/golang/glog"

)

// Load the routes.
func Load() {
	router.Get("/fabric", Index)
	router.Post("/fabric", PostFabric)
}

// Index displays the home page.
func Index(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	v := c.View.New("fabric/index")
	//if c.Sess.Values["id"] != nil {
	//	v.Vars["first_name"] = c.Sess.Values["first_name"]
	//}

	v.Render(w, r)
}

type Data struct{
	Ext  int
	Type int
}
// Index displays the home page.
func PostFabric(w http.ResponseWriter, r *http.Request) {
	//K8sInit()
	//c := flight.Context(w, r)
	var data Data
	body, err := ioutil.ReadAll(r.Body)
	if err != nil{
		fmt.Print(err)
		return
	}
	json.Unmarshal(body, &data)
	fmt.Print(data)

	//v.Render(w, r)
}


func K8sInit() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	fmt.Println(*kubeconfig)
	glog.V(2).Info(*kubeconfig)
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		glog.Error(err.Error())
		// panic(err.Error())
		return
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Error(err.Error())
		return
		// panic(err.Error())
	}
	// start namespace
	name:= &v1.Namespace{}
	name.APIVersion = "v1"
	name.Kind = "apitest"
	name.ObjectMeta.Name = "test"

	resName, err := clientset.CoreV1().Namespaces().Create(name)
	if err != nil{
		glog.Error(err)
		return
	}
	glog.Info(resName)

	// end namespace
	for {
		pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		namespace := "default"
		pod := "example-xxxxx"
		_, err = clientset.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %s in namespace %s: %v\n",
				pod, namespace, statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
		}

		time.Sleep(10 * time.Second)
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}