# Librato & AppOptics as Prometheus remote storage provider

[![CircleCI](https://circleci.com/gh/solarwinds/p2l.svg?style=svg&circle-token=de9c33d8cfa8724aadc105c798d57dca9060dc81)](https://circleci.com/gh/solarwinds/p2l)

An implementation of a Prometheus [remote storage adapter](/prometheus/prometheus/tree/master/documentation/examples/remote_storage/remote_storage_adapter) for Librato and AppOptics.
# Deployment
Two methods of deployment supported:
1. Deployment of a binary to indvidual system
1. (Recommended) Deployment via Docker container

## Deploying as a Container
```docker run -p 4567 solarwinds/prom2swi-cloud```

## Configuring Prometheus

To configure Prometheus to send samples to this binary, add the following to your `prometheus.yml`:

```yaml
# Remote write configuration
remote_write:
  - url: "http://localhost:4567/receive"
