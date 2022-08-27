# Git HTTP Server operator
This operator it's a controller for 
[GIT Http Server](https://github.com/jarpsimoes/git-http-server), 
providing support to simplify deployments inside complex clusters. 
With this operator will be able to deploy HTML5 apps with "NoOps" paradigm.

## Description
Git Http Server Operator is written Go with [operator-sdk](https://sdk.operatorframework.io/):

- **operator-sdk version:** v1.22.2 
- **commit:** da3346113a8a75e11225f586482934000504a60f
- **kubernetes version:** v1.24.1
- **go version:** go1.18.4

Should be compatible with kubernetes version: Kubernetes +1.18*

When GitHttpServer was deployed will be provided inside cluster:
- Deployment with GitHttpServer
- Service to exposed HTTP_PORT

## Getting Started

### Requirements:
- Should be able to connect with target cluster via kubectl (command line)
- Should have permissions to deploy operators inside cluster

### 1. Install Operator inside cluster
Download manifests and apply on cluster:
```shell
$ curl https://raw.githubusercontent.com/jarpsimoes/git-http-server-operator/main/dist/git-http-server-operator.yaml | kubectl apply -f -
```
Check if operator is running:
```shell
$ kubectl get nodes -n operator-system

NAME                                   READY    STATUS    RESTARTS   AGE
operator-controller-manager-####-####   2/2     Running   0          2m5s
```

### 2. Deploy simple git http server
Create operator descriptor file as sample.yaml:
```yaml
apiVersion: jarpsimoes.github.io/v1alpha1
kind: GitHttpServer
metadata:
  name: githttpserver-sample
spec:
  repo-url: https://github.com/jarpsimoes/html_sample.git
```

Deploy GitHttpServer:
```shell
$ kubectl apply -f sample.yaml
```

Check installed components:

**Deployment:**
```shell
$ kubectl decribe deployment githttpserver-sample
Name:                   githttpserver-sample-deployment
Namespace:              default
CreationTimestamp:      Sat, 27 Aug 2022 02:24:14 +0100
Labels:                 <none>
Annotations:            deployment.kubernetes.io/revision: 1
Selector:               app=githttpserver-sample,operator=git-http-server-operator,tier=backend
Replicas:               1 desired | 1 updated | 1 total | 1 available | 0 unavailable
StrategyType:           RollingUpdate
MinReadySeconds:        0
RollingUpdateStrategy:  25% max unavailable, 25% max surge
Pod Template:
  Labels:  app=githttpserver-sample
           operator=git-http-server-operator
           tier=backend
  Containers:
   githttpserver-sample-pod:
    Image:      jarpsimoes/git_http_server:latest
    Port:       8081/TCP
    Host Port:  0/TCP
    Liveness:   http-get http://:8081/_health delay=0s timeout=1s period=10s #success=1 #failure=3
    Startup:    http-get http://:8081/_health delay=0s timeout=1s period=10s #success=1 #failure=3
    Environment:
      PATH_CLONE:          _clone
      PATH_PULL:           _pull
      PATH_VERSION:        _version
      PATH_WEBHOOK:        _hook
      PATH_HEALTH:         _health
      REPO_BRANCH:         main
      REPO_TARGET_FOLDER:  target-git
      REPO_URL:            https://github.com/jarpsimoes/jarpsimoes.github.io.git
      HTTP_PORT:           8081
    Mounts:                <none>
  Volumes:                 <none>
Conditions:
  Type           Status  Reason
  ----           ------  ------
  Available      True    MinimumReplicasAvailable
  Progressing    True    NewReplicaSetAvailable
OldReplicaSets:  <none>
NewReplicaSet:   githttpserver-sample-deployment-c5669cc57 (1/1 replicas created)
Events:
  Type    Reason             Age   From                   Message
  ----    ------             ----  ----                   -------
  Normal  ScalingReplicaSet  7m2s  deployment-controller  Scaled up replica set githttpserver-sample-deployment-c5669cc57 to 1
```
**Service:**
```shell
$ kubectl describe service githttpserver-sample-service
Name:              githttpserver-sample-service
Namespace:         default
Labels:            <none>
Annotations:       <none>
Selector:          app=githttpserver-sample,operator=git-http-server-operator,tier=backend
Type:              ClusterIP
IP Family Policy:  SingleStack
IP Families:       IPv4
IP:                10.3.32.138
IPs:               10.3.32.138
Port:              <unset>  80/TCP
TargetPort:        8081/TCP
Endpoints:         10.1.0.34:8081
Session Affinity:  None
Events:            <none>

```

**Note:** TO-DO - Configurations available

**Note:** TO-DO - Complex implementations

