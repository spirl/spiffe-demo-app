# SPIFFE Demo App

SPIFFE Demo application is a simple demo app that uses SPIFFE Workload API.
The app consists of two parts: frontend and backend.
Frontend provides simple functionality to view the content of SPIFFE X509-SVID, JWT-SVID, and SPIFFE Trust Bundle through simple UI.

## Prerequisites

You'll need `kubectl` and `helm` installed. Follow the official documentation for both projects on how to install them:

* [Helm](https://helm.sh/docs/intro/install/)
* [kubectl](https://kubernetes.io/docs/tasks/tools/)

## Usage

### Add SPIFFE Demo Helm Chart

```bash
helm repo add spiffe-demo https://elinesterov.github.io/spiffe-demo-app
```

### Install the chart

```bash
$ helm install spiffe-demo spiffe-demo/spiffe-demo-app

NAME: spiffe-demo
LAST DEPLOYED: Wed May 24 10:34:58 2023
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

That will install the app in the `default` namespace. If you want to install to a different namespace, you need to create it first via `kubectl create namespace foo' and then add `--namespace` flag to the above command:

```bash
$ helm install spiffe-demo spiffe-demo/spiffe-demo-app --namespace foo

NAME: spiffe-demo
LAST DEPLOYED: Wed May 24 10:39:28 2023
NAMESPACE: foo
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

### Connect to the frontend

The easiest way to use [kubernetes port forwarding](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster/):

```bash
$ kubectl port-forward  svc/spiffe-demo-service 8080:80

Forwarding from 127.0.0.1:8080 -> 8080
Forwarding from [::1]:8080 -> 8080
```

If you installed a specific namespace - don't forget about `--namespace flag.

Now you can point your browser to [http://localhost:8080](http://localhost:8080) to reach a frontend.

### Uninstall

```bash
helm delete spiffe-demo
```

## Values

See [charts/spiffe-demo-app/README.md](charts/spiffe-demo-app/README.md)
