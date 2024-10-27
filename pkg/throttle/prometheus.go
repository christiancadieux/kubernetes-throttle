package throttle

import (
	"bufio"
	"context"
	"fmt"
	"github.com/comcast-ccp-containers/clog"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	CPUPATH = "/sys/fs/cgroup/cpu,cpuacct"
)

type PodMetrics struct {
	Namespace      string
	Name           string
	nr_periods     int64
	nr_throttled   int64
	throttled_time int64
	limit          int64
	usage          int64
}

// SendPrometheusNodeData - Send throttle and df data
func SendPrometheusNodeData(w http.ResponseWriter, logger *clog.Logger, ctx context.Context,
	kubeClient *kubernetes.Clientset, nodeName string, cgroupPath string) error {

	pods, err := kubeClient.CoreV1().Pods("").List(context.Background(), v1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})
	if err != nil {
		return err
	}

	podMetrics := []PodMetrics{}
	nsTotals := map[string]*PodMetrics{}

	// per-node metrics
	usageCpuUser, usageCpuSys := readNodeCpu(logger)
	if usageCpuUser != nil {
		fmt.Fprintf(w, "# HELP node_cpu_usage_count CPU accountUsage User per Cpu\n")
		fmt.Fprintf(w, "# TYPE node_cpu_usage_count counter\n")
		for ix, user := range usageCpuUser {
			if user > 0 {
				fmt.Fprintf(w, `node_cpu_usage_count{type="user", cpu="%03d"} %d`+"\n", ix, user)
			}
		}
		for ix, sys := range usageCpuSys {
			if sys > 0 {
				fmt.Fprintf(w, `node_cpu_usage_count{type="sys", cpu="%03d"} %d`+"\n", ix, sys)
			}
		}
	}

	// per pod/namespace metrics
	for _, pod := range pods.Items {
		nr_periods, nr_throttled, throttled_time, limit, usage, found := readCgroup(logger, pod, cgroupPath)
		if !found {
			continue
		}
		podMetrics = append(podMetrics, PodMetrics{pod.Namespace, pod.Name, nr_periods,
			nr_throttled, throttled_time, limit, usage})
		if _, ok := nsTotals[pod.Namespace]; !ok {
			nsTotals[pod.Namespace] = &PodMetrics{}
		}
		nsTotals[pod.Namespace].nr_throttled += nr_throttled
		nsTotals[pod.Namespace].nr_periods += nr_periods
		nsTotals[pod.Namespace].throttled_time += throttled_time
		nsTotals[pod.Namespace].limit += limit
		nsTotals[pod.Namespace].usage += usage

	}

	// PER POD

	fmt.Fprintf(w, "# HELP pod_cpu_nr_throttled_count CPU Throttled periods per Pod\n")
	fmt.Fprintf(w, "# TYPE pod_cpu_nr_throttled_count counter\n")
	for _, m := range podMetrics {
		if m.nr_periods > 0 {
			fmt.Fprintf(w, `pod_cpu_nr_throttled_count{namespace="%s", pod="%s"} %d`+"\n", m.Namespace, m.Name, m.nr_throttled)
		}
	}
	fmt.Fprintf(w, "# HELP pod_cpu_nr_periods_count CPU periods per Pod\n")
	fmt.Fprintf(w, "# TYPE pod_cpu_nr_periods_count counter\n")
	for _, m := range podMetrics {
		if m.nr_periods > 0 {
			fmt.Fprintf(w, `pod_cpu_nr_periods_count{namespace="%s", pod="%s"} %d`+"\n", m.Namespace, m.Name, m.nr_periods)
		}
	}
	fmt.Fprintf(w, "# HELP pod_cpu_throttled_time_count CPU throttled time per Pod\n")
	fmt.Fprintf(w, "# TYPE pod_cpu_throttled_time_count counter\n")
	for _, m := range podMetrics {
		if m.nr_periods > 0 {
			fmt.Fprintf(w, `pod_cpu_throttled_time_count{namespace="%s", pod="%s"} %d`+"\n", m.Namespace, m.Name, m.throttled_time)
		}
	}

	fmt.Fprintf(w, "# HELP pod_cpu_limit CPU limit Pod\n")
	fmt.Fprintf(w, "# TYPE pod_cpu_limit gauge\n")
	for _, m := range podMetrics {
		fmt.Fprintf(w, `pod_cpu_limit{namespace="%s", pod="%s"} %d`+"\n", m.Namespace, m.Name, m.limit)
	}
	fmt.Fprintf(w, "# HELP pod_cpu_usage_count CPU accountUsage Pod\n")
	fmt.Fprintf(w, "# TYPE pod_cpu_usage_count counter\n")
	for _, m := range podMetrics {
		fmt.Fprintf(w, `pod_cpu_usage_count{namespace="%s", pod="%s"} %d`+"\n", m.Namespace, m.Name, m.usage)
	}

	// PER NAMESPACE

	keys := []string{}
	for ns, _ := range nsTotals {
		keys = append(keys, ns)
	}
	sort.Strings(keys)

	fmt.Fprintf(w, "# HELP namespace_cpu_nr_throttled_count CPU Throttled periods per Namespace\n")
	fmt.Fprintf(w, "# TYPE namespace_cpu_nr_throttled_count counter\n")
	for _, ns := range keys {
		m := nsTotals[ns]
		if m.nr_periods > 0 {
			fmt.Fprintf(w, `namespace_cpu_nr_throttled_count{namespace="%s"} %d`+"\n", ns, m.nr_throttled)
		}
	}
	fmt.Fprintf(w, "# HELP namespace_cpu_nr_periods_count CPU periods per Namespace from containers\n")
	fmt.Fprintf(w, "# TYPE namespace_cpu_nr_periods_count counter\n")
	for _, ns := range keys {
		m := nsTotals[ns]
		if m.nr_periods > 0 {
			fmt.Fprintf(w, `namespace_cpu_nr_periods_count{namespace="%s"} %d`+"\n", ns, m.nr_periods)
		}
	}
	fmt.Fprintf(w, "# HELP namespace_cpu_throttled_time_count CPU throttled time per Namespace from containers\n")
	fmt.Fprintf(w, "# TYPE namespace_cpu_throttled_time_count counter\n")
	for _, ns := range keys {
		m := nsTotals[ns]
		if m.nr_periods > 0 {
			fmt.Fprintf(w, `namespace_cpu_throttled_time_count{namespace="%s"} %d`+"\n", ns, m.throttled_time)
		}
	}

	fmt.Fprintf(w, "# HELP namespace_cpu_limit CPU limit Namespace\n")
	fmt.Fprintf(w, "# TYPE namespace_cpu_limit gauge\n")
	for _, ns := range keys {
		m := nsTotals[ns]
		fmt.Fprintf(w, `namespace_cpu_limit{namespace="%s"} %d`+"\n", ns, m.limit)
	}
	fmt.Fprintf(w, "# HELP namespace_cpu_usage_count CPU accountUsage per Namespace\n")
	fmt.Fprintf(w, "# TYPE namespace_cpu_usage_count counter\n")
	for _, ns := range keys {
		m := nsTotals[ns]
		fmt.Fprintf(w, `namespace_cpu_usage_count{namespace="%s"} %d`+"\n", ns, m.usage)
	}

	return nil
}

