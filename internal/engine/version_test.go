package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantVersion string
		wantEngine  string
		wantError   bool
	}{
		{name: "no hint", input: "1.5.7", wantVersion: "1.5.7"},
		{name: "tofu hint", input: "1.6.0 (tofu)", wantVersion: "1.6.0", wantEngine: "tofu"},
		{name: "terraform hint", input: "1.5.7 (terraform)", wantVersion: "1.5.7", wantEngine: "terraform"},
		{name: "extra whitespace", input: " 1.6.0  ( tofu ) ", wantVersion: "1.6.0", wantEngine: "tofu"},
		{name: "no whitespace", input: "1.6.0(tofu)", wantVersion: "1.6.0", wantEngine: "tofu"},
		{name: "latest", input: "latest (tofu)", wantVersion: "latest", wantEngine: "tofu"},
		{name: "unknown engine", input: "1.6.0 (pulumi)", wantError: true},
		{name: "empty hint", input: "1.6.0 ()", wantVersion: "1.6.0 ()"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, e, err := ParseVersion(tt.input)
			if tt.wantError {
				assert.ErrorIs(t, err, ErrUnknownEngine)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantVersion, version)
			if tt.wantEngine == "" {
				assert.Nil(t, e)
			} else {
				require.NotNil(t, e)
				assert.Equal(t, tt.wantEngine, e.Name)
			}
		})
	}
}

func TestMinVersion(t *testing.T) {
	assert.Equal(t, MinEngineVersion, Terraform().MinVersion())
	assert.Equal(t, MinTofuVersion, Tofu().MinVersion())
}
