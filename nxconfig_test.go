package nxconfig

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadFromEnv(t *testing.T) {
	type cfg struct {
		Str       string
		Int       int
		UInt      uint
		NotMapped float64
		Nested    struct {
			Float float32
			Str2  string
		}
	}

	expected := cfg{
		Str:       "foo",
		Int:       -910,
		UInt:      10,
		NotMapped: 0.0,
		Nested: struct {
			Float float32
			Str2  string
		}{3.49, "bar"},
	}

	var input cfg

	err := Load(&input, WithEnv([]string{
		"STR=foo",
		"INT=-910",
		"U_INT=10",
		"NESTED_FLOAT=3.49",
		"NESTED_STR2=bar",
	}))
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, expected, input)
}

func TestLoadFromEnvWithPrefix(t *testing.T) {
	type cfg struct {
		Str       string
		Int       int
		UInt      uint
		NotMapped float64
		Nested    struct {
			Float float32
			Str2  string
		}
	}

	expected := cfg{
		Str:       "foo",
		Int:       -910,
		UInt:      10,
		NotMapped: 0.0,
		Nested: struct {
			Float float32
			Str2  string
		}{3.49, "bar"},
	}

	var input cfg

	err := Load(&input, WithEnv([]string{
		"PREFIX_STR=foo",
		"PREFIX_INT=-910",
		"PREFIX_U_INT=10",
		"PREFIX_NESTED_FLOAT=3.49",
		"PREFIX_NESTED_STR2=bar",
	}), WithPrefix("PREFIX"))
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, expected, input)
}

func TestLoadFromArgs(t *testing.T) {
	type cfg struct {
		Str       string
		Int       int
		UInt      uint
		NotMapped float64
		Nested    struct {
			Float float32
			Str2  string
		}
	}

	expected := cfg{
		Str:       "foo",
		Int:       -910,
		UInt:      10,
		NotMapped: 0.0,
		Nested: struct {
			Float float32
			Str2  string
		}{3.49, "bar"},
	}

	var input cfg

	err := Load(&input, WithArgs([]string{
		"--str", "foo",
		"--int=-910",
		"--u-int", "10",
		"--nested-float", "3.49",
		"--nested-str2", "bar",
	}))
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, expected, input)
}

func TestLoadFromArgsAndEnv(t *testing.T) {
	type cfg struct {
		Str       string
		Override  string
		Int       int
		UInt      uint
		NotMapped float64
		Nested    struct {
			Float float32
			Str2  string
		}
	}

	expected := cfg{
		Str:       "foo",
		Override:  "baz",
		Int:       -910,
		UInt:      10,
		NotMapped: 0.0,
		Nested: struct {
			Float float32
			Str2  string
		}{3.49, "bar"},
	}

	var input cfg

	err := Load(&input, WithArgs([]string{
		"--int", "-910",
		"--str", "foo",
		"--override", "baz",
		"--nested-str2", "bar",
	}), WithEnv([]string{
		"OVERRIDE=wrong",
		"U_INT=10",
		"NESTED_FLOAT=3.49",
	}))
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, expected, input)
}

func TestLoadUnexpoprted(t *testing.T) {
	type cfg struct {
		Str        string
		unexported string
		Nested     struct {
			Float float32
			ui    uint32
		}
		unexportedNested struct {
			Str2 string
		}
	}

	expected := cfg{
		Str:        "foo",
		unexported: "",
		Nested: struct {
			Float float32
			ui    uint32
		}{45.6, 0},
		unexportedNested: struct {
			Str2 string
		}{""},
	}

	var input cfg

	err := Load(&input, WithArgs([]string{
		"--str", "foo",
		"--unexported", "baz",
		"--nested-float", "45.6",
		"--unexported-nested-str2", "wrong",
	}))
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, expected, input)
}

func TestLoadStructTags(t *testing.T) {
	type cfg struct {
		Str       string
		Int       int `nxconfig:"Different"`
		UInt      uint
		NotMapped float64
		Nested    struct {
			Float float32
			Str2  string `nxconfig:"NoLongerNested"`
		}
	}

	expected := cfg{
		Str:       "foo",
		Int:       -910,
		UInt:      10,
		NotMapped: 0.0,
		Nested: struct {
			Float float32
			Str2  string `nxconfig:"NoLongerNested"`
		}{3.49, "bar"},
	}

	var input cfg

	err := Load(&input, WithArgs([]string{
		"--different", "-910",
		"--str", "foo",
		"--override", "baz",
		"--no-longer-nested", "bar",
	}), WithEnv([]string{
		"U_INT=10",
		"NESTED_FLOAT=3.49",
	}))
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, expected, input)
}

func TestToKebabCase(t *testing.T) {
	tt := []struct {
		input  string
		expect string
	}{
		{"", ""},
		{"already-kebab", "already-kebab"},
		{"A", "a"},
		{"AA", "aa"},
		{"AaAa", "aa-aa"},
		{"HTTPRequest", "http-request"},
		{"BatteryLifeValue", "battery-life-value"},
		{"Id0Value", "id0-value"},
		{"ID0Value", "id0-value"},
	}

	for _, ts := range tt {
		require.Equal(t, ts.expect, toKebabCase(ts.input))
	}
}

func TestToUpperSnakeCase(t *testing.T) {
	tt := []struct {
		input  string
		expect string
	}{
		{"", ""},
		{"UPPER_SNAKE", "UPPER_SNAKE"},
		{"A", "A"},
		{"AA", "AA"},
		{"AaAa", "AA_AA"},
		{"HTTPRequest", "HTTP_REQUEST"},
		{"BatteryLifeValue", "BATTERY_LIFE_VALUE"},
		{"Id0Value", "ID0_VALUE"},
		{"ID0Value", "ID0_VALUE"},
	}

	for _, ts := range tt {
		require.Equal(t, ts.expect, toUpperSnakeCase(ts.input))
	}
}
