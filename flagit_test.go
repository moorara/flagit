package flagit

import (
	"errors"
	"flag"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type flags struct {
	unexported  string
	WithoutFlag string

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
	url1, _ := url.Parse("service-1:8080")
	url2, _ := url.Parse("service-2:8080")

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
		{"URL", *url1, true},
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
		{"URLSlice", []url.URL{*url1, *url2}, true},
		{"Unsupported", time.Now(), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			typ := reflect.TypeOf(tc.field)

			assert.Equal(t, tc.expected, isTypeSupported(typ))
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

func TestSetFieldValue(t *testing.T) {
	d90m := 90 * time.Minute
	d120m := 120 * time.Minute
	d4h := 4 * time.Hour
	d8h := 8 * time.Hour

	url1, _ := url.Parse("service-1:8080")
	url2, _ := url.Parse("service-2:8080")
	url3, _ := url.Parse("service-3:8080")
	url4, _ := url.Parse("service-4:8080")

	tests := []struct {
		name           string
		s              *flags
		values         map[string]string
		expectedResult bool
		expected       *flags
	}{
		{
			"NewValues",
			&flags{
				String:        "default",
				Bool:          false,
				Float32:       3.1415,
				Float64:       3.14159265359,
				Int:           -2147483648,
				Int8:          -128,
				Int16:         -32768,
				Int32:         -2147483648,
				Int64:         -9223372036854775808,
				Uint:          4294967295,
				Uint8:         255,
				Uint16:        65535,
				Uint32:        4294967295,
				Uint64:        18446744073709551615,
				Duration:      d90m,
				URL:           *url1,
				StringSlice:   []string{"milad", "mona"},
				BoolSlice:     []bool{false, true},
				Float32Slice:  []float32{3.1415, 2.7182},
				Float64Slice:  []float64{3.14159265359, 2.71828182845},
				IntSlice:      []int{-2147483648, 2147483647},
				Int8Slice:     []int8{-128, 127},
				Int16Slice:    []int16{-32768, 32767},
				Int32Slice:    []int32{-2147483648, 2147483647},
				Int64Slice:    []int64{-9223372036854775808, 9223372036854775807},
				UintSlice:     []uint{0, 4294967295},
				Uint8Slice:    []uint8{0, 255},
				Uint16Slice:   []uint16{0, 65535},
				Uint32Slice:   []uint32{0, 4294967295},
				Uint64Slice:   []uint64{0, 18446744073709551615},
				DurationSlice: []time.Duration{d90m, d120m},
				URLSlice:      []url.URL{*url1, *url2},
			},
			map[string]string{
				"String":        "content",
				"Bool":          "true",
				"Float32":       "2.7182",
				"Float64":       "2.7182818284",
				"Int":           "2147483647",
				"Int8":          "127",
				"Int16":         "32767",
				"Int32":         "2147483647",
				"Int64":         "9223372036854775807",
				"Uint":          "2147483648",
				"Uint8":         "128",
				"Uint16":        "32768",
				"Uint32":        "2147483648",
				"Uint64":        "9223372036854775808",
				"Duration":      "4h",
				"URL":           "service-3:8080",
				"StringSlice":   "mona,milad",
				"BoolSlice":     "true,false",
				"Float32Slice":  "2.7182,3.1415",
				"Float64Slice":  "2.71828182845,3.14159265359",
				"IntSlice":      "2147483647,-2147483648",
				"Int8Slice":     "127,-128",
				"Int16Slice":    "32767,-32768",
				"Int32Slice":    "2147483647,-2147483648",
				"Int64Slice":    "9223372036854775807,-9223372036854775808",
				"UintSlice":     "4294967295,0",
				"Uint8Slice":    "255,0",
				"Uint16Slice":   "65535,0",
				"Uint32Slice":   "4294967295,0",
				"Uint64Slice":   "18446744073709551615,0",
				"DurationSlice": "4h,8h",
				"URLSlice":      "service-3:8080,service-4:8080",
			},
			true,
			&flags{
				String:        "content",
				Bool:          true,
				Float32:       2.7182,
				Float64:       2.7182818284,
				Int:           2147483647,
				Int8:          127,
				Int16:         32767,
				Int32:         2147483647,
				Int64:         9223372036854775807,
				Uint:          2147483648,
				Uint8:         128,
				Uint16:        32768,
				Uint32:        2147483648,
				Uint64:        9223372036854775808,
				Duration:      d4h,
				URL:           *url3,
				StringSlice:   []string{"mona", "milad"},
				BoolSlice:     []bool{true, false},
				Float32Slice:  []float32{2.7182, 3.1415},
				Float64Slice:  []float64{2.71828182845, 3.14159265359},
				IntSlice:      []int{2147483647, -2147483648},
				Int8Slice:     []int8{127, -128},
				Int16Slice:    []int16{32767, -32768},
				Int32Slice:    []int32{2147483647, -2147483648},
				Int64Slice:    []int64{9223372036854775807, -9223372036854775808},
				UintSlice:     []uint{4294967295, 0},
				Uint8Slice:    []uint8{255, 0},
				Uint16Slice:   []uint16{65535, 0},
				Uint32Slice:   []uint32{4294967295, 0},
				Uint64Slice:   []uint64{18446744073709551615, 0},
				DurationSlice: []time.Duration{d4h, d8h},
				URLSlice:      []url.URL{*url3, *url4},
			},
		},
		{
			"NoNewValues",
			&flags{
				String:        "content",
				Bool:          true,
				Float32:       2.7182,
				Float64:       2.7182818284,
				Int:           2147483647,
				Int8:          127,
				Int16:         32767,
				Int32:         2147483647,
				Int64:         9223372036854775807,
				Uint:          2147483648,
				Uint8:         128,
				Uint16:        32768,
				Uint32:        2147483648,
				Uint64:        9223372036854775808,
				Duration:      d4h,
				URL:           *url3,
				StringSlice:   []string{"mona", "milad"},
				BoolSlice:     []bool{true, false},
				Float32Slice:  []float32{2.7182, 3.1415},
				Float64Slice:  []float64{2.71828182845, 3.14159265359},
				IntSlice:      []int{2147483647, -2147483648},
				Int8Slice:     []int8{127, -128},
				Int16Slice:    []int16{32767, -32768},
				Int32Slice:    []int32{2147483647, -2147483648},
				Int64Slice:    []int64{9223372036854775807, -9223372036854775808},
				UintSlice:     []uint{4294967295, 0},
				Uint8Slice:    []uint8{255, 0},
				Uint16Slice:   []uint16{65535, 0},
				Uint32Slice:   []uint32{4294967295, 0},
				Uint64Slice:   []uint64{18446744073709551615, 0},
				DurationSlice: []time.Duration{d4h, d8h},
				URLSlice:      []url.URL{*url3, *url4},
			},
			map[string]string{
				"String":        "content",
				"Bool":          "true",
				"Float32":       "2.7182",
				"Float64":       "2.7182818284",
				"Int":           "2147483647",
				"Int8":          "127",
				"Int16":         "32767",
				"Int32":         "2147483647",
				"Int64":         "9223372036854775807",
				"Uint":          "2147483648",
				"Uint8":         "128",
				"Uint16":        "32768",
				"Uint32":        "2147483648",
				"Uint64":        "9223372036854775808",
				"Duration":      "4h",
				"URL":           "service-3:8080",
				"StringSlice":   "mona,milad",
				"BoolSlice":     "true,false",
				"Float32Slice":  "2.7182,3.1415",
				"Float64Slice":  "2.71828182845,3.14159265359",
				"IntSlice":      "2147483647,-2147483648",
				"Int8Slice":     "127,-128",
				"Int16Slice":    "32767,-32768",
				"Int32Slice":    "2147483647,-2147483648",
				"Int64Slice":    "9223372036854775807,-9223372036854775808",
				"UintSlice":     "4294967295,0",
				"Uint8Slice":    "255,0",
				"Uint16Slice":   "65535,0",
				"Uint32Slice":   "4294967295,0",
				"Uint64Slice":   "18446744073709551615,0",
				"DurationSlice": "4h,8h",
				"URLSlice":      "service-3:8080,service-4:8080",
			},
			false,
			&flags{
				String:        "content",
				Bool:          true,
				Float32:       2.7182,
				Float64:       2.7182818284,
				Int:           2147483647,
				Int8:          127,
				Int16:         32767,
				Int32:         2147483647,
				Int64:         9223372036854775807,
				Uint:          2147483648,
				Uint8:         128,
				Uint16:        32768,
				Uint32:        2147483648,
				Uint64:        9223372036854775808,
				Duration:      d4h,
				URL:           *url3,
				StringSlice:   []string{"mona", "milad"},
				BoolSlice:     []bool{true, false},
				Float32Slice:  []float32{2.7182, 3.1415},
				Float64Slice:  []float64{2.71828182845, 3.14159265359},
				IntSlice:      []int{2147483647, -2147483648},
				Int8Slice:     []int8{127, -128},
				Int16Slice:    []int16{32767, -32768},
				Int32Slice:    []int32{2147483647, -2147483648},
				Int64Slice:    []int64{9223372036854775807, -9223372036854775808},
				UintSlice:     []uint{4294967295, 0},
				Uint8Slice:    []uint8{255, 0},
				Uint16Slice:   []uint16{65535, 0},
				Uint32Slice:   []uint32{4294967295, 0},
				Uint64Slice:   []uint64{18446744073709551615, 0},
				DurationSlice: []time.Duration{d4h, d8h},
				URLSlice:      []url.URL{*url3, *url4},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vStruct := reflect.ValueOf(tc.s).Elem()
			for i := 0; i < vStruct.NumField(); i++ {
				v := vStruct.Field(i)
				f := vStruct.Type().Field(i)

				// Only consider those fields that are exported, supported, and have flag tag
				if v.CanSet() && isTypeSupported(v.Type()) && f.Tag.Get(flagTag) != "" {
					f := fieldInfo{
						value:   v,
						name:    f.Name,
						listSep: ",",
					}

					res := setFieldValue(f, tc.values[f.name])
					assert.Equal(t, tc.expectedResult, res)
				}
			}

			assert.Equal(t, tc.expected, tc.s)
		})
	}
}

func TestIterateOnFields(t *testing.T) {
	tests := []struct {
		name               string
		s                  interface{}
		expectedError      error
		expectedFieldNames []string
		expectedFlagNames  []string
		expectedListSeps   []string
	}{
		{
			name:          "OK",
			s:             &flags{},
			expectedError: nil,
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
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fieldNames := []string{}
			flagNames := []string{}
			listSeps := []string{}

			vStruct, err := validateStruct(tc.s)
			assert.NoError(t, err)

			err = iterateOnFields(vStruct, func(f fieldInfo) error {
				fieldNames = append(fieldNames, f.name)
				flagNames = append(flagNames, f.flag)
				listSeps = append(listSeps, f.listSep)
				return nil
			})

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedFieldNames, fieldNames)
			assert.Equal(t, tc.expectedFlagNames, flagNames)
			assert.Equal(t, tc.expectedListSeps, listSeps)
		})
	}
}

func TestPopulate(t *testing.T) {
	d90m := 90 * time.Minute
	d120m := 120 * time.Minute

	url1, _ := url.Parse("service-1:8080")
	url2, _ := url.Parse("service-2:8080")

	tests := []struct {
		name          string
		args          []string
		s             interface{}
		expectedError error
		expected      *flags
	}{
		{
			"NonStruct",
			[]string{"app"},
			new(string),
			errors.New("non-struct type: you should pass a pointer to a struct type"),
			&flags{},
		},
		{
			"NonPointer",
			[]string{"app"},
			flags{},
			errors.New("non-pointer type: you should pass a pointer to a struct type"),
			&flags{},
		},
		{
			"FromDefaults",
			[]string{"app"},
			&flags{
				unexported:    "internal",
				String:        "default",
				Bool:          false,
				Float32:       3.1415,
				Float64:       3.14159265359,
				Int:           -2147483648,
				Int8:          -128,
				Int16:         -32768,
				Int32:         -2147483648,
				Int64:         -9223372036854775808,
				Uint:          4294967295,
				Uint8:         255,
				Uint16:        65535,
				Uint32:        4294967295,
				Uint64:        18446744073709551615,
				Duration:      d90m,
				URL:           *url1,
				StringSlice:   []string{"milad", "mona"},
				BoolSlice:     []bool{false, true},
				Float32Slice:  []float32{3.1415, 2.7182},
				Float64Slice:  []float64{3.14159265359, 2.71828182845},
				IntSlice:      []int{-2147483648, 2147483647},
				Int8Slice:     []int8{-128, 127},
				Int16Slice:    []int16{-32768, 32767},
				Int32Slice:    []int32{-2147483648, 2147483647},
				Int64Slice:    []int64{-9223372036854775808, 9223372036854775807},
				UintSlice:     []uint{0, 4294967295},
				Uint8Slice:    []uint8{0, 255},
				Uint16Slice:   []uint16{0, 65535},
				Uint32Slice:   []uint32{0, 4294967295},
				Uint64Slice:   []uint64{0, 18446744073709551615},
				DurationSlice: []time.Duration{d90m, d120m},
				URLSlice:      []url.URL{*url1, *url2},
			},
			nil,
			&flags{
				unexported:    "internal",
				String:        "default",
				Bool:          false,
				Float32:       3.1415,
				Float64:       3.14159265359,
				Int:           -2147483648,
				Int8:          -128,
				Int16:         -32768,
				Int32:         -2147483648,
				Int64:         -9223372036854775808,
				Uint:          4294967295,
				Uint8:         255,
				Uint16:        65535,
				Uint32:        4294967295,
				Uint64:        18446744073709551615,
				Duration:      d90m,
				URL:           *url1,
				StringSlice:   []string{"milad", "mona"},
				BoolSlice:     []bool{false, true},
				Float32Slice:  []float32{3.1415, 2.7182},
				Float64Slice:  []float64{3.14159265359, 2.71828182845},
				IntSlice:      []int{-2147483648, 2147483647},
				Int8Slice:     []int8{-128, 127},
				Int16Slice:    []int16{-32768, 32767},
				Int32Slice:    []int32{-2147483648, 2147483647},
				Int64Slice:    []int64{-9223372036854775808, 9223372036854775807},
				UintSlice:     []uint{0, 4294967295},
				Uint8Slice:    []uint8{0, 255},
				Uint16Slice:   []uint16{0, 65535},
				Uint32Slice:   []uint32{0, 4294967295},
				Uint64Slice:   []uint64{0, 18446744073709551615},
				DurationSlice: []time.Duration{d90m, d120m},
				URLSlice:      []url.URL{*url1, *url2},
			},
		},
		{
			"FromFlags#1",
			[]string{
				"app",
				"-string", "content",
				"-bool",
				"-float32", "3.1415",
				"-float64", "3.14159265359",
				"-int", "-2147483648",
				"-int8", "-128",
				"-int16", "-32768",
				"-int32", "-2147483648",
				"-int64", "-9223372036854775808",
				"-uint", "4294967295",
				"-uint8", "255",
				"-uint16", "65535",
				"-uint32", "4294967295",
				"-uint64", "18446744073709551615",
				"-duration", "90m",
				"-url", "service-1:8080",
				"-string-slice", "milad,mona",
				"-bool-slice", "false,true",
				"-float32-slice", "3.1415,2.7182",
				"-float64-slice", "3.14159265359,2.71828182845",
				"-int-slice", "-2147483648,2147483647",
				"-int8-slice", "-128,127",
				"-int16-slice", "-32768,32767",
				"-int32-slice", "-2147483648,2147483647",
				"-int64-slice", "-9223372036854775808,9223372036854775807",
				"-uint-slice", "0,4294967295",
				"-uint8-slice", "0,255",
				"-uint16-slice", "0,65535",
				"-uint32-slice", "0,4294967295",
				"-uint64-slice", "0,18446744073709551615",
				"-duration-slice", "90m,120m",
				"-url-slice", "service-1:8080,service-2:8080",
			},
			&flags{},
			nil,
			&flags{
				unexported:    "",
				String:        "content",
				Bool:          true,
				Float32:       3.1415,
				Float64:       3.14159265359,
				Int:           -2147483648,
				Int8:          -128,
				Int16:         -32768,
				Int32:         -2147483648,
				Int64:         -9223372036854775808,
				Uint:          4294967295,
				Uint8:         255,
				Uint16:        65535,
				Uint32:        4294967295,
				Uint64:        18446744073709551615,
				Duration:      d90m,
				URL:           *url1,
				StringSlice:   []string{"milad", "mona"},
				BoolSlice:     []bool{false, true},
				Float32Slice:  []float32{3.1415, 2.7182},
				Float64Slice:  []float64{3.14159265359, 2.71828182845},
				IntSlice:      []int{-2147483648, 2147483647},
				Int8Slice:     []int8{-128, 127},
				Int16Slice:    []int16{-32768, 32767},
				Int32Slice:    []int32{-2147483648, 2147483647},
				Int64Slice:    []int64{-9223372036854775808, 9223372036854775807},
				UintSlice:     []uint{0, 4294967295},
				Uint8Slice:    []uint8{0, 255},
				Uint16Slice:   []uint16{0, 65535},
				Uint32Slice:   []uint32{0, 4294967295},
				Uint64Slice:   []uint64{0, 18446744073709551615},
				DurationSlice: []time.Duration{d90m, d120m},
				URLSlice:      []url.URL{*url1, *url2},
			},
		},
		{
			"FromFlags#2",
			[]string{
				"app",
				"--string", "content",
				"--bool",
				"--float32", "3.1415",
				"--float64", "3.14159265359",
				"--int", "-2147483648",
				"--int8", "-128",
				"--int16", "-32768",
				"--int32", "-2147483648",
				"--int64", "-9223372036854775808",
				"--uint", "4294967295",
				"--uint8", "255",
				"--uint16", "65535",
				"--uint32", "4294967295",
				"--uint64", "18446744073709551615",
				"--duration", "90m",
				"--url", "service-1:8080",
				"--string-slice", "milad,mona",
				"--bool-slice", "false,true",
				"--float32-slice", "3.1415,2.7182",
				"--float64-slice", "3.14159265359,2.71828182845",
				"--int-slice", "-2147483648,2147483647",
				"--int8-slice", "-128,127",
				"--int16-slice", "-32768,32767",
				"--int32-slice", "-2147483648,2147483647",
				"--int64-slice", "-9223372036854775808,9223372036854775807",
				"--uint-slice", "0,4294967295",
				"--uint8-slice", "0,255",
				"--uint16-slice", "0,65535",
				"--uint32-slice", "0,4294967295",
				"--uint64-slice", "0,18446744073709551615",
				"--duration-slice", "90m,120m",
				"--url-slice", "service-1:8080,service-2:8080",
			},
			&flags{},
			nil,
			&flags{
				unexported:    "",
				String:        "content",
				Bool:          true,
				Float32:       3.1415,
				Float64:       3.14159265359,
				Int:           -2147483648,
				Int8:          -128,
				Int16:         -32768,
				Int32:         -2147483648,
				Int64:         -9223372036854775808,
				Uint:          4294967295,
				Uint8:         255,
				Uint16:        65535,
				Uint32:        4294967295,
				Uint64:        18446744073709551615,
				Duration:      d90m,
				URL:           *url1,
				StringSlice:   []string{"milad", "mona"},
				BoolSlice:     []bool{false, true},
				Float32Slice:  []float32{3.1415, 2.7182},
				Float64Slice:  []float64{3.14159265359, 2.71828182845},
				IntSlice:      []int{-2147483648, 2147483647},
				Int8Slice:     []int8{-128, 127},
				Int16Slice:    []int16{-32768, 32767},
				Int32Slice:    []int32{-2147483648, 2147483647},
				Int64Slice:    []int64{-9223372036854775808, 9223372036854775807},
				UintSlice:     []uint{0, 4294967295},
				Uint8Slice:    []uint8{0, 255},
				Uint16Slice:   []uint16{0, 65535},
				Uint32Slice:   []uint32{0, 4294967295},
				Uint64Slice:   []uint64{0, 18446744073709551615},
				DurationSlice: []time.Duration{d90m, d120m},
				URLSlice:      []url.URL{*url1, *url2},
			},
		},
		{
			"FromFlags#3",
			[]string{
				"app",
				"-string=content",
				"-bool",
				"-float32=3.1415",
				"-float64=3.14159265359",
				"-int=-2147483648",
				"-int8=-128",
				"-int16=-32768",
				"-int32=-2147483648",
				"-int64=-9223372036854775808",
				"-uint=4294967295",
				"-uint8=255",
				"-uint16=65535",
				"-uint32=4294967295",
				"-uint64=18446744073709551615",
				"-duration=90m",
				"-url=service-1:8080",
				"-string-slice=milad,mona",
				"-bool-slice=false,true",
				"-float32-slice=3.1415,2.7182",
				"-float64-slice=3.14159265359,2.71828182845",
				"-int-slice=-2147483648,2147483647",
				"-int8-slice=-128,127",
				"-int16-slice=-32768,32767",
				"-int32-slice=-2147483648,2147483647",
				"-int64-slice=-9223372036854775808,9223372036854775807",
				"-uint-slice=0,4294967295",
				"-uint8-slice=0,255",
				"-uint16-slice=0,65535",
				"-uint32-slice=0,4294967295",
				"-uint64-slice=0,18446744073709551615",
				"-duration-slice=90m,120m",
				"-url-slice=service-1:8080,service-2:8080",
			},
			&flags{},
			nil,
			&flags{
				unexported:    "",
				String:        "content",
				Bool:          true,
				Float32:       3.1415,
				Float64:       3.14159265359,
				Int:           -2147483648,
				Int8:          -128,
				Int16:         -32768,
				Int32:         -2147483648,
				Int64:         -9223372036854775808,
				Uint:          4294967295,
				Uint8:         255,
				Uint16:        65535,
				Uint32:        4294967295,
				Uint64:        18446744073709551615,
				Duration:      d90m,
				URL:           *url1,
				StringSlice:   []string{"milad", "mona"},
				BoolSlice:     []bool{false, true},
				Float32Slice:  []float32{3.1415, 2.7182},
				Float64Slice:  []float64{3.14159265359, 2.71828182845},
				IntSlice:      []int{-2147483648, 2147483647},
				Int8Slice:     []int8{-128, 127},
				Int16Slice:    []int16{-32768, 32767},
				Int32Slice:    []int32{-2147483648, 2147483647},
				Int64Slice:    []int64{-9223372036854775808, 9223372036854775807},
				UintSlice:     []uint{0, 4294967295},
				Uint8Slice:    []uint8{0, 255},
				Uint16Slice:   []uint16{0, 65535},
				Uint32Slice:   []uint32{0, 4294967295},
				Uint64Slice:   []uint64{0, 18446744073709551615},
				DurationSlice: []time.Duration{d90m, d120m},
				URLSlice:      []url.URL{*url1, *url2},
			},
		},
		{
			"FromFlags#4",
			[]string{
				"app",
				"--string=content",
				"--bool",
				"--float32=3.1415",
				"--float64=3.14159265359",
				"--int=-2147483648",
				"--int8=-128",
				"--int16=-32768",
				"--int32=-2147483648",
				"--int64=-9223372036854775808",
				"--uint=4294967295",
				"--uint8=255",
				"--uint16=65535",
				"--uint32=4294967295",
				"--uint64=18446744073709551615",
				"--duration=90m",
				"--url=service-1:8080",
				"--string-slice=milad,mona",
				"--bool-slice=false,true",
				"--float32-slice=3.1415,2.7182",
				"--float64-slice=3.14159265359,2.71828182845",
				"--int-slice=-2147483648,2147483647",
				"--int8-slice=-128,127",
				"--int16-slice=-32768,32767",
				"--int32-slice=-2147483648,2147483647",
				"--int64-slice=-9223372036854775808,9223372036854775807",
				"--uint-slice=0,4294967295",
				"--uint8-slice=0,255",
				"--uint16-slice=0,65535",
				"--uint32-slice=0,4294967295",
				"--uint64-slice=0,18446744073709551615",
				"--duration-slice=90m,120m",
				"--url-slice=service-1:8080,service-2:8080",
			},
			&flags{},
			nil,
			&flags{
				unexported:    "",
				String:        "content",
				Bool:          true,
				Float32:       3.1415,
				Float64:       3.14159265359,
				Int:           -2147483648,
				Int8:          -128,
				Int16:         -32768,
				Int32:         -2147483648,
				Int64:         -9223372036854775808,
				Uint:          4294967295,
				Uint8:         255,
				Uint16:        65535,
				Uint32:        4294967295,
				Uint64:        18446744073709551615,
				Duration:      d90m,
				URL:           *url1,
				StringSlice:   []string{"milad", "mona"},
				BoolSlice:     []bool{false, true},
				Float32Slice:  []float32{3.1415, 2.7182},
				Float64Slice:  []float64{3.14159265359, 2.71828182845},
				IntSlice:      []int{-2147483648, 2147483647},
				Int8Slice:     []int8{-128, 127},
				Int16Slice:    []int16{-32768, 32767},
				Int32Slice:    []int32{-2147483648, 2147483647},
				Int64Slice:    []int64{-9223372036854775808, 9223372036854775807},
				UintSlice:     []uint{0, 4294967295},
				Uint8Slice:    []uint8{0, 255},
				Uint16Slice:   []uint16{0, 65535},
				Uint32Slice:   []uint32{0, 4294967295},
				Uint64Slice:   []uint64{0, 18446744073709551615},
				DurationSlice: []time.Duration{d90m, d120m},
				URLSlice:      []url.URL{*url1, *url2},
			},
		},
	}

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			os.Args = tc.args

			err := Populate(tc.s)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, tc.s)
			}
		})
	}
}

