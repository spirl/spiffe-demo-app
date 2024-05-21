# SPIFFE Demo App

SPIFFE Demo application is a simple demo app that uses SPIFFE Workload API.
The app consists of two parts: frontend and backend.
Frontend provides simple functionality to view the content of SPIFFE X509-SVID, JWT-SVID, and SPIFFE Trust Bundle through simple UI.

You can use it to troubleshoot things like SPIFFE federation (because you can see the trust bundle content):

![SPIFFE trust bundle content](img/trust-bundle.png?raw=true "SPIFFE trust bundle content")

Or easily view X509-SVID or JWT-SVID details:

![X509-SVID details](img/x509-svid.png?raw=true "X509-SVID details")

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

That will install the app in the `default` namespace. If you want to install to a different namespace, you need to add `--create-namespace` flag to the above command:

```bash
$ helm install spiffe-demo spiffe-demo/spiffe-demo-app --namespace foo --create-namespace

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
> helm delete spiffe-demo
```

## Example of using with SPIRE Helm Chart

SPIFFE community support SPIRE Helm Chart which provides a way to install SPIRE in k8s cluster for a quick start.
Check for the latest version of SPIRE Helm Chart installation instructions [here](https://artifacthub.io/packages/helm/spiffe/spire#install-instructions).

I'll use [kind](https://kind.sigs.k8s.io/) as an example but you can use any other k8s distribution with a little bit of adjusment (e.g. don't use `port-forwarding`)

### 1. Create kind cluster

```bash
$ kind create cluster --name spire-demo
Creating cluster "spire-demo" ...
 ✓ Ensuring node image (kindest/node:v1.26.3) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-spire-demo"
You can now use your cluster with:

kubectl cluster-info --context kind-spire-demo
```

### 2. Install SPIRE Helm Charts

Install SPIRE:

```bash
helm upgrade --install -n spire-server spire-crds spire-crds --repo https://spiffe.github.io/helm-charts-hardened/ --create-namespace
helm upgrade --install -n spire-server spire spire --repo https://spiffe.github.io/helm-charts-hardened/
```

### 3. Install spiffe-demo-app

[SPIRE Helm Chart](https://github.com/spiffe/helm-charts-hardened/tree/main/charts/spire) by default uses SPIFFE Workload API socket with a name `spire-agent.sock`. However, `spiffe-demo-app` by [default](charts/spiffe-demo-app/README.md) expects the name for a socket `agent.sock`. Therefore we need to change it.

To install SPIFFE Demo App you can use the following command with a few parameters or you may choose to set them in your `values.yaml` for a helm chart.

```bash
helm upgrade --install spiffe-demo spiffe-demo-app --repo https://elinesterov.github.io/spiffe-demo-app -n spiffe-demo --create-namespace --set app.spiffeSocketName=spire-agent.sock --set app.spiffeCSIDriverInjectionEnabled=false --set app.spiffeCSIDriverVolume=true
```

* `app.spiffeSocketName` - is a name of the SPIFFE Workload API socket that is exposed by SPIFFE Agent. By default, SPIRE Helm Chart uses `spire-agent.sock` name for the socket. However, `spiffe-demo-app` expects `agent.sock` name. Therefore we need to change it.

* `app.spiffeCSIDriverInjectionEnabled` - is a flag that enables or disables the SPIRL COntroller Manager that injects SPIFFE CSIDriver Volume and env variable. Since SPIRE Helm Cart doesn't have such a feature we need to disable it.

* `app.spiffeCSIDriverVolume` - flag that enables SPIFFE CSI Driver Volume mount and also set `SPIFFE_ENDPOINT_SOCKET` environment variable to the path of the SPIFFE Workload API socket. This is required for the `spiffe-demo-app` to connect to the SPIFFE Agent. Since we cannot use SPIRL Controller Manager with SPIRE we should manually set these values.

### 4. Connecting to the spiffe-demo-app

Since we use kind cluster in this example, the easiest way to connect to a frontend is by using kubectl port-frowarding feature. If you use minikube, eks or any other flavour of k8s you might be able to use node port or LoadBalancer.

With kind run the following command to enable port-frwarding to the spiffe-demo-app frontend:

```bash
$ kubectl -n spiffe-demo port-forward  svc/spiffe-demo-service 8080:80

Forwarding from 127.0.0.1:8080 -> 8080
Forwarding from [::1]:8080 -> 8080
```

Now you can point your browser to [http://localhost:8080](http://localhost:8080) to connect to the frontend. Clock buttons to see JWT-SVID, X509-SVID and expore SPIFFE Trust Bundle in a very simple way.

### 5. Using busybox for troubleshooting

Sometimes you need to run a few simple shell commands to list mounted agent socket content or for any other reason. We have `busybox` container withtin the spiffe-demo-app deployemnt that is disabled by default. In order to enable it you can add `--set app.enableBusybox=true` when installing `spiffe-demo-app` helm chart.

### 6. Cleaning up

Cleaning up is simple:

```bash
helm -n spiffe-demo delete spiffe-demo
helm -n spire delete spire
kubectl delete namespace spire
kubectl delete namespace spiffe-demo
kind delete cluster --name spire-dem
```

or you can just execute the lates command in case of kind.

## Values

See [charts/spiffe-demo-app/README.md](charts/spiffe-demo-app/README.md)
