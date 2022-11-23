# Debug hcloud reverse-proxy
This repo contains a mutatingwebhook that uses the ENV variable HCLOUD_ENDPOINT to redirect all requests that normally go to api.hetzner.cloud to a reverse-proxy that stores and forwards the request. The code for the webhook is located under `webhook`. A manifest to use the webhook and to redirect requests to a reverse proxy can be found under `k8s/mutatingwebhook.yaml`. To analyze a cluster you can just:
```
kubectl apply -f https://raw.githubusercontent.com/23technologies/debug-hcloud-api/305bb8f8504d87a296acf3a29eed3a334c4a39be/k8s/mutatingwebhook.yaml
```

## Reverse proxy configuration
Furthermore this repository contains files to configure a reverse proxy and to setup prometheus and grafana for it. To create a reverse proxy first setup cert-manager and the letsencrypt-prod clusterissuer. Next install nginx via helm chart and the values located under `k8s/helm/nginx.yaml`. These values install nginx along with a sidecar that reads the accesslog of nginx and exposes them as a prometheus-metric, to make these metrics available in the whole cluster please apply `k8s/reverse-proxy-metrics-service.yaml`. To configure the sidecar properly please apply `k8s/configmap-nginx-accesslog-to-metrics-exporter.yaml. To create the reverse-proxy ingress, please apply `k8s/ingress-mitm-reverse-proxy.yaml`. If you want the remaining tokens as an additional metric also apply `k8s/hcloud-remaining-tokens-metric.yaml`, make sure that you add the token that you want to monitor to the included script. You can scrape your metrics with prometheus and the provided values under `k8s/helm/prometheus.yaml`. You can graph them with grafana with the grafana helmchart and the values under `k8s/helm/grafana.yaml`
