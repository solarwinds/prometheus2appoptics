# Prometheus remote storage provider for AppOptics

![CircleCI](https://circleci.com/gh/solarwinds/prometheus2appoptics.svg?style=svg&circle-token=51448f9d74b885c408a0831b4f81134a422f0f5c)

**NOTE**: this is considered **EXPERIMENTAL** and is not yet recommended for production systems.

An implementation of a Prometheus [remote storage adapter](https://github.com/prometheus/prometheus/tree/master/documentation/examples/remote_storage/remote_storage_adapter) for AppOptics.

`prometheus2appoptics` is a web application that handles incoming payloads of Prometheus Sample data and then converts it into AppOptics Measurement semantics and pushes that up to AppOptics' REST API in rate-limit-compliant batches.

**Assumptions:**

* All Prometheus `Labels` can be converted into AppOptics `Measurement Tags`
* Any Prometheus sample w/ a NaN value is worthless and can be discarded
* prometheus2appoptics will handle any difference in throughput between Prometheus' remote storage flush rate and AppOptics's ingestion limit. As of Prometheus 1.7, the storage rate of remote storage sample queues is not configurable, but such options exist internally and the exposure of those options is [scheduled for the 2.0 release](https://github.com/prometheus/prometheus/issues/3095))

## Deployment
Two methods of deployment supported:

1. Deployment of a binary to indvidual system
2. (Recommended) Deployment via Docker container

### Deploying as a Container
`docker run --env ACCESS_TOKEN=<APPOPTICS_TOKEN> --env SEND_STATS=true -p 4567 solarwinds/prometheus2appoptics`

### Configuring Prometheus

To configure Prometheus to send samples to this binary, add the following to your `prometheus.yml`:

```yaml
# Remote write configuration
remote_write:
  - url: "http://<STORAGE_ADAPTER_HOST>:<STORAGE_ADAPTER_PORT>/receive"
```

## Development

Create the bin with

`make build`

prometheus2appoptics supports [several runtime flags](https://github.com/solarwinds/prometheus2appoptics/blob/master/config/config.go#L29-L32) for configuration:

```
--bind-port (the port the HTTP handler will bind to - defaults to 4567)
--send-stats (sends stats to AppOptics if true, to stdout if false - defaults to false)
--access-email (email address associated with API token - defaults to "")
--access-token (API token string - defaults to "")
```

#### Prometheus
* Install Prometheus by downloading the [latest stable release](https://prometheus.io/download)
* Untar the download and put it anywhere you want.
* Open `prometheus.yml` and configure it with running services.
* To just have *some* data, you can [install the "random" RPC process from Prometheus](https://prometheus.io/docs/introduction/getting_started/#starting-up-some-sample-targets) and run several of them at once
* You can also run the [node exporter](https://github.com/prometheus/node_exporter) on your local system for local stats. Remember to set up a target section in the config file.

# Questions/Comments?
Please [open an issue](https://github.com/solarwinds/prometheus2appoptics/issues/new), we'd love to hear from you. As a SolarWinds Innovation Project, this adapter is supported in a best-effort fashion.
