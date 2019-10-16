package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/util/homedir"
	"kmodules.xyz/client-go/tools/clientcmd"
	"kmodules.xyz/client-go/tools/portforward"
)

var (
	kubeconfigPath = func() string {
		kubecfg := os.Getenv("KUBECONFIG")
		if kubecfg != "" {
			return kubecfg
		}
		return filepath.Join(homedir.HomeDir(), ".kube", "config")
	}()
	kubeContext = ""
)

func main() {
	config, err := clientcmd.BuildConfigFromContext(kubeconfigPath, kubeContext)
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}
	kc := kubernetes.NewForConfigOrDie(config)
	tunnel := portforward.NewTunnel(kc.CoreV1().RESTClient(), config, "default", "nginx", 80)
	defer tunnel.Close()
	err = tunnel.ForwardPort()
	if err != nil {
		log.Fatalln(err)
	}

	url := fmt.Sprintf("http://127.0.0.1:%d", tunnel.Local)
	fmt.Println("url = ", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	b, err := httputil.DumpResponse(resp, true)
	if err == nil {
		fmt.Println(string(b))
	} else {
		fmt.Println(err)
	}
	fmt.Println()
}
