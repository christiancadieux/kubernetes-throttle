package main

import (
	"context"
	"fmt"
	"github.com/christiancadieux/kubernetes-throttle/pkg/client"
	"github.com/christiancadieux/kubernetes-throttle/pkg/throttle"
	"github.com/christiancadieux/kubernetes-throttle/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"os"
	"time"
)

const (
	CGROUP0        = "/sys/fs/cgroup/cpu,cpuacct/kubepods"
	ENV_NAME       = "ENV_NAME"
	SERVER_PORT    = "SERVER_PORT"
	DEFAULT_PORT   = "9191"
	SERVER_LOGGING = "SERVER_LOGGING"
	MY_NODE_NAME   = "MY_NODE_NAME"
)

func main() {
	ctx, cxl := context.WithCancel(context.Background())
	defer cxl()
	logger := logrus.New()

	r := mux.NewRouter()
	port := DEFAULT_PORT
	portEnv := os.Getenv(SERVER_PORT)
	if portEnv != "" {
		port = portEnv
	}

	site := os.Getenv(ENV_NAME)
	if site == "" {
		site = utils.MyNodeName()
	}
	cgroupPath := CGROUP0
	cgroupEnv := os.Getenv("CGROUP_PATH")
	if cgroupEnv != "" {
		cgroupPath = cgroupEnv
	}
	loggingOn := os.Getenv(SERVER_LOGGING) == "Y"

	cgroupVersion := "v1"
	// find the cgroup
	if stat, err := os.Stat(cgroupPath + "/cpu,cpuacct"); err == nil && stat.IsDir() {
		cgroupPath += "/cpu,cpuacct/kubepods"
	} else {
		cgroupPath += "/kubepods" // cgroupv2
		cgroupVersion = "v2"
	}
	logger.Infof("Using cgroupPath=%s ", cgroupPath)

	var err error
	kubeClient, err := client.LoadClient()
	if err != nil {
		fmt.Println("error", err)
		logger.Error(err)
		os.Exit(1)
	}

	nodeName := os.Getenv(MY_NODE_NAME)
	if nodeName == "" {
		logger.Error("env MY_NODE_NAME is not defined")
		os.Exit(1)
	}
	logger.Infof("Server http://%s:%s/metrics", nodeName, port)

	r.HandleFunc("/metrics", prometheusNode(logger, ctx, kubeClient, nodeName, loggingOn, cgroupPath, cgroupVersion))

	http.ListenAndServe(":"+port, r)
}

// called by prometheus every minute
func prometheusNode(logger *logrus.Logger, ctx context.Context, kubeClient *kubernetes.Clientset,
	nodeName string, loggingOn bool, cgroupPath string, cgroupVersion string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		throttle.SendPrometheusNodeData(w, logger, ctx, kubeClient, nodeName, cgroupPath, cgroupVersion)
		if loggingOn {
			lapse := time.Now().Sub(start)
			logger.Infof("%s - %s %s - %v", r.RemoteAddr, r.Method, r.RequestURI, lapse)
		}

	}
}
