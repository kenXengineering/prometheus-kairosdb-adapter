package adapter

import (
	"math"
	"testing"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/prompb"
)

var (
	samples = model.Samples{
		{
			Metric: model.Metric{
				model.MetricNameLabel: "testmetric",
				"test_label":          "test_label_value1",
			},
			Timestamp: model.Time(123456789123),
			Value:     1.23,
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "testmetric",
				"test_label":          "test_label_value2",
			},
			Timestamp: model.Time(123456789123),
			Value:     5.1234,
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "nan_value",
			},
			Timestamp: model.Time(123456789123),
			Value:     model.SampleValue(math.NaN()),
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "pos_inf_value",
			},
			Timestamp: model.Time(123456789123),
			Value:     model.SampleValue(math.Inf(1)),
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "neg_inf_value",
			},
			Timestamp: model.Time(123456789123),
			Value:     model.SampleValue(math.Inf(-1)),
		},
	}
	expectedBody    = `[{"name":"testmetric","tags":{"test_label":"test_label_value1"},"datapoints":[[123456789123,1.23]]},{"name":"testmetric","tags":{"test_label":"test_label_value2"},"datapoints":[[123456789123,5.1234]]}]`
	numMetrics      = 5
	numValidMetrics = 2

	prompbMetrics = prompb.WriteRequest{
		Timeseries: []*prompb.TimeSeries{
			{
				Labels: []*prompb.Label{
					{
						Name:  model.MetricNameLabel,
						Value: "testmetric",
					},
					{
						Name:  "test_label",
						Value: "test_label_value1",
					},
				},
				Samples: []*prompb.Sample{
					{
						Timestamp: 123456789123,
						Value:     1.23,
					},
				},
			},
			{
				Labels: []*prompb.Label{
					{
						Name:  model.MetricNameLabel,
						Value: "testmetric",
					},
					{
						Name:  "test_label",
						Value: "test_label_value2",
					},
				},
				Samples: []*prompb.Sample{
					{
						Timestamp: 123456789123,
						Value:     5.1234,
					},
				},
			},
			{
				Labels: []*prompb.Label{
					{
						Name:  model.MetricNameLabel,
						Value: "nan_value",
					},
				},
				Samples: []*prompb.Sample{
					{
						Timestamp: 123456789123,
						Value:     math.NaN(),
					},
				},
			},
			{
				Labels: []*prompb.Label{
					{
						Name:  model.MetricNameLabel,
						Value: "pos_inf_value",
					},
				},
				Samples: []*prompb.Sample{
					{
						Timestamp: 123456789123,
						Value:     math.Inf(1),
					},
				},
			},
			{
				Labels: []*prompb.Label{
					{
						Name:  model.MetricNameLabel,
						Value: "neg_inf_value",
					},
				},
				Samples: []*prompb.Sample{
					{
						Timestamp: 123456789123,
						Value:     math.Inf(-1),
					},
				},
			},
		},
	}
)

func TestBuildKairosDBMetrics(t *testing.T) {
	builder := BuildKairosDBMetrics(samples)
	if len(builder.GetMetrics()) != numValidMetrics {
		t.Fatalf("Error building metrics, expected %d metrics, got %d", numValidMetrics, len(builder.GetMetrics()))
	}
}
