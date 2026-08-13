# spiffe-demo-app

![Version: 0.5.0](https://img.shields.io/badge/Version-0.5.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square)

A Helm chart to install spiffe-demo-app

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| app | object | `{"enableBusybox":false,"enableDebug":false,"healthCheckOnly":false,"spiffeCSIDriverInjectionEnabled":true,"spiffeCSIDriverVolume":false,"spiffeSocketMountPathDir":"/spiffe-workload-api","spiffeSocketName":"agent.sock"}` | The App configuration |
| app.enableBusybox | bool | `false` | Enable busybox container |
| app.enableDebug | bool | `false` | Enable spirldbg |
| app.healthCheckOnly | bool | `false` | Enable health check only mode |
| app.spiffeCSIDriverInjectionEnabled | bool | `true` | SPIRL SPIFFE CSI injeciton enabled |
| app.spiffeCSIDriverVolume | bool | `false` | Add SPIFFE CSIdirver volume |
| app.spiffeSocketMountPathDir | string | `"/spiffe-workload-api"` | SPIFFE Workload API socket mount path |
| app.spiffeSocketName | string | `"agent.sock"` | SPIFFE Workload API socket name exposed by the agent the resulting default socket path will be /spiffe-workload-api/agent.sock spiffeSocketName is a filename from the socket path for the agent |
| busyboxImage | string | `"busybox:1.38.0"` | Image used for busybox when enabled. Pinned to a release tag; ":latest" is mutable. Bump deliberately when updating. |
| image.pullPolicy | string | `"IfNotPresent"` | The image pull policy |
| image.registry | string | `"ghcr.io"` | The OCI registry to pull the image from |
| image.repository | string | `"spirl/spiffe-demo-app"` | The repository within the registry |
| image.tag | string | `"v0.3.3"` | The image tag to pull |
| podSecurity | object | `{"restricted":false}` | PodSecurity hardening. OFF by default: the rendered pod is unchanged from previous chart versions. Set podSecurity.restricted=true to apply the securityContext blocks below and make the app pod compliant with the PodSecurity "restricted" profile (required for namespaces labeled pod-security.kubernetes.io/enforce: restricted, e.g. GKE Autopilot). |
| podSecurity.restricted | bool | `false` | Apply the restricted-compliant securityContext blocks below. |
| podSecurityContext | object | `{"runAsNonRoot":true,"runAsUser":65532,"seccompProfile":{"type":"RuntimeDefault"}}` | Pod-level security context, applied only when podSecurity.restricted=true. The app image runs as the nonroot UID 65532; runAsUser is pinned to match it. WARNING: this whole map is a single override target. Setting it to {} (or a partial map) replaces the entire securityContext and drops "restricted" compliance. "restricted" only requires runAsNonRoot, not a specific UID, so if you override image.tag with an image using a different nonroot UID, update runAsUser to match (or drop it and rely on the image's own USER). |
| securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"readOnlyRootFilesystem":true}` | Container-level security context for the app and busybox containers, applied only when podSecurity.restricted=true. WARNING: this whole map is a single override target. Setting it to {} (or a partial map) drops "restricted" compliance. |
| service | object | `{"port":80,"type":"LoadBalancer"}` | The service type to use |
| spiffeCSIDriver | object | `{"enabled":false,"image":{"registry":"ghcr.io","repository":"spiffe/spiffe-csi-driver","tag":"0.2.3"},"nodeDriverRegistrarImage":{"registry":"registry.k8s.io","repository":"sig-storage/csi-node-driver-registrar","tag":"v2.6.0"}}` | SPIFFE CSI driver support WARNING: when enabled this renders a DaemonSet that is privileged, uses hostPath mounts with Bidirectional mount propagation, and has no pod securityContext. CSI node drivers genuinely require this privilege, so the DaemonSet CANNOT satisfy PodSecurity "restricted" and will be rejected by PSA admission (and by GKE Autopilot outright). Do not enable it in a namespace labeled pod-security.kubernetes.io/enforce: restricted. Run it in a separate privileged/baseline namespace, or use the managed SPIRL injection path. The app Deployment itself (default render, with this disabled) is restricted-compliant. |
| spiffeCSIDriver.enabled | bool | `false` | Enable/disable SPIFFE CSI driver support |
| spiffeCSIDriver.image | object | `{"registry":"ghcr.io","repository":"spiffe/spiffe-csi-driver","tag":"0.2.3"}` | The SPIFFE CSI driver image. Third-party and upgraded independently of this repo's own image, so it is tracked by Renovate. |
| spiffeCSIDriver.nodeDriverRegistrarImage | object | `{"registry":"registry.k8s.io","repository":"sig-storage/csi-node-driver-registrar","tag":"v2.6.0"}` | The node-driver-registrar sidecar image, likewise third-party. |
| spirldbg.image | object | `{"pullPolicy":"IfNotPresent","registry":"ghcr.io","repository":"spirl/spirldbg","tag":"v0.0.18"}` | The spirldbg debug sidecar image. Nested under an "image" key so Renovate's helm-values manager tracks it; the manager only recognises keys whose name ends in "image". |
| spirldbg.image.tag | string | `"v0.0.18"` | Pinned to a released tag. ":latest" is mutable and resolves to an old build carrying CVEs. Bump deliberately. |
| spirldbgSecurityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"runAsNonRoot":true,"runAsUser":65532}` | Container-level security context for the spirldbg debug sidecar, applied only when podSecurity.restricted=true. The spirldbg image runs as root, so runAsNonRoot/runAsUser are set here to override the pod-level context and keep the pod "restricted"-compliant. The root filesystem is left writable so the debug shell can create temp files. NOTE: spirldbg's native tooling may need privileges; the chart's default debug usage (an exec-in shell) works as nonroot. |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.14.2](https://github.com/norwoodj/helm-docs/releases/v1.14.2)
