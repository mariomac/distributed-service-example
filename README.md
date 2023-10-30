# distributed-service-example
An example deployment and instrumentation of a distributed service that
is automatically instrumented by Beyla.

Requirements:
* Local Kubernetes (use K3d, as Kind does not work well with Beyla as DaemonSet)

## Run the test environment and services

```
k3d cluster create
make build-all push-all-k3d
kubectl apply -f all-services.yml
```

## Instrument all the services as a Daemonset

First, make a copy of the `grafana-credentials.template.yml` file and
fill the gaps with your grafana endpoint and credentials information:

```
cp grafana-credentials.template.yml grafana-credentials.yml
vi grafana-credentials.yml
# after editing...
kubectl apply -f grafana-credentials.yml
```

Then you can deploy Grafana Beyla along with the grafana agent. Maybe you are interested to provide
a finer-grained configuration by editing the `beyla-config` configmap
that is inside the `beyla-daemonset.yml` file:

```
kubectl apply -f grafana-agent.yml
kubectl apply -f beyla-daemonset.yml
```


