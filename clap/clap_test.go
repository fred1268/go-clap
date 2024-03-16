package clap_test

import (
	"errors"
	"reflect"
	"slices"
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

func TestMultipleMandatoryNotFound(t *testing.T) {
	type config struct {
		Image int `clap:"--image,mandatory"`
		Num   int `clap:"--number,mandatory"`
	}
	wanted := []string{"image", "number"}

	cfg := &config{}
	var err error
	var results *clap.Results

	if results, err = clap.Parse([]string{"-i", "photo.png"}, cfg); err != nil {
		if !errors.Is(err, clap.ErrMandatoryArgument) {
			t.Errorf("parsing error: %s", err)
			return
		}
		got := results.Mandatory
		if len(got) != 2 || !slices.Contains(got, wanted[0]) || !slices.Contains(got, wanted[1]) {
			t.Errorf("wanted: '%v', got '%v'", wanted, got)
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
		Recursive   bool     `clap:"--recursive,-r"`
		Verbose     bool     `clap:"--verbose,-v"`
		Size        int      `clap:"--size,-s"`
		Directories []string `clap:"trailing"`
	}
	cfg := &config{}
	var err error
	var results *clap.Results
	if results, err = clap.Parse([]string{
		"--extensions", "jpg", "png", "bmp", "-v", "-s", "10", "$home/temp", "$home/tmp", "/tmp",
	}, cfg); err != nil {
		t.Errorf("parsing error: %s", err)
	}
	t.Logf("t: %v\n", results)
	wanted := &config{
		Extensions: []string{"jpg", "png", "bmp"}, Verbose: true, Size: 10,
		Directories: []string{"$home/temp", "$home/tmp", "/tmp"},
	}
	if !reflect.DeepEqual(cfg, wanted) {
		t.Errorf("wanted: '%v', got '%v'", wanted, cfg)
	}
}

func TestBooleans(t *testing.T) {
	type config struct {
		Recursive bool `clap:"--recursive,-R"`
	}
	cfg := &config{
		Recursive: false,
	}
	var err error
	var results *clap.Results
	if results, err = clap.Parse([]string{"--recursive"}, cfg); err != nil {
		t.Errorf("parsing error: %s", err)
	}
	t.Logf("t: %v\n", results)
	wanted := &config{Recursive: true}
	if !reflect.DeepEqual(cfg, wanted) {
		t.Errorf("wanted: '%v', got '%v'", wanted, cfg)
	}
	// with --no-recursive
	cfg = &config{
		Recursive: true,
	}
	if results, err = clap.Parse([]string{"--no-recursive"}, cfg); err != nil {
		t.Errorf("parsing error: %s", err)
	}
	t.Logf("t: %v\n", results)
	wanted = &config{Recursive: false}
	if !reflect.DeepEqual(cfg, wanted) {
		t.Errorf("wanted: '%v', got '%v'", wanted, cfg)
	}
}

func TestTypes(t *testing.T) {
	type config struct {
		String      string    `clap:"--string"`
		Int         int       `clap:"--int"`
		Int8        int8      `clap:"--int8"`
		Int16       int16     `clap:"--int16"`
		Int32       int32     `clap:"--int32"`
		Int64       int64     `clap:"--int64"`
		UInt        uint      `clap:"--uint"`
		UInt8       uint8     `clap:"--uint8"`
		UInt16      uint16    `clap:"--uint16"`
		UInt32      uint32    `clap:"--uint32"`
		UInt64      uint64    `clap:"--uint64"`
		Float32     float32   `clap:"--float32"`
		Float64     float64   `clap:"--float64"`
		Bool        bool      `clap:"--bool"`
		DefaultTrue bool      `clap:"--defaulttrue"`
		StringSlice []string  `clap:"--string-slice"`
		IntSlice    []int     `clap:"--int-slice"`
		StringArray [2]string `clap:"--string-array"`
		IntArray    [3]int    `clap:"--int-array"`
		Trailing    []string  `clap:"trailing"`
	}
	cfg := &config{
		DefaultTrue: true,
	}
	var err error
	var results *clap.Results
	if results, err = clap.Parse([]string{
		"--string", "str", "--int", "10", "--int8", "8", "--int16", "16", "--int32", "32", "--int64", "64",
		"--uint", "12", "--uint8", "65535", "--uint16", "65535", "--uint32", "65535", "--uint64", "65535",
		"--float32", "12.32", "--float64", "12.64", "--bool", "--no-defaulttrue", "--string-slice", "a", "b", "c",
		"--int-slice", "10", "11", "12", "--string-array", "a", "b", "--int-array", "10", "11", "12",
		"w", "x", "y", "z",
	}, cfg); err != nil {
		t.Errorf("parsing error: %s", err)
	}
	t.Logf("t: %v\n", results)
	wanted := &config{
		String:      "str",
		Int:         10,
		Int8:        8,
		Int16:       16,
		Int32:       32,
		Int64:       64,
		UInt:        12,
		UInt8:       255,
		UInt16:      65535,
		UInt32:      65535,
		UInt64:      65535,
		Float32:     12.32,
		Float64:     12.64,
		Bool:        true,
		DefaultTrue: false,
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
