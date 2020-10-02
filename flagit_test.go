package flagit

import (
	"errors"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type flags struct {
	unexported    string
	String        string          `flag:"string"`
	Bool          bool            `flag:"bool"`
	Float32       float32         `flag:"float32"`
	Float64       float64         `flag:"float64"`
	Int           int             `flag:"int"`
	Int8          int8            `flag:"int8"`
	Int16         int16           `flag:"int16"`
	Int32         int32           `flag:"int32"`
	Int64         int64           `flag:"int64"`
	Uint          uint            `flag:"uint"`
	Uint8         uint8           `flag:"uint8"`
	Uint16        uint16          `flag:"uint16"`
	Uint32        uint32          `flag:"uint32"`
	Uint64        uint64          `flag:"uint64"`
	Duration      time.Duration   `flag:"duration"`
	URL           url.URL         `flag:"url"`
	StringSlice   []string        `flag:"string-slice"`
	BoolSlice     []bool          `flag:"bool-slice"`
	Float32Slice  []float32       `flag:"float32-slice"`
	Float64Slice  []float64       `flag:"float64-slice"`
	IntSlice      []int           `flag:"int-slice"`
	Int8Slice     []int8          `flag:"int8-slice"`
	Int16Slice    []int16         `flag:"int16-slice"`
	Int32Slice    []int32         `flag:"int32-slice"`
	Int64Slice    []int64         `flag:"int64-slice"`
	UintSlice     []uint          `flag:"uint-slice"`
	Uint8Slice    []uint8         `flag:"uint8-slice"`
	Uint16Slice   []uint16        `flag:"uint16-slice"`
	Uint32Slice   []uint32        `flag:"uint32-slice"`
	Uint64Slice   []uint64        `flag:"uint64-slice"`
	DurationSlice []time.Duration `flag:"duration-slice"`
	URLSlice      []url.URL       `flag:"url-slice"`
}

func TestValidateStruct(t *testing.T) {
	tests := []struct {
		name          string
		s             interface{}
		expectedError error
	}{
		{
			"NonStruct",
			new(string),
			errors.New("non-struct type: you should pass a pointer to a struct type"),
		},
		{
			"NonPointer",
			flags{},
			errors.New("non-pointer type: you should pass a pointer to a struct type"),
		},
		{
			"OK",
			&flags{},
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v, err := validateStruct(tc.s)

			if tc.expectedError == nil {
				assert.NotNil(t, v)
				assert.NoError(t, err)
			} else {
				assert.Empty(t, v)
				assert.Equal(t, tc.expectedError, err)
			}
		})
	}
}

func TestIsTypeSupported(t *testing.T) {
	service1URL, _ := url.Parse("service-1:8080")
	service2URL, _ := url.Parse("service-2:8080")

	tests := []struct {
		name     string
		field    interface{}
		expected bool
	}{
		{"String", "dummy", true},
		{"Bool", true, true},
		{"Float32", float32(3.1415), true},
		{"Float64", float64(3.14159265359), true},
		{"Int", int(-2147483648), true},
		{"Int8", int8(-128), true},
		{"Int16", int16(-32768), true},
		{"Int32", int32(-2147483648), true},
		{"Int64", int64(-9223372036854775808), true},
		{"Duration", time.Hour, true},
		{"Uint", uint(4294967295), true},
		{"Uint8", uint8(255), true},
		{"Uint16", uint16(65535), true},
		{"Uint32", uint32(4294967295), true},
		{"Uint64", uint64(18446744073709551615), true},
		{"URL", *service1URL, true},
		{"StringSlice", []string{"foo", "bar"}, true},
		{"BoolSlice", []bool{true, false}, true},
		{"Float32Slice", []float32{3.1415, 2.7182}, true},
		{"Float64Slice", []float64{3.14159265359, 2.71828182845}, true},
		{"IntSlice", []int{}, true},
		{"Int8Slice", []int8{}, true},
		{"Int16Slice", []int16{}, true},
		{"Int32Slice", []int32{}, true},
		{"Int64Slice", []int64{}, true},
		{"DurationSlice", []time.Duration{}, true},
		{"UintSlice", []uint{}, true},
		{"Uint8Slice", []uint8{}, true},
		{"Uint16Slice", []uint16{}, true},
		{"Uint32Slice", []uint32{}, true},
		{"Uint64Slice", []uint64{}, true},
		{"URLSlice", []url.URL{*service1URL, *service2URL}, true},
		{"Unsupported", time.Now(), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tp := reflect.TypeOf(tc.field)

			assert.Equal(t, tc.expected, isTypeSupported(tp))
		})
	}
}

