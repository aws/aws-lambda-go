package lambdacontext

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientContextUnmarshalJSON(t *testing.T) {
	t.Run("non-string custom values are serialized to string", func(t *testing.T) {
		input := `{
			"Client": {"installation_id": "install1"},
			"custom": {
				"key1": "stringval",
				"key2": {"nested": "object"},
				"key3": 42
			}
		}`
		var cc ClientContext
		err := json.Unmarshal([]byte(input), &cc)
		require.NoError(t, err)
		assert.Equal(t, "install1", cc.Client.InstallationID)
		assert.Equal(t, "stringval", cc.Custom["key1"])
		assert.JSONEq(t, `{"nested":"object"}`, cc.Custom["key2"])
		assert.Equal(t, "42", cc.Custom["key3"])
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		var cc ClientContext
		err := json.Unmarshal([]byte(`not valid json`), &cc)
		assert.Error(t, err)
	})
}
