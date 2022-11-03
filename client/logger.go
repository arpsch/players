package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// LoggerRoundTripper represents the logger middleware
type LoggerRoundTripper struct {
	core   http.RoundTripper
	logger io.Writer
}

// NewLogRoundTripper constructs the loging round tripper with
// given logger option
func NewLogRoundTripper(rt http.RoundTripper, l io.Writer) (*LoggerRoundTripper, error) {

	if rt == nil || l == nil {
		return nil, errors.New("failed to set loging round tripper")
	}
	return &LoggerRoundTripper{
		core:   rt,
		logger: l,
	}, nil
}

// RoundTrip adds the request duration with request details
func (l LoggerRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {

	start := time.Now()

	resp, _ := l.core.RoundTrip(r)

	fmt.Fprintf(l.logger, "[request duration: %s] %s %s\n", time.Since(start).String(), r.Method, r.URL.String())

	return resp, nil
}
