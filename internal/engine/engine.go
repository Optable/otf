// Package engine manages the CLI engine binaries that carry out run operations.
package engine

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"net/url"
)

// MinEngineVersion specifies the minimum engine version accepted by OTF.
const MinEngineVersion = "1.2.0"

// MinTofuVersion specifies the minimum version accepted for the tofu engine:
// tofu forked from terraform at 1.6.0 and no earlier release exists.
const MinTofuVersion = "1.6.0"

var (
	// Default is the default for setting the default engine.
	//
	// NOTE: the actual default engine that has been set by the user should be
	// retrieved via the daemon config.
	Default = Terraform()
	// ErrInvalidVersion is returned when a engine version string is
	// not a semantic version string (major.minor.patch).
	ErrInvalidVersion = errors.New("invalid engine version")
	// ErrUnknownEngine is returned when a string does not name a known engine.
	ErrUnknownEngine = errors.New("no engine found with that name: must be either 'terraform' or 'tofu'")
)

func Engines() []*Engine {
	return []*Engine{
		Terraform(),
		Tofu(),
	}
}

// Lookup returns the engine with the given name.
func Lookup(name string) (*Engine, error) {
	var e Engine
	if err := e.set(name); err != nil {
		return nil, err
	}
	return &e, nil
}

// Engine represents a CLI capable of carrying out infrastructure as code
// operations, e.g. terraform.
type Engine struct {
	Name           string
	DefaultVersion string
	client         Client
}

// Client provides access to the engine's upstream services.
type Client interface {
	// getLatestVersion retrieves the latest available (semantic) version.
	getLatestVersion(context.Context) (string, error)
	// sourceURL returns the URL for retrieving a given version of the engine
	// binary.
	sourceURL(version string) *url.URL
}

func (e *Engine) String() string { return e.Name }

// MinVersion is the earliest version of this engine accepted by OTF.
func (e *Engine) MinVersion() string {
	if e.Name == tofuName {
		return MinTofuVersion
	}
	return MinEngineVersion
}

func (e *Engine) Type() string { return "engine" }

func (e *Engine) Set(v string) error {
	return e.set(v)
}

func (e *Engine) MarshalText() ([]byte, error) {
	return []byte(e.Name), nil
}

func (e *Engine) UnmarshalText(text []byte) error {
	return e.set(string(text))
}

func (e *Engine) Scan(text any) error {
	if text == nil {
		return nil
	}
	s, ok := text.(string)
	if !ok {
		return fmt.Errorf("expected database value to be a string: %#v", text)
	}
	return e.set(s)
}

func (e *Engine) Value() (driver.Value, error) {
	if e == nil {
		return nil, nil
	}
	return e.Name, nil
}

func (e *Engine) set(v string) error {
	switch v {
	case terraformName:
		*e = *Terraform()
	case tofuName:
		*e = *Tofu()
	default:
		return fmt.Errorf("%w: %s", ErrUnknownEngine, v)
	}
	return nil
}
