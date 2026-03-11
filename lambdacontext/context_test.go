package lambdacontext

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientContextUnmarshalJSON_InvalidJSON(t *testing.T) {
	var cc ClientContext
	err := json.Unmarshal([]byte(`not valid json`), &cc)
	assert.Error(t, err)
	assert.Empty(t, cc.Custom)
}
