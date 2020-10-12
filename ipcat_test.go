package main

import (
    "testing"
    "net/http/httptest"
    "net/http"
)




func IpcatResponseStatus(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        if r.Method != "GET"{
            t.Errorf("Expected 'GET' request, got '%s'", r.Method)
        }
        if r.URL.EscapedPath() != "/" {
            t.Errorf("Expected request to '/', got '%s'", r.URL.EscapedPath())
        }
    }))
 
    defer ts.Close()
 
}
