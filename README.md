# Prometheus KairosDB Adapter

[![Build Status](https://travis-ci.org/chosenken/prometheus-kairosdb-adapter.svg?branch=master)](https://travis-ci.org/chosenken/prometheus-kairosdb-adapter)

This is a write adapter that receives samples via Prometheus remote write protocol and stores them in KairosDB.  Any labels attached to the metric will be added to the KairosDB metric as a Tag.

As of now this adapter only supports writing to KairosDB.

Building
---

```bash
make
```

Usage
---
```
Prometheus write adapter for KairosDB

Usage:
  prometheus-kairosdb-adapter [flags]
  prometheus-kairosdb-adapter [command]

Available Commands:
  echo        Prints out received metrics from prometheus
  help        Help about any command

Flags:
  -d, --debug                 Enable Debug
  -h, --help                  help for prometheus-kairosdb-adapter
      --kairosdb-url string   KairosDB URL
  -p, --listen-port int       Listen Port (default 9201)
```

Running
---
```bash
./prometheus_kairosdb_adapter --kairosdb-url http://localhost:8080
```

Configure Prometheus
---
To configure Prometheus to send samples to the KairosDB adapter, add the following to your `prometheus.yaml` file:

```yaml
remote_write:
  - url: "http://localhost:9201/write"
```

Metrics
---

Prometheus metrics are exported on the `/metrics` path.  The standard prometheus client library metrics are provided, along with the metric `prometheus_kairosdb_ignored_samples_total`.
The metric is the total number of samples not sent to KairosDB due to unsupported float values (Inf, -Inf, NaN)."