func TestGetFlagValue(t *testing.T) {
	tests := []struct {
		args              []string
		flagName          string
		expectedFlagValue string
	}{
		{[]string{"app", "invalid"}, "invalid", ""},

		{[]string{"app", "-enabled"}, "enabled", "true"},
		{[]string{"app", "--enabled"}, "enabled", "true"},
		{[]string{"app", "-enabled=false"}, "enabled", "false"},
		{[]string{"app", "--enabled=false"}, "enabled", "false"},
		{[]string{"app", "-enabled", "false"}, "enabled", "false"},
		{[]string{"app", "--enabled", "false"}, "enabled", "false"},

		{[]string{"app", "-number=-10"}, "number", "-10"},
		{[]string{"app", "--number=-10"}, "number", "-10"},
		{[]string{"app", "-number", "-10"}, "number", "-10"},
		{[]string{"app", "--number", "-10"}, "number", "-10"},

		{[]string{"app", "-text=content"}, "text", "content"},
		{[]string{"app", "--text=content"}, "text", "content"},
		{[]string{"app", "-text", "content"}, "text", "content"},
		{[]string{"app", "--text", "content"}, "text", "content"},

		{[]string{"app", "-enabled", "-text=content"}, "enabled", "true"},
		{[]string{"app", "--enabled", "--text=content"}, "enabled", "true"},
		{[]string{"app", "-enabled", "-text", "content"}, "enabled", "true"},
		{[]string{"app", "--enabled", "--text", "content"}, "enabled", "true"},

		{[]string{"app", "-name-list=alice,bob"}, "name-list", "alice,bob"},
		{[]string{"app", "--name-list=alice,bob"}, "name-list", "alice,bob"},
		{[]string{"app", "-name-list", "alice,bob"}, "name-list", "alice,bob"},
		{[]string{"app", "--name-list", "alice,bob"}, "name-list", "alice,bob"},
	}

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	for _, tc := range tests {
		os.Args = tc.args
		flagValue := getFlagValue(tc.flagName)

		assert.Equal(t, tc.expectedFlagValue, flagValue)
	}
}

func TestIterateOnFields(t *testing.T) {
	tests := []struct {
		name               string
		s                  interface{}
		expectedFieldNames []string
		expectedFlagNames  []string
		expectedListSeps   []string
		expectedError      error
	}{
		{
			name: "OK",
			s:    &flags{},
			expectedFieldNames: []string{
				"String",
				"Bool",
				"Float32", "Float64",
				"Int", "Int8", "Int16", "Int32", "Int64",
				"Uint", "Uint8", "Uint16", "Uint32", "Uint64",
				"Duration", "URL",
				"StringSlice",
				"BoolSlice",
				"Float32Slice", "Float64Slice",
				"IntSlice", "Int8Slice", "Int16Slice", "Int32Slice", "Int64Slice",
				"UintSlice", "Uint8Slice", "Uint16Slice", "Uint32Slice", "Uint64Slice",
				"DurationSlice", "URLSlice",
			},
			expectedFlagNames: []string{
				"string",
				"bool",
				"float32", "float64",
				"int", "int8", "int16", "int32", "int64",
				"uint", "uint8", "uint16", "uint32", "uint64",
				"duration", "url",
				"string-slice",
				"bool-slice",
				"float32-slice", "float64-slice",
				"int-slice", "int8-slice", "int16-slice", "int32-slice", "int64-slice",
				"uint-slice", "uint8-slice", "uint16-slice", "uint32-slice", "uint64-slice",
				"duration-slice", "url-slice",
			},
			expectedListSeps: []string{
				",",
				",",
				",", ",",
				",", ",", ",", ",", ",",
				",", ",", ",", ",", ",",
				",", ",",
				",",
				",",
				",", ",",
				",", ",", ",", ",", ",",
				",", ",", ",", ",", ",",
				",", ",",
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fieldNames := []string{}
			flagNames := []string{}
			listSeps := []string{}

			vStruct, err := validateStruct(tc.s)
			assert.NoError(t, err)

			iterateOnFields(vStruct, func(v reflect.Value, fieldName, flagName, listSep string) {
				// values = append(values, v)
				fieldNames = append(fieldNames, fieldName)
				flagNames = append(flagNames, flagName)
				listSeps = append(listSeps, listSep)
			})

			assert.Equal(t, tc.expectedFieldNames, fieldNames)
			assert.Equal(t, tc.expectedFlagNames, flagNames)
			assert.Equal(t, tc.expectedListSeps, listSeps)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}