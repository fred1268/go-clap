package clap_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/fred1268/go-clap/clap"
)

func TestEmpty(t *testing.T) {
	type config struct{}
	cfg := &config{}
	var err error
	var results *clap.Results
	if results, err = clap.Parse([]string{}, cfg); err != nil {
		t.Errorf("parsing error: %s", err)
	}
	t.Logf("t: %v\n", results)
	wanted := &config{}
	if !reflect.DeepEqual(cfg, wanted) {
		t.Errorf("wanted: '%v', got '%v'", wanted, cfg)
	}
}

func TestEmptyAndUselessFields(t *testing.T) {
	type config struct {
		AField       int
		AnotherField string
	}
	cfg := &config{}
	var err error
	var results *clap.Results
	if results, err = clap.Parse([]string{}, cfg); err != nil {
		t.Errorf("parsing error: %s", err)
	}
	t.Logf("t: %v\n", results)
	wanted := &config{}
	if !reflect.DeepEqual(cfg, wanted) {
		t.Errorf("wanted: '%v', got '%v'", wanted, cfg)
	}
}

func TestInvalidType(t *testing.T) {
	type config struct {
		Image int `clap:"--image"`
	}
	cfg := &config{}
	var err error
	var results *clap.Results
	if results, err = clap.Parse([]string{"--image", "photo.png"}, cfg); err != nil {
		if !errors.Is(err, clap.ErrUnexpectedArgument) {
			t.Errorf("parsing error: %s", err)
			return
		}
		t.Logf("t: %v\n", err)
	}
	t.Logf("t: %v\n", results)
}

func TestMandatoryNotFound(t *testing.T) {
	type config struct {
		Image int `clap:"--image,mandatory"`
	}
	cfg := &config{}
	var err error
	var results *clap.Results
	if results, err = clap.Parse([]string{"-i", "photo.png"}, cfg); err != nil {
		if !errors.Is(err, clap.ErrMandatoryArgument) {
			t.Errorf("parsing error: %s", err)
			return
		}
		t.Logf("t: %v\n", err)
	}
	t.Logf("t: %v\n", results)
}

func TestInvalidTag(t *testing.T) {
	type config struct {
		Test string `clap:"--test,-wrongtag"`
	}
	cfg := &config{}
	var err error
	var results *clap.Results
	if results, err = clap.Parse([]string{"--test", "list"}, cfg); err != nil {
		if !errors.Is(err, clap.ErrInvalidTag) {
			t.Errorf("parsing error: %s", err)
			return
		}
		t.Logf("t: %v\n", err)
	}
	t.Logf("t: %v\n", results)
}

func TestTooManyOptions(t *testing.T) {
	type config struct {
		Test string `clap:"--test,-t,another"`
	}
	cfg := &config{}
	var err error
	var results *clap.Results
	if results, err = clap.Parse([]string{"--test", "list"}, cfg); err != nil {
		if !errors.Is(err, clap.ErrInvalidTag) {
			t.Errorf("parsing error: %s", err)
			return
		}
		t.Logf("parsing error: %s", err)
	}
	t.Logf("t: %v\n", results)
}

func TestShortName(t *testing.T) {
	type config struct {
		Short    string   `clap:",-S,mandatory"`
		Integer  int      `clap:",-I,mandatory"`
		Optional bool     `clap:",-O"`
		Slice    []string `clap:",-C,mandatory"`
	}
	cfg := &config{}
	var err error
	var results *clap.Results
	if results, err = clap.Parse([]string{
		"-S", "shortparam", "-O", "-I", "42", "-C", "a",
		"b", "c",
	}, cfg); err != nil {
		t.Errorf("parsing error: %s", err)
	}
	t.Logf("t: %v\n", results)
	wanted := &config{Short: "shortparam", Integer: 42, Optional: true, Slice: []string{"a", "b", "c"}}
	if !reflect.DeepEqual(cfg, wanted) {
		t.Errorf("wanted: '%v', got '%v'", wanted, cfg)
	}
}

func TestLongName(t *testing.T) {
	type config struct {
		Test  string `clap:"--test,mandatory"`
		Slice []int  `clap:"--slice"`
	}
	cfg := &config{}
	var err error
	var results *clap.Results
	if results, err = clap.Parse([]string{
		"image", "--test", "value", "-O", "-I", "42", "--slice", "22", "32",
		"42",
	}, cfg); err != nil {
		t.Errorf("parsing error: %s", err)
	}
	t.Logf("t: %v\n", results)
	wanted := &config{Test: "value", Slice: []int{22, 32, 42}}
	if !reflect.DeepEqual(cfg, wanted) {
		t.Errorf("wanted: '%v', got '%v'", wanted, cfg)
	}
	if !results.HasWarnings() || len(results.Ignored) != 4 {
		t.Errorf("wrong error number / type")
	}
}

