package clap_test

import (
	"testing"

	"github.com/fred1268/go-clap/clap"
)

func TestNonEmptyTrailing(t *testing.T) {
	t.Parallel()
	type config struct {
		Field    string   `clap:"--string"`
		Trailing []string `clap:"trailing,-t"`
	}
	cfg := &config{}
	var err error
	var results *clap.Results
	if results, err = clap.Parse([]string{"--string", "hello"}, cfg); err == nil {
		t.Errorf("unexpected trailing: %s", err)
	}
	t.Logf("t: %v\n", results)
}

func TestNonStringTrailing(t *testing.T) {
	t.Parallel()
	type config struct {
		Field    string `clap:"--string"`
		Trailing int    `clap:"trailing"`
	}
	cfg := &config{}
	var err error
	var results *clap.Results
	if results, err = clap.Parse([]string{"--string", "hello"}, cfg); err == nil {
		t.Errorf("unexpected trailing: %s", err)
	}
	t.Logf("t: %v\n", results)
}

func TestInvalidShortName(t *testing.T) {
	t.Parallel()
	type config struct {
		Field string `clap:",-s,foo,bar"`
	}
	cfg := &config{}
	var err error
	var results *clap.Results
	if results, err = clap.Parse([]string{"-s", "hello"}, cfg); err == nil {
		t.Errorf("unexpected valid shortname: %s", err)
	}
	t.Logf("t: %v\n", results)
}

func TestNoShortName(t *testing.T) {
	t.Parallel()
	type config struct {
		Field string `clap:",-"`
	}
	cfg := &config{}
	var err error
	var results *clap.Results
	if results, err = clap.Parse([]string{"-s", "hello"}, cfg); err == nil {
		t.Errorf("unexpected valid shortname: %s", err)
	}
	t.Logf("t: %v\n", results)
}

func TestInvalidShortNameParameter(t *testing.T) {
	t.Parallel()
	type config struct {
		Field string `clap:",-s,unexpected"`
	}
	cfg := &config{}
	var err error
	var results *clap.Results
	if results, err = clap.Parse([]string{"-s", "hello"}, cfg); err == nil {
		t.Errorf("unexpected valid shortname: %s", err)
	}
	t.Logf("t: %v\n", results)
}