func readNodeCpu(logger *clog.Logger) ([]int64, []int64) {

	// PER CPU
	usageAllCpu, err := readLines(CPUPATH + "/cpuacct.usage_all")
	if err != nil {
		logger.Errorf("readline cpuacct.usage_all error", err)
		return nil, nil
	}
	size := 0
	for _, line := range usageAllCpu {
		words := strings.Split(line, " ")
		if len(words) != 3 || words[0] == "cpu" {
			continue
		}
		size++
	}
	usageCpuUser := make([]int64, size)
	usageCpuSys := make([]int64, size)

	for _, line := range usageAllCpu {
		words := strings.Split(line, " ")
		if len(words) != 3 || words[0] == "cpu" {
			continue
		}
		cpuno, err := strconv.ParseInt(words[0], 10, 64)
		cpuUser, err := strconv.ParseInt(words[1], 10, 64)
		if err != nil {
			logger.Errorf("readline cpuNo error", err)
			continue
		}
		cpuSys, err := strconv.ParseInt(words[2], 10, 64)
		usageCpuUser[cpuno] = cpuUser
		usageCpuSys[cpuno] = cpuSys
	}
	return usageCpuUser, usageCpuSys
}

func readCgroup(logger *clog.Logger, pod corev1.Pod, cgroupPath string) (int64, int64, int64, int64, int64, bool) {
	var nr_periods int64 = 0
	var nr_throttled int64 = 0
	var throttled_time int64 = 0
	var limit int64 = 0
	var usage int64 = 0

	var groups = []string{"burstable", "besteffort", "guaranteed"}

	found := false
	for _, g := range groups {
		path := fmt.Sprintf("%s/%s/pod%s", cgroupPath, g, pod.UID)
		_, err := os.Open(path)
		if err != nil {
			continue
		}
		cpustatLines, err := readLines(path + "/cpu.stat")
		if err != nil {
			logger.Errorf("readline cpu.stat error", err)
			continue
		}

		found = true
		for _, line := range cpustatLines {
			words := strings.Split(line, " ")
			i, err := strconv.ParseInt(words[1], 10, 64)
			// i, err := strconv.Atoi(words[1])
			if err != nil {
				logger.Errorf("Atoi error", err)
				break
			}
			if words[0] == "nr_periods" {
				nr_periods = i

			} else if words[0] == "throttled_time" {
				throttled_time = i

			} else if words[0] == "nr_throttled" {
				nr_throttled = i
			}
		}
		// read /cpu.cfs_quota_us = k8s resources.limit
		quotaValue, err := os.ReadFile(path + "/cpu.cfs_quota_us")
		if err != nil || len(quotaValue) < 2 {
			logger.Errorf("readLines /cpu.cfs_quota_us error - %v", err)
			continue
		}
		i, err := strconv.ParseInt(string(quotaValue[:len(quotaValue)-1]), 10, 64) // remove newline
		if err != nil {
			logger.Errorf("ParseInt error %s - %v", string(quotaValue), err)
			continue
		}
		if i > 0 {
			limit = i
		}
		// read /cpuacct.usage
		usageValue, err := os.ReadFile(path + "/cpuacct.usage")
		if err != nil || len(usageValue) < 2 {
			logger.Errorf("readLines cpuacct.usage error - %v", err)
			continue
		}
		i, err = strconv.ParseInt(string(usageValue[:len(usageValue)-1]), 10, 64) // remove newline
		if err != nil {
			logger.Errorf("ParseInt error %s - %v", string(usageValue), err)
			continue
		}
		if i > 0 {
			usage = i
		}

		break
	}
	return nr_periods, nr_throttled, throttled_time, limit, usage, found
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