func TestShortAndLong(t *testing.T) {
	type config struct {
		Slice []int `clap:"--slice,-S"`
	}
	cfg := &config{}
	var err error
	var results *clap.Results
	if results, err = clap.Parse([]string{
		"image", "-x", "shortparam", "-O", "-I", "42", "-S", "22", "32",
		"42",
	}, cfg); err != nil {
		t.Errorf("parsing error: %s", err)
	}
	t.Logf("t: %v\n", results)
	wanted := &config{Slice: []int{22, 32, 42}}
	if !reflect.DeepEqual(cfg, wanted) {
		t.Errorf("wanted: '%v', got '%v'", wanted, cfg)
	}
	if !results.HasWarnings() || len(results.Ignored) != 6 {
		t.Errorf("wrong error number / type")
	}
}

func TestComplete(t *testing.T) {
	type config struct {
		Extensions  []string `clap:"--extensions,-e,mandatory"`
		Recursive   bool     `clap:"--recusrive,-r"`
		Verbose     bool     `clap:"--verbose,-v"`
		Directories []string `clap:"trailing"`
	}
	cfg := &config{}
	var err error
	var results *clap.Results
	if results, err = clap.Parse([]string{
		"--extensions", "jpg", "png", "bmp", "-v", "$home/temp", "$home/tmp", "/tmp",
	}, cfg); err != nil {
		t.Errorf("parsing error: %s", err)
	}
	t.Logf("t: %v\n", results)
	wanted := &config{
		Extensions: []string{"jpg", "png", "bmp"}, Verbose: true,
		Directories: []string{"$home/temp", "$home/tmp", "/tmp"},
	}
	if !reflect.DeepEqual(cfg, wanted) {
		t.Errorf("wanted: '%v', got '%v'", wanted, cfg)
	}
}

func TestTypes(t *testing.T) {
	type config struct {
		String      string    `clap:"--string"`
		Int         int       `clap:"--int"`
		Float       float64   `clap:"--float"`
		Bool        bool      `clap:"--bool"`
		StringSlice []string  `clap:"--string-slice"`
		IntSlice    []int     `clap:"--int-slice"`
		StringArray [2]string `clap:"--string-array"`
		IntArray    [3]int    `clap:"--int-array"`
		Trailing    []string  `clap:"trailing"`
	}
	cfg := &config{}
	var err error
	var results *clap.Results
	if results, err = clap.Parse([]string{
		"--string", "str", "--int", "10", "--float", "12.3", "--bool", "--string-slice", "a", "b", "c",
		"--int-slice", "10", "11", "12", "--string-array", "a", "b", "--int-array", "10", "11", "12",
		"w", "x", "y", "z",
	}, cfg); err != nil {
		t.Errorf("parsing error: %s", err)
	}
	t.Logf("t: %v\n", results)
	wanted := &config{
		String:      "str",
		Int:         10,
		Float:       12.3,
		Bool:        true,
		StringSlice: []string{"a", "b", "c"},
		IntSlice:    []int{10, 11, 12},
		StringArray: [2]string{"a", "b"},
		IntArray:    [3]int{10, 11, 12},
		Trailing:    []string{"w", "x", "y", "z"},
	}
	if !reflect.DeepEqual(cfg, wanted) {
		t.Errorf("wanted: '%v', got '%v'", wanted, cfg)
	}
}

func TestReadme(t *testing.T) {
	type config struct {
		Cookie      string    `clap:"--cookie"`
		HTTPOnly    bool      `clap:"--httpOnly"`
		Secure      bool      `clap:"--secure"`
		Origins     [4]string `clap:"--origins,-O,mandatory"`
		Port        int       `clap:",-P,mandatory"`
		ConfigFiles []string  `clap:"trailing"`
	}
	cfg := &config{Secure: true}
	var err error
	var results *clap.Results
	if results, err = clap.Parse([]string{
		"-P", "8080", "--cookie", "clapcookie", "--httpOnly", "--origins", "http://localhost:5137",
		"https://localhost:5173", "http://localhost:3000", "https://localhost:3000",
		"config-db.json", "config-log.json",
	}, cfg); err != nil {
		t.Errorf("parsing error: %s", err)
	}
	t.Logf("t: %v\n", results)
	wanted := &config{
		Cookie:   "clapcookie",
		HTTPOnly: true,
		Secure:   true,
		Origins: [4]string{
			"http://localhost:5137", "https://localhost:5173",
			"http://localhost:3000", "https://localhost:3000",
		},
		Port:        8080,
		ConfigFiles: []string{"config-db.json", "config-log.json"},
	}
	if !reflect.DeepEqual(cfg, wanted) {
		t.Errorf("wanted: '%v', got '%v'", wanted, cfg)
	}
}
