package adapter

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHandleMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Fatalf("Unexpected method: expected POST, got %s", r.Method)
			}
			if r.URL.Path != "/api/v1/datapoints" {
				t.Fatalf("Unexpected path: expected %s, got %s", "/api/v1/datapoints", r.URL.Path)
			}
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Error reading body: %s", err)
			}
			if string(b) != expectedBody {
				t.Fatalf("Unexpected request body; expected:\n\n%s\n\ngot:\n\n%s", expectedBody, string(b))
			}
			// KairosDB API returns a 204 No Content on a successful post
			w.WriteHeader(http.StatusNoContent)
		},
	))
	defer server.Client()

	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("Unable to parse server URL %s: %s", server.URL, err)
	}

	client, err := NewClient(&Options{
		KairosDBURL: serverURL.String(),
	})
	if err != nil {
		t.Fatalf("Error creating server: %s", err)
	}

	err = client.HandleMetrics(&prompbMetrics)
	if err != nil {
		t.Fatal(err)
	}
}
