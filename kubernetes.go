package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func ExecuteRemoteCommand() (string, string, error) {

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.CoreV1().Pods("postgrest").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	var command string = "killall -SIGUSR1 postgrest"
	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}

	for _, pod := range pods.Items {
		if pod.Name == "api" {
			continue
		}

		request := clientset.CoreV1().RESTClient().
			Post().
			Namespace("postgrest").
			Resource("pods").
			Name(pod.Name).
			SubResource("exec").Param("container", "process-kill").
			VersionedParams(&v1.PodExecOptions{
				Command: []string{"/bin/sh", "-c", command},
				Stdin:   false,
				Stdout:  true,
				Stderr:  true,
				TTY:     true,
			}, scheme.ParameterCodec)
		exec, err := remotecommand.NewSPDYExecutor(config, "POST", request.URL())
		if err = exec.Stream(remotecommand.StreamOptions{
			Stdout: buf,
			Stderr: errBuf,
		}); err != nil {
			return "", "", fmt.Errorf("%w Failed executing command %s on %v/%v", err, command)
		}
	}
	fmt.Println("Good Kill Process")

	return buf.String(), errBuf.String(), nil

}

func GetPods(c *gin.Context) {

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.CoreV1().Pods("postgrest").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	var list []string
	for _, pod := range pods.Items {
		list = append(list, pod.Name)
	}
	c.IndentedJSON(http.StatusOK, list)

}

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.POST("/", func(c *gin.Context) { ExecuteRemoteCommand() })
	router.GET("/", GetPods)
	router.Run("0.0.0.0:3000")
}
