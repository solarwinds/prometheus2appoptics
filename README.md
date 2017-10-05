# Librato & AppOptics as Prometheus remote storage provider
An implementation of a Prometheus [remote storage adapter](/prometheus/prometheus/tree/master/documentation/examples/remote_storage/remote_storage_adapter) for Librato and AppOptics.

`p2l` is a web application that handles incoming payloads of Prometheus Sample data and then converts it into Librato Measurement semantics and pushes that up to Librato's REST API in rate-limit-compliant batches.

**Assumptions:**

* All Prometheus `Labels` can be converted into Librato `Measurement Tags`
* Any Prometheus sample w/ a NaN value is worthless and can be discarded
* p2l will handle any difference in throughput between Prometheus' remote storage flush rate and Librato's ingestion limit. As of Prometheus 1.7, the storage rate of remote storage sample queues is not configurable, but such options exist internally and the exposure of those options is [scheduled for the 2.0 release](https://github.com/prometheus/prometheus/issues/3095))

## Deployment
Two methods of deployment supported:

1. Deployment of a binary to indvidual system
2. (Recommended) Deployment via Docker container

### Deploying as a Container
`docker run -p 4567 solarwinds/prom2swi-cloud`

### Configuring Prometheus

To configure Prometheus to send samples to this binary, add the following to your `prometheus.yml`:

```yaml
# Remote write configuration
remote_write:
  - url: "http://<STORAGE_ADAPTER_HOST>:<STORAGE_ADAPTER_PORT>/receive"
```

## Development

#### dep
[dep](https://github.com/golang/dep) is the new official dependency tool for Go. It's still in the prototype phase, but it's totally usable. You can `brew install dep` on macOS or you can build yourself on any system that can run Go:

`go get -u github.com/golang/dep/cmd/dep`


#### p2l
Assuming you have a standard Go environment with a checkout of the p2l code in the normal place and the `dep` tool in your `$PATH`, you can install the project's dependencies with:

`dep ensure`

Then create the bin with

`make`

p2l supports [several runtime flags](https://github.com/solarwinds/p2l/blob/master/config/config.go#L18-L21) for configuration:

```
--bind-port (the port the HTTP handler will bind to - defaults to 4567)
--send-stats (sends stats to Librato if true, to stdout if false - defaults to false)
--access-email (email address associated with API token - defaults to "")
--access-token (API token string - defaults to "")
```

#### Prometheus
* Install Prometheus by downloading the [latest stable release](https://github.com/prometheus/prometheus/releases/tag/v1.7.2)
* Untar the download and put it anywhere you want
* Open `prometheus.yml` and configure it with running services.
* To just have *some* data, you can [install the "random" RPC process from Prometheus](https://prometheus.io/docs/introduction/getting_started/#starting-up-some-sample-targets) and run several of them at once
* You can also run the [node exporter](https://github.com/prometheus/node_exporter) on your local system for local stats. Remember to set up a target section in the config file.