func TestRegisterFlags(t *testing.T) {
	d90m := 90 * time.Minute
	d120m := 120 * time.Minute

	url1, _ := url.Parse("service-1:8080")
	url2, _ := url.Parse("service-2:8080")

	fs := flag.NewFlagSet("app", flag.ContinueOnError)
	fs.String("string", "", "")

	tests := []struct {
		name          string
		args          []string
		fs            *flag.FlagSet
		s             interface{}
		expectedError error
		expected      *flags
	}{
		{
			"NonStruct",
			[]string{"app"},
			new(flag.FlagSet),
			new(string),
			errors.New("non-struct type: you should pass a pointer to a struct type"),
			&flags{},
		},
		{
			"NonPointer",
			[]string{"app"},
			new(flag.FlagSet),
			flags{},
			errors.New("non-pointer type: you should pass a pointer to a struct type"),
			&flags{},
		},
		{
			"FlagAlreadyRegistered",
			[]string{"app"},
			fs,
			&flags{},
			errors.New("flag already registered: string"),
			&flags{},
		},
		{
			"FromDefaults",
			[]string{"app"},
			new(flag.FlagSet),
			&flags{
				unexported:    "internal",
				String:        "default",
				Bool:          false,
				Float32:       3.1415,
				Float64:       3.14159265359,
				Int:           -2147483648,
				Int8:          -128,
				Int16:         -32768,
				Int32:         -2147483648,
				Int64:         -9223372036854775808,
				Uint:          4294967295,
				Uint8:         255,
				Uint16:        65535,
				Uint32:        4294967295,
				Uint64:        18446744073709551615,
				Duration:      d90m,
				URL:           *url1,
				StringSlice:   []string{"milad", "mona"},
				BoolSlice:     []bool{false, true},
				Float32Slice:  []float32{3.1415, 2.7182},
				Float64Slice:  []float64{3.14159265359, 2.71828182845},
				IntSlice:      []int{-2147483648, 2147483647},
				Int8Slice:     []int8{-128, 127},
				Int16Slice:    []int16{-32768, 32767},
				Int32Slice:    []int32{-2147483648, 2147483647},
				Int64Slice:    []int64{-9223372036854775808, 9223372036854775807},
				UintSlice:     []uint{0, 4294967295},
				Uint8Slice:    []uint8{0, 255},
				Uint16Slice:   []uint16{0, 65535},
				Uint32Slice:   []uint32{0, 4294967295},
				Uint64Slice:   []uint64{0, 18446744073709551615},
				DurationSlice: []time.Duration{d90m, d120m},
				URLSlice:      []url.URL{*url1, *url2},
			},
			nil,
			&flags{
				unexported:    "internal",
				String:        "default",
				Bool:          false,
				Float32:       3.1415,
				Float64:       3.14159265359,
				Int:           -2147483648,
				Int8:          -128,
				Int16:         -32768,
				Int32:         -2147483648,
				Int64:         -9223372036854775808,
				Uint:          4294967295,
				Uint8:         255,
				Uint16:        65535,
				Uint32:        4294967295,
				Uint64:        18446744073709551615,
				Duration:      d90m,
				URL:           *url1,
				StringSlice:   []string{"milad", "mona"},
				BoolSlice:     []bool{false, true},
				Float32Slice:  []float32{3.1415, 2.7182},
				Float64Slice:  []float64{3.14159265359, 2.71828182845},
				IntSlice:      []int{-2147483648, 2147483647},
				Int8Slice:     []int8{-128, 127},
				Int16Slice:    []int16{-32768, 32767},
				Int32Slice:    []int32{-2147483648, 2147483647},
				Int64Slice:    []int64{-9223372036854775808, 9223372036854775807},
				UintSlice:     []uint{0, 4294967295},
				Uint8Slice:    []uint8{0, 255},
				Uint16Slice:   []uint16{0, 65535},
				Uint32Slice:   []uint32{0, 4294967295},
				Uint64Slice:   []uint64{0, 18446744073709551615},
				DurationSlice: []time.Duration{d90m, d120m},
				URLSlice:      []url.URL{*url1, *url2},
			},
		},
		{
			"FromFlags#1",
			[]string{
				"app",
				"-string", "content",
				"-bool",
				"-float32", "3.1415",
				"-float64", "3.14159265359",
				"-int", "-2147483648",
				"-int8", "-128",
				"-int16", "-32768",
				"-int32", "-2147483648",
				"-int64", "-9223372036854775808",
				"-uint", "4294967295",
				"-uint8", "255",
				"-uint16", "65535",
				"-uint32", "4294967295",
				"-uint64", "18446744073709551615",
				"-duration", "90m",
				"-url", "service-1:8080",
				"-string-slice", "milad,mona",
				"-bool-slice", "false,true",
				"-float32-slice", "3.1415,2.7182",
				"-float64-slice", "3.14159265359,2.71828182845",
				"-int-slice", "-2147483648,2147483647",
				"-int8-slice", "-128,127",
				"-int16-slice", "-32768,32767",
				"-int32-slice", "-2147483648,2147483647",
				"-int64-slice", "-9223372036854775808,9223372036854775807",
				"-uint-slice", "0,4294967295",
				"-uint8-slice", "0,255",
				"-uint16-slice", "0,65535",
				"-uint32-slice", "0,4294967295",
				"-uint64-slice", "0,18446744073709551615",
				"-duration-slice", "90m,120m",
				"-url-slice", "service-1:8080,service-2:8080",
			},
			new(flag.FlagSet),
			&flags{},
			nil,
			&flags{
				unexported:    "",
				String:        "content",
				Bool:          true,
				Float32:       3.1415,
				Float64:       3.14159265359,
				Int:           -2147483648,
				Int8:          -128,
				Int16:         -32768,
				Int32:         -2147483648,
				Int64:         -9223372036854775808,
				Uint:          4294967295,
				Uint8:         255,
				Uint16:        65535,
				Uint32:        4294967295,
				Uint64:        18446744073709551615,
				Duration:      d90m,
				URL:           *url1,
				StringSlice:   []string{"milad", "mona"},
				BoolSlice:     []bool{false, true},
				Float32Slice:  []float32{3.1415, 2.7182},
				Float64Slice:  []float64{3.14159265359, 2.71828182845},
				IntSlice:      []int{-2147483648, 2147483647},
				Int8Slice:     []int8{-128, 127},
				Int16Slice:    []int16{-32768, 32767},
				Int32Slice:    []int32{-2147483648, 2147483647},
				Int64Slice:    []int64{-9223372036854775808, 9223372036854775807},
				UintSlice:     []uint{0, 4294967295},
				Uint8Slice:    []uint8{0, 255},
				Uint16Slice:   []uint16{0, 65535},
				Uint32Slice:   []uint32{0, 4294967295},
				Uint64Slice:   []uint64{0, 18446744073709551615},
				DurationSlice: []time.Duration{d90m, d120m},
				URLSlice:      []url.URL{*url1, *url2},
			},
		},
		{
			"FromFlags#2",
			[]string{
				"app",
				"--string", "content",
				"--bool",
				"--float32", "3.1415",
				"--float64", "3.14159265359",
				"--int", "-2147483648",
				"--int8", "-128",
				"--int16", "-32768",
				"--int32", "-2147483648",
				"--int64", "-9223372036854775808",
				"--uint", "4294967295",
				"--uint8", "255",
				"--uint16", "65535",
				"--uint32", "4294967295",
				"--uint64", "18446744073709551615",
				"--duration", "90m",
				"--url", "service-1:8080",
				"--string-slice", "milad,mona",
				"--bool-slice", "false,true",
				"--float32-slice", "3.1415,2.7182",
				"--float64-slice", "3.14159265359,2.71828182845",
				"--int-slice", "-2147483648,2147483647",
				"--int8-slice", "-128,127",
				"--int16-slice", "-32768,32767",
				"--int32-slice", "-2147483648,2147483647",
				"--int64-slice", "-9223372036854775808,9223372036854775807",
				"--uint-slice", "0,4294967295",
				"--uint8-slice", "0,255",
				"--uint16-slice", "0,65535",
				"--uint32-slice", "0,4294967295",
				"--uint64-slice", "0,18446744073709551615",
				"--duration-slice", "90m,120m",
				"--url-slice", "service-1:8080,service-2:8080",
			},
			new(flag.FlagSet),
			&flags{},
			nil,
			&flags{
				unexported:    "",
				String:        "content",
				Bool:          true,
				Float32:       3.1415,
				Float64:       3.14159265359,
				Int:           -2147483648,
				Int8:          -128,
				Int16:         -32768,
				Int32:         -2147483648,
				Int64:         -9223372036854775808,
				Uint:          4294967295,
				Uint8:         255,
				Uint16:        65535,
				Uint32:        4294967295,
				Uint64:        18446744073709551615,
				Duration:      d90m,
				URL:           *url1,
				StringSlice:   []string{"milad", "mona"},
				BoolSlice:     []bool{false, true},
				Float32Slice:  []float32{3.1415, 2.7182},
				Float64Slice:  []float64{3.14159265359, 2.71828182845},
				IntSlice:      []int{-2147483648, 2147483647},
				Int8Slice:     []int8{-128, 127},
				Int16Slice:    []int16{-32768, 32767},
				Int32Slice:    []int32{-2147483648, 2147483647},
				Int64Slice:    []int64{-9223372036854775808, 9223372036854775807},
				UintSlice:     []uint{0, 4294967295},
				Uint8Slice:    []uint8{0, 255},
				Uint16Slice:   []uint16{0, 65535},
				Uint32Slice:   []uint32{0, 4294967295},
				Uint64Slice:   []uint64{0, 18446744073709551615},
				DurationSlice: []time.Duration{d90m, d120m},
				URLSlice:      []url.URL{*url1, *url2},
			},
		},
		{
			"FromFlags#3",
			[]string{
				"app",
				"-string=content",
				"-bool",
				"-float32=3.1415",
				"-float64=3.14159265359",
				"-int=-2147483648",
				"-int8=-128",
				"-int16=-32768",
				"-int32=-2147483648",
				"-int64=-9223372036854775808",
				"-uint=4294967295",
				"-uint8=255",
				"-uint16=65535",
				"-uint32=4294967295",
				"-uint64=18446744073709551615",
				"-duration=90m",
				"-url=service-1:8080",
				"-string-slice=milad,mona",
				"-bool-slice=false,true",
				"-float32-slice=3.1415,2.7182",
				"-float64-slice=3.14159265359,2.71828182845",
				"-int-slice=-2147483648,2147483647",
				"-int8-slice=-128,127",
				"-int16-slice=-32768,32767",
				"-int32-slice=-2147483648,2147483647",
				"-int64-slice=-9223372036854775808,9223372036854775807",
				"-uint-slice=0,4294967295",
				"-uint8-slice=0,255",
				"-uint16-slice=0,65535",
				"-uint32-slice=0,4294967295",
				"-uint64-slice=0,18446744073709551615",
				"-duration-slice=90m,120m",
				"-url-slice=service-1:8080,service-2:8080",
			},
			new(flag.FlagSet),
			&flags{},
			nil,
			&flags{
				unexported:    "",
				String:        "content",
				Bool:          true,
				Float32:       3.1415,
				Float64:       3.14159265359,
				Int:           -2147483648,
				Int8:          -128,
				Int16:         -32768,
				Int32:         -2147483648,
				Int64:         -9223372036854775808,
				Uint:          4294967295,
				Uint8:         255,
				Uint16:        65535,
				Uint32:        4294967295,
				Uint64:        18446744073709551615,
				Duration:      d90m,
				URL:           *url1,
				StringSlice:   []string{"milad", "mona"},
				BoolSlice:     []bool{false, true},
				Float32Slice:  []float32{3.1415, 2.7182},
				Float64Slice:  []float64{3.14159265359, 2.71828182845},
				IntSlice:      []int{-2147483648, 2147483647},
				Int8Slice:     []int8{-128, 127},
				Int16Slice:    []int16{-32768, 32767},
				Int32Slice:    []int32{-2147483648, 2147483647},
				Int64Slice:    []int64{-9223372036854775808, 9223372036854775807},
				UintSlice:     []uint{0, 4294967295},
				Uint8Slice:    []uint8{0, 255},
				Uint16Slice:   []uint16{0, 65535},
				Uint32Slice:   []uint32{0, 4294967295},
				Uint64Slice:   []uint64{0, 18446744073709551615},
				DurationSlice: []time.Duration{d90m, d120m},
				URLSlice:      []url.URL{*url1, *url2},
			},
		},
		{
			"FromFlags#4",
			[]string{
				"app",
				"--string=content",
				"--bool",
				"--float32=3.1415",
				"--float64=3.14159265359",
				"--int=-2147483648",
				"--int8=-128",
				"--int16=-32768",
				"--int32=-2147483648",
				"--int64=-9223372036854775808",
				"--uint=4294967295",
				"--uint8=255",
				"--uint16=65535",
				"--uint32=4294967295",
				"--uint64=18446744073709551615",
				"--duration=90m",
				"--url=service-1:8080",
				"--string-slice=milad,mona",
				"--bool-slice=false,true",
				"--float32-slice=3.1415,2.7182",
				"--float64-slice=3.14159265359,2.71828182845",
				"--int-slice=-2147483648,2147483647",
				"--int8-slice=-128,127",
				"--int16-slice=-32768,32767",
				"--int32-slice=-2147483648,2147483647",
				"--int64-slice=-9223372036854775808,9223372036854775807",
				"--uint-slice=0,4294967295",
				"--uint8-slice=0,255",
				"--uint16-slice=0,65535",
				"--uint32-slice=0,4294967295",
				"--uint64-slice=0,18446744073709551615",
				"--duration-slice=90m,120m",
				"--url-slice=service-1:8080,service-2:8080",
			},
			new(flag.FlagSet),
			&flags{},
			nil,
			&flags{
				unexported:    "",
				String:        "content",
				Bool:          true,
				Float32:       3.1415,
				Float64:       3.14159265359,
				Int:           -2147483648,
				Int8:          -128,
				Int16:         -32768,
				Int32:         -2147483648,
				Int64:         -9223372036854775808,
				Uint:          4294967295,
				Uint8:         255,
				Uint16:        65535,
				Uint32:        4294967295,
				Uint64:        18446744073709551615,
				Duration:      d90m,
				URL:           *url1,
				StringSlice:   []string{"milad", "mona"},
				BoolSlice:     []bool{false, true},
				Float32Slice:  []float32{3.1415, 2.7182},
				Float64Slice:  []float64{3.14159265359, 2.71828182845},
				IntSlice:      []int{-2147483648, 2147483647},
				Int8Slice:     []int8{-128, 127},
				Int16Slice:    []int16{-32768, 32767},
				Int32Slice:    []int32{-2147483648, 2147483647},
				Int64Slice:    []int64{-9223372036854775808, 9223372036854775807},
				UintSlice:     []uint{0, 4294967295},
				Uint8Slice:    []uint8{0, 255},
				Uint16Slice:   []uint16{0, 65535},
				Uint32Slice:   []uint32{0, 4294967295},
				Uint64Slice:   []uint64{0, 18446744073709551615},
				DurationSlice: []time.Duration{d90m, d120m},
				URLSlice:      []url.URL{*url1, *url2},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := RegisterFlags(tc.fs, tc.s)
			tc.fs.Parse(tc.args[1:])

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, tc.s)
			}
		})
	}
}
