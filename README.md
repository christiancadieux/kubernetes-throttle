
# Node Throttling Metrics Service


Runs on each node and returns data from the pseudo-files in /sys/fs/cgroup/cpu,cpuacct/kubepods/.

This service runs as a Kubernets daemonset on each node to extract cpu activity and throttling. See /Kubernetes for the daemonset spec.

Generate throttling percentage graph using Prometheus:

![graph]( doc/graph1.png)


## Prometheus Format

The endpoints /metrics returns prometheus-compatible metrics for each node called.


## Prometheus /metrics example

```
$ curl node1-ip:9750/metrics
# HELP pod_cpu_nr_throttled_count CPU Throttled periods per Pod
# TYPE pod_cpu_nr_throttled_count counter
pod_cpu_nr_throttled_count{namespace="cc-test", pod="rdei-tools-7f7b4f97b7-m222r"} 25399
pod_cpu_nr_throttled_count{namespace="kube-system", pod="calicoctl-pldms"} 0
pod_cpu_nr_throttled_count{namespace="kube-system", pod="coredns-f9f669486-rp8ss"} 0
pod_cpu_nr_throttled_count{namespace="kube-system", pod="node-local-dns-krzks"} 0
pod_cpu_nr_throttled_count{namespace="kube-system", pod="prometheus-kube-state-metrics-6dcf587fb8-qwpj2"} 1
pod_cpu_nr_throttled_count{namespace="kube-system", pod="prometheus-node-exporter-7jl6l"} 20889
pod_cpu_nr_throttled_count{namespace="kube-system", pod="prometheus-server-569655877c-698gp"} 57
pod_cpu_nr_throttled_count{namespace="kube-system", pod="rdei-throttle-qvt89"} 0
pod_cpu_nr_throttled_count{namespace="packetsink", pod="packetsink-server-588d7f955f-59c76"} 0
pod_cpu_nr_throttled_count{namespace="packetsink", pod="packetsink-server-588d7f955f-pk7p7"} 0
pod_cpu_nr_throttled_count{namespace="rdei-pod-monitor", pod="cniviptools-green-6f998b66d5-xw2hz"} 0
pod_cpu_nr_throttled_count{namespace="rdei-pod-monitor", pod="kube-ztest-ws-green"} 0
pod_cpu_nr_throttled_count{namespace="rdei-system", pod="relay-0"} 0
pod_cpu_nr_throttled_count{namespace="rdei-system", pod="unicorns-green-754d8f4fb9-fxhft"} 0
pod_cpu_nr_throttled_count{namespace="thirutest", pod="grafana-test-57985ffdcc-k9sj9"} 0
pod_cpu_nr_throttled_count{namespace="thirutest", pod="grafana-test-57985ffdcc-nkbjv"} 1
# HELP pod_cpu_nr_periods_count CPU periods per Pod
# TYPE pod_cpu_nr_periods_count counter
pod_cpu_nr_periods_count{namespace="cc-test", pod="rdei-tools-7f7b4f97b7-m222r"} 30417
pod_cpu_nr_periods_count{namespace="kube-system", pod="calicoctl-pldms"} 96
pod_cpu_nr_periods_count{namespace="kube-system", pod="coredns-f9f669486-rp8ss"} 339846
pod_cpu_nr_periods_count{namespace="kube-system", pod="node-local-dns-krzks"} 403663
pod_cpu_nr_periods_count{namespace="kube-system", pod="prometheus-kube-state-metrics-6dcf587fb8-qwpj2"} 45010
pod_cpu_nr_periods_count{namespace="kube-system", pod="prometheus-node-exporter-7jl6l"} 33166
pod_cpu_nr_periods_count{namespace="kube-system", pod="prometheus-server-569655877c-698gp"} 1457720
pod_cpu_nr_periods_count{namespace="kube-system", pod="rdei-throttle-qvt89"} 16822
pod_cpu_nr_periods_count{namespace="packetsink", pod="packetsink-server-588d7f955f-59c76"} 1859
pod_cpu_nr_periods_count{namespace="packetsink", pod="packetsink-server-588d7f955f-pk7p7"} 239
pod_cpu_nr_periods_count{namespace="rdei-pod-monitor", pod="cniviptools-green-6f998b66d5-xw2hz"} 9155
pod_cpu_nr_periods_count{namespace="rdei-pod-monitor", pod="kube-ztest-ws-green"} 10618
pod_cpu_nr_periods_count{namespace="rdei-system", pod="relay-0"} 6410
pod_cpu_nr_periods_count{namespace="rdei-system", pod="unicorns-green-754d8f4fb9-fxhft"} 8762
pod_cpu_nr_periods_count{namespace="thirutest", pod="grafana-test-57985ffdcc-k9sj9"} 88271
pod_cpu_nr_periods_count{namespace="thirutest", pod="grafana-test-57985ffdcc-nkbjv"} 83162
# HELP pod_cpu_throttled_time_count CPU throttled time per Pod
# TYPE pod_cpu_throttled_time_count counter
pod_cpu_throttled_time_count{namespace="cc-test", pod="rdei-tools-7f7b4f97b7-m222r"} 453939788654
pod_cpu_throttled_time_count{namespace="kube-system", pod="calicoctl-pldms"} 0
pod_cpu_throttled_time_count{namespace="kube-system", pod="coredns-f9f669486-rp8ss"} 65650677
pod_cpu_throttled_time_count{namespace="kube-system", pod="node-local-dns-krzks"} 0
pod_cpu_throttled_time_count{namespace="kube-system", pod="prometheus-kube-state-metrics-6dcf587fb8-qwpj2"} 442019620
pod_cpu_throttled_time_count{namespace="kube-system", pod="prometheus-node-exporter-7jl6l"} 32594600726574
pod_cpu_throttled_time_count{namespace="kube-system", pod="prometheus-server-569655877c-698gp"} 17302879631
pod_cpu_throttled_time_count{namespace="kube-system", pod="rdei-throttle-qvt89"} 0
pod_cpu_throttled_time_count{namespace="packetsink", pod="packetsink-server-588d7f955f-59c76"} 0
pod_cpu_throttled_time_count{namespace="packetsink", pod="packetsink-server-588d7f955f-pk7p7"} 0
pod_cpu_throttled_time_count{namespace="rdei-pod-monitor", pod="cniviptools-green-6f998b66d5-xw2hz"} 0
pod_cpu_throttled_time_count{namespace="rdei-pod-monitor", pod="kube-ztest-ws-green"} 0
pod_cpu_throttled_time_count{namespace="rdei-system", pod="relay-0"} 0
pod_cpu_throttled_time_count{namespace="rdei-system", pod="unicorns-green-754d8f4fb9-fxhft"} 0
pod_cpu_throttled_time_count{namespace="thirutest", pod="grafana-test-57985ffdcc-k9sj9"} 0
pod_cpu_throttled_time_count{namespace="thirutest", pod="grafana-test-57985ffdcc-nkbjv"} 37566065
# HELP pod_cpu_limit CPU limit Pod
# TYPE pod_cpu_limit gauge
pod_cpu_limit{namespace="cc-test", pod="rdei-tools-7f7b4f97b7-m222r"} 190000
pod_cpu_limit{namespace="kube-system", pod="calico-node-9t6fn"} 0
pod_cpu_limit{namespace="kube-system", pod="calicoctl-pldms"} 200000
pod_cpu_limit{namespace="kube-system", pod="coredns-f9f669486-rp8ss"} 100000
pod_cpu_limit{namespace="kube-system", pod="node-local-dns-krzks"} 400000
pod_cpu_limit{namespace="kube-system", pod="prometheus-kube-state-metrics-6dcf587fb8-qwpj2"} 100000
pod_cpu_limit{namespace="kube-system", pod="prometheus-node-exporter-7jl6l"} 70000
pod_cpu_limit{namespace="kube-system", pod="prometheus-server-569655877c-698gp"} 350000
pod_cpu_limit{namespace="kube-system", pod="rdei-throttle-qvt89"} 100000
pod_cpu_limit{namespace="kube-system", pod="sumatra-daemonset-bfcxw"} 0
pod_cpu_limit{namespace="packetsink", pod="packetsink-server-588d7f955f-59c76"} 400000
pod_cpu_limit{namespace="packetsink", pod="packetsink-server-588d7f955f-pk7p7"} 400000
pod_cpu_limit{namespace="rdei-pod-monitor", pod="cniviptools-green-6f998b66d5-xw2hz"} 100000
pod_cpu_limit{namespace="rdei-pod-monitor", pod="kube-ztest-ws-green"} 100000
pod_cpu_limit{namespace="rdei-system", pod="relay-0"} 100000
pod_cpu_limit{namespace="rdei-system", pod="unicorns-green-754d8f4fb9-fxhft"} 50000
pod_cpu_limit{namespace="remove-2", pod="rdei-tools-755665f54b-9z67f"} 0
pod_cpu_limit{namespace="thirutest", pod="grafana-test-57985ffdcc-k9sj9"} 100000
pod_cpu_limit{namespace="thirutest", pod="grafana-test-57985ffdcc-nkbjv"} 100000
# HELP pod_cpu_usage_count CPU accountUsage Pod
# TYPE pod_cpu_usage_count counter
pod_cpu_usage_count{namespace="cc-test", pod="rdei-tools-7f7b4f97b7-m222r"} 5287686829731
pod_cpu_usage_count{namespace="kube-system", pod="calico-node-9t6fn"} 4862826389980
pod_cpu_usage_count{namespace="kube-system", pod="calicoctl-pldms"} 186241469
pod_cpu_usage_count{namespace="kube-system", pod="coredns-f9f669486-rp8ss"} 247442494147
pod_cpu_usage_count{namespace="kube-system", pod="node-local-dns-krzks"} 479810311147
pod_cpu_usage_count{namespace="kube-system", pod="prometheus-kube-state-metrics-6dcf587fb8-qwpj2"} 47492602482
pod_cpu_usage_count{namespace="kube-system", pod="prometheus-node-exporter-7jl6l"} 1695860483875
pod_cpu_usage_count{namespace="kube-system", pod="prometheus-server-569655877c-698gp"} 5799780155924
pod_cpu_usage_count{namespace="kube-system", pod="rdei-throttle-qvt89"} 58473370117
pod_cpu_usage_count{namespace="kube-system", pod="sumatra-daemonset-bfcxw"} 1148600271227
pod_cpu_usage_count{namespace="packetsink", pod="packetsink-server-588d7f955f-59c76"} 22151805081
pod_cpu_usage_count{namespace="packetsink", pod="packetsink-server-588d7f955f-pk7p7"} 1268603205
pod_cpu_usage_count{namespace="rdei-pod-monitor", pod="cniviptools-green-6f998b66d5-xw2hz"} 4902624301
pod_cpu_usage_count{namespace="rdei-pod-monitor", pod="kube-ztest-ws-green"} 5772520485
pod_cpu_usage_count{namespace="rdei-system", pod="relay-0"} 13293519649
pod_cpu_usage_count{namespace="rdei-system", pod="unicorns-green-754d8f4fb9-fxhft"} 4677023997
pod_cpu_usage_count{namespace="remove-2", pod="rdei-tools-755665f54b-9z67f"} 32589937
pod_cpu_usage_count{namespace="thirutest", pod="grafana-test-57985ffdcc-k9sj9"} 75150075488
pod_cpu_usage_count{namespace="thirutest", pod="grafana-test-57985ffdcc-nkbjv"} 73092534502
# HELP namespace_cpu_nr_throttled_count CPU Throttled periods per Namespace
# TYPE namespace_cpu_nr_throttled_count counter
namespace_cpu_nr_throttled_count{namespace="cc-test"} 25399
namespace_cpu_nr_throttled_count{namespace="kube-system"} 20947
namespace_cpu_nr_throttled_count{namespace="packetsink"} 0
namespace_cpu_nr_throttled_count{namespace="rdei-pod-monitor"} 0
namespace_cpu_nr_throttled_count{namespace="rdei-system"} 0
namespace_cpu_nr_throttled_count{namespace="thirutest"} 1
# HELP namespace_cpu_nr_periods_count CPU periods per Namespace from containers
# TYPE namespace_cpu_nr_periods_count counter
namespace_cpu_nr_periods_count{namespace="cc-test"} 30417
namespace_cpu_nr_periods_count{namespace="kube-system"} 2296323
namespace_cpu_nr_periods_count{namespace="packetsink"} 2098
namespace_cpu_nr_periods_count{namespace="rdei-pod-monitor"} 19773
namespace_cpu_nr_periods_count{namespace="rdei-system"} 15172
namespace_cpu_nr_periods_count{namespace="thirutest"} 171433
# HELP namespace_cpu_throttled_time_count CPU throttled time per Namespace from containers
# TYPE namespace_cpu_throttled_time_count counter
namespace_cpu_throttled_time_count{namespace="cc-test"} 453939788654
namespace_cpu_throttled_time_count{namespace="kube-system"} 32612411276502
namespace_cpu_throttled_time_count{namespace="packetsink"} 0
namespace_cpu_throttled_time_count{namespace="rdei-pod-monitor"} 0
namespace_cpu_throttled_time_count{namespace="rdei-system"} 0
namespace_cpu_throttled_time_count{namespace="thirutest"} 37566065
# HELP namespace_cpu_limit CPU limit Namespace
# TYPE namespace_cpu_limit gauge
namespace_cpu_limit{namespace="cc-test"} 190000
namespace_cpu_limit{namespace="kube-system"} 1320000
namespace_cpu_limit{namespace="packetsink"} 800000
namespace_cpu_limit{namespace="rdei-pod-monitor"} 200000
namespace_cpu_limit{namespace="rdei-system"} 150000
namespace_cpu_limit{namespace="remove-2"} 0
namespace_cpu_limit{namespace="thirutest"} 200000
# HELP namespace_cpu_usage_count CPU accountUsage per Namespace
# TYPE namespace_cpu_usage_count counter
namespace_cpu_usage_count{namespace="cc-test"} 5287686829731
namespace_cpu_usage_count{namespace="kube-system"} 14340472320368
namespace_cpu_usage_count{namespace="packetsink"} 23420408286
namespace_cpu_usage_count{namespace="rdei-pod-monitor"} 10675144786
namespace_cpu_usage_count{namespace="rdei-system"} 17970543646
namespace_cpu_usage_count{namespace="remove-2"} 32589937
namespace_cpu_usage_count{namespace="thirutest"} 148242609990

```
