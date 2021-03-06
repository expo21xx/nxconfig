package nxconfig

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadFromEnv(t *testing.T) {
	type cfg struct {
		Str       string
		Int       int
		UInt      uint
		NotMapped float64
		Duration  time.Duration
		StrSlice  []string
		Nested    struct {
			Float float32
			Str2  string
		}
		NestedPtr *struct {
			Int int
		}
	}

	expected := cfg{
		Str:       "foo",
		Int:       -910,
		UInt:      10,
		NotMapped: 0.0,
		Duration:  time.Hour * 2,
		StrSlice:  []string{"a", "b"},
		Nested: struct {
			Float float32
			Str2  string
		}{3.49, "bar"},
		NestedPtr: &struct {
			Int int
		}{89},
	}

	var input cfg

	err := Load(&input, WithEnv([]string{
		"STR=foo",
		"INT=-910",
		"U_INT=10",
		"DURATION=2h",
		"NESTED_FLOAT=3.49",
		"NESTED_STR2=bar",
		"STR_SLICE=a,b",
		"NESTED_PTR_INT=89",
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
		NestedPtr *struct {
			Int int
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
		NestedPtr: &struct {
			Int int
		}{89},
	}

	var input cfg

	err := Load(&input, WithEnv([]string{
		"PREFIX_STR=foo",
		"PREFIX_INT=-910",
		"PREFIX_U_INT=10",
		"PREFIX_NESTED_FLOAT=3.49",
		"PREFIX_NESTED_STR2=bar",
		"PREFIX_NESTED_PTR_INT=89",
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
		Duration  time.Duration
		StrSlice  []string
		Nested    struct {
			Float float32
			Str2  string
		}
		NestedPtr *struct {
			Int int
		}
	}

	expected := cfg{
		Str:       "foo",
		Int:       -910,
		UInt:      10,
		NotMapped: 0.0,
		Duration:  time.Minute * 5,
		StrSlice:  []string{"a", "b"},
		Nested: struct {
			Float float32
			Str2  string
		}{3.49, "bar"},
		NestedPtr: &struct {
			Int int
		}{89},
	}

	var input cfg

	err := Load(&input, WithArgs([]string{
		"--str", "foo",
		"--int=-910",
		"--u-int", "10",
		"--duration", "5m",
		"--nested-float", "3.49",
		"--nested-str2", "bar",
		"--str-slice=a,b",
		"--nested-ptr-int", "89",
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
		StrSlice  []string
		Nested    struct {
			Float float32
			Str2  string
		}
		NestedPtr *struct {
			Int int
		}
	}

	expected := cfg{
		Str:       "foo",
		Override:  "baz",
		Int:       -910,
		UInt:      10,
		NotMapped: 0.0,
		StrSlice:  []string{"b", "c"},
		Nested: struct {
			Float float32
			Str2  string
		}{3.49, "bar"},
		NestedPtr: &struct {
			Int int
		}{89},
	}

	var input cfg

	err := Load(&input, WithArgs([]string{
		"--int", "-910",
		"--str", "foo",
		"--override", "baz",
		"--nested-str2", "bar",
		"--str-slice=b",
		"--str-slice=c",
		"--nested-ptr-int=89",
	}), WithEnv([]string{
		"OVERRIDE=wrong",
		"U_INT=10",
		"NESTED_FLOAT=3.49",
		"STR_SLICE=a,b",
		"NESTED_PTR_INT=87",
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
		Int       int `name:"Different"`
		UInt      uint
		NoValue   float64 `default:"89.7"`
		NotMapped float64 `name:""`
		Nested    struct {
			Float float32 `default:"789.9"`
			Str2  string  `name:"NoLongerNested"`
		}
		NestedPtr *struct {
			Int int `name:"ForPtrNested"`
		}
	}

	expected := cfg{
		Str:       "foo",
		Int:       -910,
		UInt:      10,
		NoValue:   89.7,
		NotMapped: 0,
		Nested: struct {
			Float float32 `default:"789.9"`
			Str2  string  `name:"NoLongerNested"`
		}{3.49, "bar"},
		NestedPtr: &struct {
			Int int `name:"ForPtrNested"`
		}{89},
	}

	var input cfg

	err := Load(&input, WithArgs([]string{
		"--different", "-910",
		"--str", "foo",
		"--override", "baz",
		"--not-mapped", "59.2",
		"--no-longer-nested", "bar",
	}), WithEnv([]string{
		"U_INT=10",
		"NESTED_FLOAT=3.49",
		"FOR_PTR_NESTED=89",
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
