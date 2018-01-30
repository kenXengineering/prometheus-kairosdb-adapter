package adapter

import (
	"io/ioutil"
	"net/http"

	"strconv"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/prompb"
)

// listenAndServe will listen on the given port and send metrics to the
// adapters HandleMetrics function
func listenAndServe(adapter Adapter, port int64) {
	logger := log.WithField("function", "listenAndServe")

	http.HandleFunc("/receive", func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Got Message")
		compressed, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error(errors.Wrap(err, "Error reading http body"))
			return
		}

		reqBuf, err := snappy.Decode(nil, compressed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Error(errors.Wrap(err, "Error decoding http body"))
			return
		}

		var req prompb.WriteRequest
		if err := proto.Unmarshal(reqBuf, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Error(errors.Wrap(err, "Error unmarshaling request"))
			return
		}

		// Send the metrics to kairosDB
		logger.Debug("Sending to Handler")
		adapter.HandleMetrics(&req)
	})
	url := ":" + strconv.FormatInt(port, 10)
	logger.WithField("url", url).Info("HTTP Server Starting")
	http.ListenAndServe(url, nil)
}
