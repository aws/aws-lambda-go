// Copyright 2026 Amazon.com, Inc. or its affiliates. All Rights Reserved

package lambda

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// These tests verify both API clients bypass any configured proxy when
// reaching the link-local Runtime/Extensions API. A non-loopback URL is used
// because Go never proxies loopback addresses, which would mask the bug.

func nonLoopbackAPIRequest(t *testing.T) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "http://192.0.2.1:9001/2018-06-01/runtime/invocation/next", nil)
	require.NoError(t, err)
	return req
}

func proxyForClient(t *testing.T, client *http.Client, req *http.Request) *string {
	t.Helper()
	rt := client.Transport
	if rt == nil {
		rt = http.DefaultTransport
	}
	transport, ok := rt.(*http.Transport)
	require.True(t, ok, "expected an *http.Transport to inspect proxy behavior")
	if transport.Proxy == nil {
		return nil
	}
	proxyURL, err := transport.Proxy(req)
	require.NoError(t, err)
	if proxyURL == nil {
		return nil
	}
	s := proxyURL.String()
	return &s
}

func TestRuntimeAPIClientBypassesProxy(t *testing.T) {
	t.Setenv("HTTP_PROXY", "http://192.0.2.100:3128")
	t.Setenv("HTTPS_PROXY", "http://192.0.2.100:3128")
	t.Setenv("http_proxy", "http://192.0.2.100:3128")
	t.Setenv("https_proxy", "http://192.0.2.100:3128")

	client := newRuntimeAPIClient("192.0.2.1:9001").httpClient
	proxy := proxyForClient(t, client, nonLoopbackAPIRequest(t))
	require.Nil(t, proxy, "Runtime API client must not route through the configured proxy, got %v", proxy)
}

func TestExtensionAPIClientBypassesProxy(t *testing.T) {
	t.Setenv("HTTP_PROXY", "http://192.0.2.100:3128")
	t.Setenv("HTTPS_PROXY", "http://192.0.2.100:3128")
	t.Setenv("http_proxy", "http://192.0.2.100:3128")
	t.Setenv("https_proxy", "http://192.0.2.100:3128")

	client := newExtensionAPIClient("192.0.2.1:9001").httpClient
	proxy := proxyForClient(t, client, nonLoopbackAPIRequest(t))
	require.Nil(t, proxy, "Extensions API client must not route through the configured proxy, got %v", proxy)
}

// wrappedRoundTripper is a RoundTripper that is not an *http.Transport, mirroring
// how instrumentation libraries (for example otelhttp) replace http.DefaultTransport.
type wrappedRoundTripper struct {
	base http.RoundTripper
}

func (w *wrappedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return w.base.RoundTrip(req)
}

// When http.DefaultTransport is a non-*http.Transport wrapper, constructing the
// clients must not panic on a type assertion, and the transport must still bypass
// the proxy.
func TestAPIClientsWithWrappedDefaultTransport(t *testing.T) {
	t.Setenv("HTTP_PROXY", "http://192.0.2.100:3128")
	t.Setenv("HTTPS_PROXY", "http://192.0.2.100:3128")
	t.Setenv("http_proxy", "http://192.0.2.100:3128")
	t.Setenv("https_proxy", "http://192.0.2.100:3128")

	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })
	http.DefaultTransport = &wrappedRoundTripper{base: original}

	for _, tc := range []struct {
		name   string
		client func() *http.Client
	}{
		{"runtime", func() *http.Client { return newRuntimeAPIClient("192.0.2.1:9001").httpClient }},
		{"extension", func() *http.Client { return newExtensionAPIClient("192.0.2.1:9001").httpClient }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var client *http.Client
			require.NotPanics(t, func() { client = tc.client() })
			proxy := proxyForClient(t, client, nonLoopbackAPIRequest(t))
			require.Nil(t, proxy, "%s API client must not route through the configured proxy, got %v", tc.name, proxy)
		})
	}
}
