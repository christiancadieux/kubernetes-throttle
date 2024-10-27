package main

import (
	"context"
	"fmt"
	"github.com/christiancadieux/kubernetes-throttle/clog"
	"github.com/christiancadieux/kubernetes-throttle/pkg/client"
	"github.com/christiancadieux/kubernetes-throttle/pkg/throttle"
	"github.com/christiancadieux/kubernetes-throttle/pkg/utils"
	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"os"
	"time"
)

const (
	CGROUP0           = "/sys/fs/cgroup/cpu,cpuacct/kubepods"
	ENV_NAME          = "ENV_NAME"
	SERVER_PORT       = "SERVER_PORT"
	PROMETHEUS_PORT   = "PROMETHEUS_PORT"
	DEFAULT_PORT      = "9191"
	DEFAULT_PROM_PORT = "10099"
	SERVER_LOGGING    = "SERVER_LOGGING"
	MY_NODE_NAME      = "MY_NODE_NAME"
)

func main() {
	ctx, cxl := context.WithCancel(context.Background())
	defer cxl()
	logger := clog.MakeSplunkLogger(clog.Fields{"s": "rdei-throttle"})

	r := mux.NewRouter()
	port := DEFAULT_PORT
	portEnv := os.Getenv(SERVER_PORT)
	if portEnv != "" {
		port = portEnv
	}
	promPort := DEFAULT_PROM_PORT
	promPortEnv := os.Getenv(PROMETHEUS_PORT)
	if promPortEnv != "" {
		promPort = promPortEnv
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

	var err error
	kubeClient, err := client.LoadClient()
	if err != nil {
		fmt.Println("error", err)
		logger.Error(err)
		os.Exit(1)
	}

	if os.Getenv("TEST") == "Y" {
		fmt.Println("Calling throttle.test()")
		err := throttle.Test(logger, kubeClient, cgroupPath)
		if err != nil {
			logger.Error(err)
		}
		os.Exit(0)
	}

	nodeName := os.Getenv(MY_NODE_NAME)
	if nodeName == "" {
		logger.Error("env MY_NODE_NAME is not defined")
		os.Exit(1)
	}
	logger.Infof("Server http://%s:%s/metrics", nodeName, port)

	r.HandleFunc("/metrics", prometheusNode(logger, ctx, kubeClient, nodeName, loggingOn, cgroupPath))

	http.ListenAndServe(":"+port, r)
}

// called by prometheus every minute
func prometheusNode(logger *clog.Logger, ctx context.Context, kubeClient *kubernetes.Clientset,
	nodeName string, loggingOn bool, cgroupPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		throttle.SendPrometheusNodeData(w, logger, ctx, kubeClient, nodeName, cgroupPath)
		if loggingOn {
			lapse := time.Now().Sub(start)
			logger.Infof("%s - %s %s - %v", r.RemoteAddr, r.Method, r.RequestURI, lapse)
		}

	}
}
