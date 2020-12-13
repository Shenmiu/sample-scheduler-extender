# sample-scheduler-extender

A sample to showcase how to create a k8s scheduler extender.

if you cluster is installed by kubeadm，maybe your `kube-scheduler's` yaml like this:
```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    component: kube-scheduler
    tier: control-plane
  name: kube-scheduler
  namespace: kube-system
spec:
  containers:
  - command:
    - kube-scheduler
    - --authentication-kubeconfig=/etc/kubernetes/scheduler.conf
    - --authorization-kubeconfig=/etc/kubernetes/scheduler.conf
    - --bind-address=127.0.0.1
    - --kubeconfig=/etc/kubernetes/scheduler.conf
    - --leader-elect=true
    - --config=/etc/kubernetes/scheduler-config.yaml
    - --v=9
    image: gcr.azk8s.cn/google_containers/kube-scheduler:v1.19.2
    imagePullPolicy: IfNotPresent
    livenessProbe:
      failureThreshold: 8
      httpGet:
        host: 127.0.0.1
        path: /healthz
        port: 10251
        scheme: HTTP
      initialDelaySeconds: 15
      timeoutSeconds: 15
    name: kube-scheduler
    resources:
      requests:
        cpu: 100m
    volumeMounts:
    - mountPath: /etc/kubernetes/scheduler.conf
      name: kubeconfig
      readOnly: true
    - mountPath: /etc/kubernetes/scheduler-config.yaml
      name: scheduler-config
      readOnly: true
  hostNetwork: true
  priorityClassName: system-cluster-critical
  volumes:
  - hostPath:
      path: /etc/kubernetes/scheduler.conf
      type: FileOrCreate
    name: kubeconfig
  - hostPath:
      path: /etc/kubernetes/scheduler-config.yaml
      type: FileOrCreate
    name: scheduler-config
status: {}
```

You must mount `scheduler-extender.yaml` to Pod.

## Notes

- Prioritize webhook won't be triggered if it's running on an one-node cluster. As it makes no sense to run priorities logic when there is only one candidate:

```go
// from k8s.io/kubernetes/pkg/scheduler/core/generic_scheduler.go
func (g *genericScheduler) Schedule(pod *v1.Pod, nodeLister algorithm.NodeLister) (string, error) {
    ...
	// When only one node after predicate, just use it.
	if len(filteredNodes) == 1 {
		metrics.SchedulingAlgorithmPriorityEvaluationDuration.Observe(metrics.SinceInMicroseconds(startPriorityEvalTime))
		return filteredNodes[0].Name, nil
    }
    ...
}
```