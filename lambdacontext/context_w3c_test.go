// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package lambdacontext

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLambdaContextW3C(t *testing.T) {
	traceparent := "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"

	cases := []struct {
		name          string
		clientContext string
		expected      map[string]string
	}{
		{
			name:          "empty when no client context",
			clientContext: "",
			expected:      map[string]string{},
		},
		{
			name:          "empty when client context has no w3c key",
			clientContext: `{"custom":{"value":"test"}}`,
			expected:      map[string]string{},
		},
		{
			name:          "baggage only",
			clientContext: `{"w3c":{"baggage":"userId=alice"}}`,
			expected:      map[string]string{"baggage": "userId=alice"},
		},
		{
			name:          "all three allowlisted fields",
			clientContext: `{"custom":{"value":"test"},"w3c":{"traceparent":"` + traceparent + `","tracestate":"rojo=00f067aa0ba902b7","baggage":"userId=alice"}}`,
			expected: map[string]string{
				"traceparent": traceparent,
				"tracestate":  "rojo=00f067aa0ba902b7",
				"baggage":     "userId=alice",
			},
		},
		{
			name:          "allowlist drops non-allowlisted keys",
			clientContext: `{"w3c":{"baggage":"keep=me","unknownField":"nope","x-custom-trace":"nope"}}`,
			expected:      map[string]string{"baggage": "keep=me"},
		},
		{
			name:          "drops allowlisted fields with non-string values",
			clientContext: `{"w3c":{"traceparent":42,"tracestate":null,"baggage":{"nested":"no"}}}`,
			expected:      map[string]string{},
		},
		{
			name:          "non-object w3c treated as empty",
			clientContext: `{"w3c":"not-an-object"}`,
			expected:      map[string]string{},
		},
		{
			name:          "array w3c treated as empty",
			clientContext: `{"w3c":["baggage=abc"]}`,
			expected:      map[string]string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			lc := &LambdaContext{}
			lc.ExtractW3C([]byte(tc.clientContext))
			assert.Equal(t, tc.expected, lc.W3C())
		})
	}
}

func TestLambdaContextW3CReturnsFreshCopy(t *testing.T) {
	lc := &LambdaContext{}
	lc.ExtractW3C([]byte(`{"w3c":{"baggage":"abc"}}`))

	// Mutating a returned map must not affect the context's internal state.
	first := lc.W3C()
	first["baggage"] = "tampered"
	first["injected"] = "nope"

	assert.Equal(t, map[string]string{"baggage": "abc"}, lc.W3C())
}

func TestLambdaContextW3CNeverNil(t *testing.T) {
	// A zero-value context (ExtractW3C never called) still yields a usable,
	// non-nil empty map.
	lc := &LambdaContext{}
	assert.Equal(t, map[string]string{}, lc.W3C())
}

func TestW3CAllowedFieldsIsImmutable(t *testing.T) {
	assert.Equal(t, []string{"traceparent", "tracestate", "baggage"}, W3CAllowedFields())
	got := W3CAllowedFields()
	got[0] = "tampered"
	got = append(got, "injected")
	assert.Equal(t, []string{"traceparent", "tracestate", "baggage"}, W3CAllowedFields())
}
