package workspace

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"slices"

	"github.com/leg100/otf/internal/engine"
	"github.com/leg100/otf/internal/semver"
)

var apiTestTerraformVersions = []string{"0.10.0", "0.11.0", "0.11.1"}

// Version is a workspace's engine version.
type Version struct {
	// Latest if true means runs use the Latest available version at time of
	// creation of the run.
	Latest bool
	// Engine is the engine selected via a hint in the version string, e.g.
	// "1.6.0 (tofu)". Nil when the string carries no hint, in which case the
	// engine is left to the caller.
	Engine *engine.Engine
	// semver is the semantic version of the engine; must be non-empty if latest
	// is false.
	//
	// TODO: use custom type
	semver string
}

// String returns the bare version, without any engine hint: the engine is
// persisted and reported separately.
func (v *Version) String() string {
	if v.Latest {
		return "latest"
	}
	return v.semver
}

func (v *Version) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v *Version) UnmarshalText(text []byte) error {
	return v.set(string(text))
}

func (v *Version) Scan(text any) error {
	if text == nil {
		return nil
	}
	s, ok := text.(string)
	if !ok {
		return fmt.Errorf("expected database value to be a string: %#v", text)
	}
	return v.set(s)
}

func (v *Version) Value() (driver.Value, error) {
	if v == nil {
		return nil, nil
	}
	return v.String(), nil
}

var errEmptyString = errors.New("value cannot be an empty string")

// checkVersion rejects a version below the minimum published by the engine. A
// nil engine falls back to the lowest minimum across all engines, deferring the
// engine-specific check to whoever resolves the engine.
func checkVersion(e *engine.Engine, version string) error {
	minVersion := engine.MinEngineVersion
	if e != nil {
		minVersion = e.MinVersion()
	}
	if semver.Compare(version, minVersion) >= 0 {
		return nil
	}
	// NOTE: we make an exception for the specific versions posted by the go-tfe
	// integration tests, which only ever use terraform.
	if minVersion == engine.MinEngineVersion && slices.Contains(apiTestTerraformVersions, version) {
		return nil
	}
	return fmt.Errorf("%w: minimum version is %s", ErrUnsupportedTerraformVersion, minVersion)
}

func (v *Version) set(value string) error {
	value, hinted, err := engine.ParseVersion(value)
	if err != nil {
		return err
	}
	v.Engine = hinted

	switch value {
	case "latest":
		v.Latest = true
	case "":
		return errEmptyString
	default:
		if !semver.IsValid(value) {
			return engine.ErrInvalidVersion
		}
		if err := checkVersion(hinted, value); err != nil {
			return err
		}
		v.semver = value
	}
	return nil
}
