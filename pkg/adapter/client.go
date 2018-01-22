// Copyright © 2018 Kenneth Herner <kherner@navistone.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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
	HandleMetrics(request *prompb.WriteRequest)
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
func (c *Client) HandleMetrics(req *prompb.WriteRequest) {
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
}
