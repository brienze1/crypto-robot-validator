package mocks

import (
	"github.com/brienze1/crypto-robot-validator/pkg/log"
	"net/http"
	"net/http/httptest"
)

type httpClient struct {
	DoCounter      int
	DoError        error
	Server         *httptest.Server
	StatusCode     int
	ServerResponse string
}

func HttpClient() *httpClient {
	return &httpClient{
		StatusCode: 200,
	}
}

func (h *httpClient) Do(req *http.Request) (*http.Response, error) {
	h.DoCounter++

	if h.DoError != nil {
		return nil, h.DoError
	}

	realClient := http.Client{}
	return realClient.Do(req)
}

func (h *httpClient) SetupServer() {
	h.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(h.StatusCode)
		log.Logger().Info("GET called", h.ServerResponse)
		_, _ = w.Write([]byte(h.ServerResponse))
	}))
}

func (h *httpClient) GetUrl() string {
	return h.Server.URL
}

func (h *httpClient) Close() {
	h.Server.Close()
}

func (h *httpClient) Reset() {
	h.Server = nil
	h.DoCounter = 0
	h.DoError = nil
	h.StatusCode = 200
}
