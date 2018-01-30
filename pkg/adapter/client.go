package adapter

import (
	kclient "github.com/ajityagaty/go-kairosdb/client"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/prompb"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("package", "adapter")
var ignoredSamples = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "prometheus_kairosdb_ignored_samples_total",
		Help: "The total number of samples not sent to KairosDB due to unsupported float values (Inf, -Inf, NaN).",
	},
)

// Adapter is an interface for a Remote Storage Adapter
type Adapter interface {
	Start()
	HandleMetrics(request *prompb.WriteRequest) error
}

// Client is a KairosDB Remote Storage Adapter
type Client struct {
	kairosClient kclient.Client
	opts         *Options
}

// Options are options to be passed to the KarisoDB Remote Storage Adapter
type Options struct {
	KairosDBURL string
	ListenPort  int64
	PrintJson   bool
}

func (opts *Options) validate() error {
	if len(opts.KairosDBURL) == 0 {
		return errors.New("kairosDB URL cannot be empty")
	}
	return nil
}

// NewClient returns a new KairosDB Remote Storage Adapter client
func NewClient(opts *Options) (*Client, error) {
	if err := opts.validate(); err != nil {
		return nil, err
	}
	prometheus.MustRegister(ignoredSamples)
	kairosDBClient := kclient.NewHttpClient(opts.KairosDBURL)
	adapterClient := &Client{
		kairosClient: kairosDBClient,
		opts:         opts}
	return adapterClient, nil
}

// Start will start the KairosDB Remote Storage Adapter
func (c *Client) Start() {
	listenAndServe(c, c.opts.ListenPort)
}

// HandleMetrics takes in the prometheus metrics and sends them to KairosDB
func (c *Client) HandleMetrics(req *prompb.WriteRequest) error {
	logger := log.WithField("function", "writeMetrics")
	logger.Debug("Entered writeMetrics")
	mb := BuildKairosDBMetrics(protoToSamples(req))
	resp, err := c.kairosClient.PushMetrics(mb)
	if err != nil {
		logger.Error(errors.Wrap(err, "Error pushing metrics"))
	}
	if resp != nil && len(resp.GetErrors()) != 0 {
		logger.WithFields(logrus.Fields{"statusCode": resp.GetStatusCode(), "errors": resp.GetErrors()}).Error("Error from response")
	}
	return err
}
