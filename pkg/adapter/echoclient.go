package adapter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/prompb"
)

// EchoClient is a simple client that echos back the received metrics
type EchoClient struct {
	opts *Options
}

// NewEchoClient returns a new EchoClient
func NewEchoClient(opts *Options) *EchoClient {
	prometheus.MustRegister(ignoredSamples)
	return &EchoClient{
		opts: opts,
	}
}

// Start will start the EchoClient
func (c *EchoClient) Start() {
	log.Info("Starting Echo Client")
	listenAndServe(c, c.opts.ListenPort)
}

// HandleMetrics will take in the prometheus metrics and print out the formatted metrics
// or the metrics in JSON that KairosDB expects if the --json flag is specified
func (c *EchoClient) HandleMetrics(req *prompb.WriteRequest) error {
	log.Debug("Got Metric")
	if c.opts.PrintJson {
		c.echoJson(req)
		return nil
	}
	c.echoFormatted(req)
	return nil
}

func (c *EchoClient) echoFormatted(req *prompb.WriteRequest) {
	log.Debug("Building Formatted")
	for _, ts := range req.Timeseries {
		m := make(model.Metric, len(ts.Labels))
		for _, l := range ts.Labels {
			m[model.LabelName(l.Name)] = model.LabelValue(l.Value)
		}
		fmt.Println(m)

		for _, s := range ts.Samples {
			fmt.Printf("  %f %d\n", s.Value, s.Timestamp)
		}
	}
}

func (c *EchoClient) echoJson(req *prompb.WriteRequest) {
	log.Debug("Building Json")
	mb := BuildKairosDBMetrics(protoToSamples(req))
	data, err := mb.Build()
	if err != nil {
		log.Error(err)
	}
	fmt.Println(string(data))
}
