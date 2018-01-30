package adapter

import (
	"math"
	"time"

	"github.com/ajityagaty/go-kairosdb/builder"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/prompb"
	"github.com/sirupsen/logrus"
)

// BuildKairosDBMetrics takes in prometheus samples and returns a KairosDB MetricBuilder
func BuildKairosDBMetrics(samples model.Samples) builder.MetricBuilder {
	logger := log.WithField("function", "BuildKairosDBMetric")
	logger.Debug("Building Metrics!")
	mb := builder.NewMetricBuilder()
	for _, s := range samples {
		v := float64(s.Value)
		if math.IsNaN(v) || math.IsInf(v, 0) {
			logger.WithFields(logrus.Fields{"value": v, "sample": s}).Debug("cannot send to KairosDB, skipping sample")
			ignoredSamples.Inc()
			continue
		}
		metric := mb.AddMetric(string(s.Metric[model.MetricNameLabel]))
		// KairosDB timestamps are in milliseconds
		metric.AddDataPoint(s.Timestamp.UnixNano()/int64(time.Millisecond), v)
		tags := tagsFromMetric(s.Metric)
		for name, value := range tags {
			// KairosDB does not like tags with empty values
			if len(value) != 0 {
				metric.AddTag(name, value)
			}
		}
	}
	return mb
}

// protoToSamples turns the prometheus protobuf values into Sample objects
func protoToSamples(req *prompb.WriteRequest) model.Samples {
	var samples model.Samples
	for _, ts := range req.Timeseries {
		metric := make(model.Metric, len(ts.Labels))
		for _, l := range ts.Labels {
			metric[model.LabelName(l.Name)] = model.LabelValue(l.Value)
		}

		for _, s := range ts.Samples {
			samples = append(samples, &model.Sample{
				Metric:    metric,
				Value:     model.SampleValue(s.Value),
				Timestamp: model.Time(s.Timestamp),
			})
		}
	}
	return samples
}

// tagsFromMetric extracts KairosDB tags from a Prometheus metric.
func tagsFromMetric(m model.Metric) map[string]string {
	tags := make(map[string]string, len(m)-1)
	for l, v := range m {
		if l != model.MetricNameLabel {
			tags[string(l)] = string(v)
		}
	}
	return tags
}
