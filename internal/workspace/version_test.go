package workspace

import (
	"testing"

	"github.com/leg100/otf/internal/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustVersion(t *testing.T, s string) *Version {
	t.Helper()

	v := &Version{}
	require.NoError(t, v.set(s))
	return v
}

func TestVersion(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantError  error
		wantString string
		wantEngine string
	}{
		{name: "valid semver", input: "1.9.3", wantString: "1.9.3"},
		{name: "track latest version", input: "latest", wantString: "latest"},
		{name: "invalid semver", input: "1,2,0", wantError: engine.ErrInvalidVersion},
		{name: "unsupported terraform version", input: "0.14.0", wantError: ErrUnsupportedTerraformVersion},
		{name: "make exception for go-tfe integration test version", input: "0.10.0", wantString: "0.10.0"},
		{name: "make exception for go-tfe integration test version", input: "0.11.0", wantString: "0.11.0"},
		{name: "make exception for go-tfe integration test version", input: "0.11.1", wantString: "0.11.1"},
		{name: "tofu hint", input: "1.6.0 (tofu)", wantString: "1.6.0", wantEngine: "tofu"},
		{name: "terraform hint", input: "1.5.7 (terraform)", wantString: "1.5.7", wantEngine: "terraform"},
		{name: "hint without whitespace", input: "1.6.0(tofu)", wantString: "1.6.0", wantEngine: "tofu"},
		{name: "hint on latest", input: "latest (tofu)", wantString: "latest", wantEngine: "tofu"},
		{name: "tofu below its minimum version", input: "1.5.7 (tofu)", wantError: ErrUnsupportedTerraformVersion},
		{name: "go-tfe exception does not apply to tofu", input: "0.11.0 (tofu)", wantError: ErrUnsupportedTerraformVersion},
		{name: "unknown engine", input: "1.6.0 (pulumi)", wantError: engine.ErrUnknownEngine},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Version{}
			err := v.set(tt.input)
			if tt.wantError != nil {
				assert.ErrorContains(t, err, tt.wantError.Error())
				return
			}
			require.NoError(t, err)
			// the hint is never round-tripped: it is reported via Engine.
			assert.Equal(t, tt.wantString, v.String())
			if tt.wantEngine == "" {
				assert.Nil(t, v.Engine)
			} else {
				require.NotNil(t, v.Engine)
				assert.Equal(t, tt.wantEngine, v.Engine.String())
			}
		})
	}
}
