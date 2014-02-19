package fake

import (
    "net/http"
    "testing"
)

type ResponseWriter struct {
    http.ResponseWriter

    t       *testing.T
    headers http.Header
    body    []byte
    status  int
}

func New(t *testing.T) *ResponseWriter {
    return &ResponseWriter{
        t:       t,
        headers: make(http.Header),
    }
}

func (r *ResponseWriter) Header() http.Header {
    return r.headers
}

func (r *ResponseWriter) Write(body []byte) (int, error) {
    r.body = body
    return len(body), nil
}

func (r *ResponseWriter) WriteHeader(status int) {
    r.status = status
}

func (r *ResponseWriter) Assert(status int, body []byte) {
  if r.status != status {
    r.t.Errorf("expected status %+v to equal %+v", r.status, status)
  }
  if string(r.body) != string(body) {
    r.t.Errorf("expected body %+v to equal %+v", string(r.body), body)
  }
}
