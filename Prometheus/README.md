

## Prometheus queries

```
http://localhost:9090/graph

current-context: anvil-stage-rdei-dev-01
namespace:  cc-test
```


```
rate(namespace_cpu_nr_throttled[5m])
rate(pod_cpu_nr_throttled_count{namespace="cc-test"}[5m])
sum(rate(pod_cpu_nr_throttled_count{namespace="cc-test"}[5m]))

# percentage of throttling
rate(pod_cpu_nr_throttled_count{namespace="cc-test"}[2m]) / (rate(pod_cpu_nr_throttled_count{namespace="cc-test"}[2m])+rate(pod_cpu_nr_periods_count{namespace="cc-test"}[2m])) * 100
```

### stress

```
while test 1=1; do    stress -c 3 -t  300;    sleep 300; done
```

![3mins-cycle](stress-01.png)


### REF

https://jaanhio.me/blog/kubernetes-cpu-requests-limits/
