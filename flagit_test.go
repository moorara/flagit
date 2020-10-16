package flagit

import (
	"errors"
	"flag"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type (
	Floats struct {
		Float32 float32 `flag:"float32"`
		Float64 float64 `flag:"float64"`
	}

	Ints struct {
		Int   int   `flag:"int"`
		Int8  int8  `flag:"int8"`
		Int16 int16 `flag:"int16"`
		Int32 int32 `flag:"int32"`
		Int64 int64 `flag:"int64"`
	}

	Uints struct {
		Uint   uint   `flag:"uint"`
		Uint8  uint8  `flag:"uint8"`
		Uint16 uint16 `flag:"uint16"`
		Uint32 uint32 `flag:"uint32"`
		Uint64 uint64 `flag:"uint64"`
	}

	FloatSlices struct {
		Float32Slice []float32 `flag:"float32-slice"`
		Float64Slice []float64 `flag:"float64-slice"`
	}

	IntSlices struct {
		IntSlice   []int   `flag:"int-slice"`
		Int8Slice  []int8  `flag:"int8-slice"`
		Int16Slice []int16 `flag:"int16-slice"`
		Int32Slice []int32 `flag:"int32-slice"`
		Int64Slice []int64 `flag:"int64-slice"`
	}

	UintSlices struct {
		UintSlice   []uint   `flag:"uint-slice"`
		Uint8Slice  []uint8  `flag:"uint8-slice"`
		Uint16Slice []uint16 `flag:"uint16-slice"`
		Uint32Slice []uint32 `flag:"uint32-slice"`
		Uint64Slice []uint64 `flag:"uint64-slice"`
	}

	SliceGroup struct {
		StringSlice   []string        `flag:"string-slice,the help text for the string-slice flag"`
		BoolSlice     []bool          `flag:"bool-slice,the help text for the bool-slice flag"`
		FloatSlices   FloatSlices     `flag:""`
		IntSlices     IntSlices       `flag:""`
		UintSlices    UintSlices      `flag:""`
		DurationSlice []time.Duration `flag:"duration-slice,the help text for the duration-slice flag"`
		URLSlice      []url.URL       `flag:"url-slice,the help text for the url-slice flag"`
		RegexpSlice   []regexp.Regexp `flag:"regexp-slice,the help text for the regexp-slice flag"`
	}

	Flags struct {
		unexported  string
		WithoutFlag string
		String      string        `flag:"string,the help text for the string flag"`
		Bool        bool          `flag:"bool,the help text for the bool flag"`
		Floats      Floats        `flag:""`
		Ints        Ints          `flag:""`
		Uints       Uints         `flag:""`
		Duration    time.Duration `flag:"duration,the help text for the duration flag"`
		URL         url.URL       `flag:"url,the help text for the url flag"`
		Regexp      regexp.Regexp `flag:"regexp,the help text for the regexp flag"`
		SliceGroup  SliceGroup    `flag:""`
	}
)

func getFields(vStruct reflect.Value, handle func(f fieldInfo)) {
	for i := 0; i < vStruct.NumField(); i++ {
		v := vStruct.Field(i)
		t := v.Type()
		f := vStruct.Type().Field(i)

		if isNestedStruct(t) {
			if _, ok := f.Tag.Lookup(flagTag); ok {
				getFields(v, handle)
			}
		}

		if !v.CanSet() || !isTypeSupported(t) && f.Tag.Get(flagTag) != "" {
			continue
		}

		handle(fieldInfo{
			value: v,
			name:  f.Name,
			sep:   ",",
		})
	}
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
			Flags{},
			errors.New("non-pointer type: you should pass a pointer to a struct type"),
		},
		{
			"OK",
			&Flags{},
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

func TestIsStructSupported(t *testing.T) {
	tests := []struct {
		name     string
		s        interface{}
		expected bool
	}{
		{
			name:     "NotSupported",
			s:        struct{}{},
			expected: false,
		},
		{
			name:     "URL",
			s:        url.URL{},
			expected: true,
		},
		{
			name:     "Regexp",
			s:        regexp.Regexp{},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tStruct := reflect.TypeOf(tc.s)

			assert.Equal(t, tc.expected, isStructSupported(tStruct))
		})
	}
}

func TestIsNestedStruct(t *testing.T) {
	vStruct := reflect.ValueOf(struct {
		Int    int
		URL    url.URL
		Regexp regexp.Regexp
		Group  struct {
			String string
		}
	}{})

	vInt := vStruct.FieldByName("Int")
	assert.False(t, isNestedStruct(vInt.Type()))

	vURL := vStruct.FieldByName("URL")
	assert.False(t, isNestedStruct(vURL.Type()))

	vRegexp := vStruct.FieldByName("Regexp")
	assert.False(t, isNestedStruct(vRegexp.Type()))

	vGroup := vStruct.FieldByName("Group")
	assert.True(t, isNestedStruct(vGroup.Type()))
}

func TestIsTypeSupported(t *testing.T) {
	url1, _ := url.Parse("service-1:8080")
	url2, _ := url.Parse("service-2:8080")

	re1 := regexp.MustCompilePOSIX("[:digit:]")
	re2 := regexp.MustCompilePOSIX("[:alpha:]")

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
		{"Regexp", *re1, true},
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
		{"RegexpSlice", []regexp.Regexp{*re1, *re2}, true},
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
		flag              string
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
		flagValue := getFlagValue(tc.flag)

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

	re1 := regexp.MustCompilePOSIX("[:digit:]")
	re2 := regexp.MustCompilePOSIX("[:alpha:]")
	re3 := regexp.MustCompilePOSIX("[:alnum:]")
	re4 := regexp.MustCompilePOSIX("[:word:]")

	tests := []struct {
		name            string
		s               *Flags
		values          map[string]string
		expectedUpdated bool
		expectError     bool
		expected        *Flags
	}{
		{
			"InvalidValues",
			&Flags{},
			map[string]string{
				"Bool":          "invalid",
				"Float32":       "invalid",
				"Float64":       "invalid",
				"Int":           "invalid",
				"Int8":          "invalid",
				"Int16":         "invalid",
				"Int32":         "invalid",
				"Int64":         "invalid",
				"Uint":          "invalid",
				"Uint8":         "invalid",
				"Uint16":        "invalid",
				"Uint32":        "invalid",
				"Uint64":        "invalid",
				"Duration":      "invalid",
				"URL":           ":invalid",
				"Regexp":        "[:invalid",
				"BoolSlice":     "invalid",
				"Float32Slice":  "invalid",
				"Float64Slice":  "invalid",
				"IntSlice":      "invalid",
				"Int8Slice":     "invalid",
				"Int16Slice":    "invalid",
				"Int32Slice":    "invalid",
				"Int64Slice":    "invalid",
				"UintSlice":     "invalid",
				"Uint8Slice":    "invalid",
				"Uint16Slice":   "invalid",
				"Uint32Slice":   "invalid",
				"Uint64Slice":   "invalid",
				"DurationSlice": "invalid",
				"URLSlice":      ":invalid",
				"RegexpSlice":   "[:invalid",
			},
			false,
			true,
			&Flags{},
		},
		{
			"NewValues",
			&Flags{
				String: "default",
				Bool:   false,
				Floats: Floats{
					Float32: 3.1415,
					Float64: 3.14159265359,
				},
				Ints: Ints{
					Int:   -2147483648,
					Int8:  -128,
					Int16: -32768,
					Int32: -2147483648,
					Int64: -9223372036854775808,
				},
				Uints: Uints{
					Uint:   4294967295,
					Uint8:  255,
					Uint16: 65535,
					Uint32: 4294967295,
					Uint64: 18446744073709551615,
				},
				Duration: d90m,
				URL:      *url1,
				Regexp:   *re1,
				SliceGroup: SliceGroup{
					StringSlice: []string{"milad", "mona"},
					BoolSlice:   []bool{false, true},
					FloatSlices: FloatSlices{
						Float32Slice: []float32{3.1415, 2.7182},
						Float64Slice: []float64{3.14159265359, 2.71828182845},
					},
					IntSlices: IntSlices{
						IntSlice:   []int{-2147483648, 2147483647},
						Int8Slice:  []int8{-128, 127},
						Int16Slice: []int16{-32768, 32767},
						Int32Slice: []int32{-2147483648, 2147483647},
						Int64Slice: []int64{-9223372036854775808, 9223372036854775807},
					},
					UintSlices: UintSlices{
						UintSlice:   []uint{0, 4294967295},
						Uint8Slice:  []uint8{0, 255},
						Uint16Slice: []uint16{0, 65535},
						Uint32Slice: []uint32{0, 4294967295},
						Uint64Slice: []uint64{0, 18446744073709551615},
					},
					DurationSlice: []time.Duration{d90m, d120m},
					URLSlice:      []url.URL{*url1, *url2},
					RegexpSlice:   []regexp.Regexp{*re1, *re2},
				},
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
				"Regexp":        "[:alnum:]",
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
				"RegexpSlice":   "[:alnum:],[:word:]",
			},
			true,
			false,
			&Flags{
				String: "content",
				Bool:   true,
				Floats: Floats{
					Float32: 2.7182,
					Float64: 2.7182818284,
				},
				Ints: Ints{
					Int:   2147483647,
					Int8:  127,
					Int16: 32767,
					Int32: 2147483647,
					Int64: 9223372036854775807,
				},
				Uints: Uints{
					Uint:   2147483648,
					Uint8:  128,
					Uint16: 32768,
					Uint32: 2147483648,
					Uint64: 9223372036854775808,
				},
				Duration: d4h,
				URL:      *url3,
				Regexp:   *re3,
				SliceGroup: SliceGroup{
					StringSlice: []string{"mona", "milad"},
					BoolSlice:   []bool{true, false},
					FloatSlices: FloatSlices{
						Float32Slice: []float32{2.7182, 3.1415},
						Float64Slice: []float64{2.71828182845, 3.14159265359},
					},
					IntSlices: IntSlices{
						IntSlice:   []int{2147483647, -2147483648},
						Int8Slice:  []int8{127, -128},
						Int16Slice: []int16{32767, -32768},
						Int32Slice: []int32{2147483647, -2147483648},
						Int64Slice: []int64{9223372036854775807, -9223372036854775808},
					},
					UintSlices: UintSlices{
						UintSlice:   []uint{4294967295, 0},
						Uint8Slice:  []uint8{255, 0},
						Uint16Slice: []uint16{65535, 0},
						Uint32Slice: []uint32{4294967295, 0},
						Uint64Slice: []uint64{18446744073709551615, 0},
					},
					DurationSlice: []time.Duration{d4h, d8h},
					URLSlice:      []url.URL{*url3, *url4},
					RegexpSlice:   []regexp.Regexp{*re3, *re4},
				},
			},
		},
		{
			"NoNewValues",
			&Flags{
				String: "content",
				Bool:   true,
				Floats: Floats{
					Float32: 2.7182,
					Float64: 2.7182818284,
				},
				Ints: Ints{
					Int:   2147483647,
					Int8:  127,
					Int16: 32767,
					Int32: 2147483647,
					Int64: 9223372036854775807,
				},
				Uints: Uints{
					Uint:   2147483648,
					Uint8:  128,
					Uint16: 32768,
					Uint32: 2147483648,
					Uint64: 9223372036854775808,
				},
				Duration: d4h,
				URL:      *url3,
				Regexp:   *re3,
				SliceGroup: SliceGroup{
					StringSlice: []string{"mona", "milad"},
					BoolSlice:   []bool{true, false},
					FloatSlices: FloatSlices{
						Float32Slice: []float32{2.7182, 3.1415},
						Float64Slice: []float64{2.71828182845, 3.14159265359},
					},
					IntSlices: IntSlices{
						IntSlice:   []int{2147483647, -2147483648},
						Int8Slice:  []int8{127, -128},
						Int16Slice: []int16{32767, -32768},
						Int32Slice: []int32{2147483647, -2147483648},
						Int64Slice: []int64{9223372036854775807, -9223372036854775808},
					},
					UintSlices: UintSlices{
						UintSlice:   []uint{4294967295, 0},
						Uint8Slice:  []uint8{255, 0},
						Uint16Slice: []uint16{65535, 0},
						Uint32Slice: []uint32{4294967295, 0},
						Uint64Slice: []uint64{18446744073709551615, 0},
					},
					DurationSlice: []time.Duration{d4h, d8h},
					URLSlice:      []url.URL{*url3, *url4},
					RegexpSlice:   []regexp.Regexp{*re3, *re4},
				},
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
				"Regexp":        "[:alnum:]",
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
				"RegexpSlice":   "[:alnum:],[:word:]",
			},
			false,
			false,
			&Flags{
				String: "content",
				Bool:   true,
				Floats: Floats{
					Float32: 2.7182,
					Float64: 2.7182818284,
				},
				Ints: Ints{
					Int:   2147483647,
					Int8:  127,
					Int16: 32767,
					Int32: 2147483647,
					Int64: 9223372036854775807,
				},
				Uints: Uints{
					Uint:   2147483648,
					Uint8:  128,
					Uint16: 32768,
					Uint32: 2147483648,
					Uint64: 9223372036854775808,
				},
				Duration: d4h,
				URL:      *url3,
				Regexp:   *re3,
				SliceGroup: SliceGroup{
					StringSlice: []string{"mona", "milad"},
					BoolSlice:   []bool{true, false},
					FloatSlices: FloatSlices{
						Float32Slice: []float32{2.7182, 3.1415},
						Float64Slice: []float64{2.71828182845, 3.14159265359},
					},
					IntSlices: IntSlices{
						IntSlice:   []int{2147483647, -2147483648},
						Int8Slice:  []int8{127, -128},
						Int16Slice: []int16{32767, -32768},
						Int32Slice: []int32{2147483647, -2147483648},
						Int64Slice: []int64{9223372036854775807, -9223372036854775808},
					},
					UintSlices: UintSlices{
						UintSlice:   []uint{4294967295, 0},
						Uint8Slice:  []uint8{255, 0},
						Uint16Slice: []uint16{65535, 0},
						Uint32Slice: []uint32{4294967295, 0},
						Uint64Slice: []uint64{18446744073709551615, 0},
					},
					DurationSlice: []time.Duration{d4h, d8h},
					URLSlice:      []url.URL{*url3, *url4},
					RegexpSlice:   []regexp.Regexp{*re3, *re4},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vStruct := reflect.ValueOf(tc.s).Elem()
			getFields(vStruct, func(f fieldInfo) {
				if val := tc.values[f.name]; val != "" {
					updated, err := setFieldValue(f, val)

					if tc.expectError {
						assert.Error(t, err)
					} else {
						assert.NoError(t, err)
						assert.Equal(t, tc.expectedUpdated, updated)
					}
				}
			})

			assert.Equal(t, tc.expected, tc.s)
		})
	}
}

func TestIterateOnFields(t *testing.T) {
	invalid := struct {
		LogLevel string `flag:"log level"`
	}{}

	tests := []struct {
		name               string
		s                  interface{}
		continueOnError    bool
		expectedError      error
		expectedFieldNames []string
		expectedFlagNames  []string
		expectedListSeps   []string
	}{
		{
			name:               "StopOnError",
			s:                  &invalid,
			continueOnError:    false,
			expectedError:      errors.New("invalid flag name: log level"),
			expectedFieldNames: []string{},
			expectedFlagNames:  []string{},
			expectedListSeps:   []string{},
		},
		{
			name:               "ContinueOnError",
			s:                  &invalid,
			continueOnError:    true,
			expectedError:      nil,
			expectedFieldNames: []string{},
			expectedFlagNames:  []string{},
			expectedListSeps:   []string{},
		},
		{
			name:            "OK",
			s:               &Flags{},
			continueOnError: false,
			expectedError:   nil,
			expectedFieldNames: []string{
				"String",
				"Bool",
				"Float32", "Float64",
				"Int", "Int8", "Int16", "Int32", "Int64",
				"Uint", "Uint8", "Uint16", "Uint32", "Uint64",
				"Duration", "URL", "Regexp",
				"StringSlice",
				"BoolSlice",
				"Float32Slice", "Float64Slice",
				"IntSlice", "Int8Slice", "Int16Slice", "Int32Slice", "Int64Slice",
				"UintSlice", "Uint8Slice", "Uint16Slice", "Uint32Slice", "Uint64Slice",
				"DurationSlice", "URLSlice", "RegexpSlice",
			},
			expectedFlagNames: []string{
				"string",
				"bool",
				"float32", "float64",
				"int", "int8", "int16", "int32", "int64",
				"uint", "uint8", "uint16", "uint32", "uint64",
				"duration", "url", "regexp",
				"string-slice",
				"bool-slice",
				"float32-slice", "float64-slice",
				"int-slice", "int8-slice", "int16-slice", "int32-slice", "int64-slice",
				"uint-slice", "uint8-slice", "uint16-slice", "uint32-slice", "uint64-slice",
				"duration-slice", "url-slice", "regexp-slice",
			},
			expectedListSeps: []string{
				",",
				",",
				",", ",",
				",", ",", ",", ",", ",",
				",", ",", ",", ",", ",",
				",", ",", ",",
				",",
				",",
				",", ",",
				",", ",", ",", ",", ",",
				",", ",", ",", ",", ",",
				",", ",", ",",
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

			err = iterateOnFields("", vStruct, tc.continueOnError, func(f fieldInfo) error {
				fieldNames = append(fieldNames, f.name)
				flagNames = append(flagNames, f.flag)
				listSeps = append(listSeps, f.sep)
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

	re1 := regexp.MustCompilePOSIX("[:digit:]")
	re2 := regexp.MustCompilePOSIX("[:alpha:]")

	tests := []struct {
		name            string
		args            []string
		s               interface{}
		continueOnError bool
		expectedError   string
		expected        *Flags
	}{
		{
			"NonStruct",
			[]string{"app"},
			new(string),
			false,
			"non-struct type: you should pass a pointer to a struct type",
			&Flags{},
		},
		{
			"NonPointer",
			[]string{"app"},
			Flags{},
			false,
			"non-pointer type: you should pass a pointer to a struct type",
			&Flags{},
		},
		{
			"FromDefaults",
			[]string{"app"},
			&Flags{
				unexported: "internal",
				String:     "default",
				Bool:       false,
				Floats: Floats{
					Float32: 3.1415,
					Float64: 3.14159265359,
				},
				Ints: Ints{
					Int:   -2147483648,
					Int8:  -128,
					Int16: -32768,
					Int32: -2147483648,
					Int64: -9223372036854775808,
				},
				Uints: Uints{
					Uint:   4294967295,
					Uint8:  255,
					Uint16: 65535,
					Uint32: 4294967295,
					Uint64: 18446744073709551615,
				},
				Duration: d90m,
				URL:      *url1,
				Regexp:   *re1,
				SliceGroup: SliceGroup{
					StringSlice: []string{"milad", "mona"},
					BoolSlice:   []bool{false, true},
					FloatSlices: FloatSlices{
						Float32Slice: []float32{3.1415, 2.7182},
						Float64Slice: []float64{3.14159265359, 2.71828182845},
					},
					IntSlices: IntSlices{
						IntSlice:   []int{-2147483648, 2147483647},
						Int8Slice:  []int8{-128, 127},
						Int16Slice: []int16{-32768, 32767},
						Int32Slice: []int32{-2147483648, 2147483647},
						Int64Slice: []int64{-9223372036854775808, 9223372036854775807},
					},
					UintSlices: UintSlices{
						UintSlice:   []uint{0, 4294967295},
						Uint8Slice:  []uint8{0, 255},
						Uint16Slice: []uint16{0, 65535},
						Uint32Slice: []uint32{0, 4294967295},
						Uint64Slice: []uint64{0, 18446744073709551615},
					},
					DurationSlice: []time.Duration{d90m, d120m},
					URLSlice:      []url.URL{*url1, *url2},
					RegexpSlice:   []regexp.Regexp{*re1, *re2},
				},
			},
			false,
			"",
			&Flags{
				unexported: "internal",
				String:     "default",
				Bool:       false,
				Floats: Floats{
					Float32: 3.1415,
					Float64: 3.14159265359,
				},
				Ints: Ints{
					Int:   -2147483648,
					Int8:  -128,
					Int16: -32768,
					Int32: -2147483648,
					Int64: -9223372036854775808,
				},
				Uints: Uints{
					Uint:   4294967295,
					Uint8:  255,
					Uint16: 65535,
					Uint32: 4294967295,
					Uint64: 18446744073709551615,
				},
				Duration: d90m,
				URL:      *url1,
				Regexp:   *re1,
				SliceGroup: SliceGroup{
					StringSlice: []string{"milad", "mona"},
					BoolSlice:   []bool{false, true},
					FloatSlices: FloatSlices{
						Float32Slice: []float32{3.1415, 2.7182},
						Float64Slice: []float64{3.14159265359, 2.71828182845},
					},
					IntSlices: IntSlices{
						IntSlice:   []int{-2147483648, 2147483647},
						Int8Slice:  []int8{-128, 127},
						Int16Slice: []int16{-32768, 32767},
						Int32Slice: []int32{-2147483648, 2147483647},
						Int64Slice: []int64{-9223372036854775808, 9223372036854775807},
					},
					UintSlices: UintSlices{
						UintSlice:   []uint{0, 4294967295},
						Uint8Slice:  []uint8{0, 255},
						Uint16Slice: []uint16{0, 65535},
						Uint32Slice: []uint32{0, 4294967295},
						Uint64Slice: []uint64{0, 18446744073709551615},
					},
					DurationSlice: []time.Duration{d90m, d120m},
					URLSlice:      []url.URL{*url1, *url2},
					RegexpSlice:   []regexp.Regexp{*re1, *re2},
				},
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
				"-regexp", "[:digit:]",
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
				"-regexp-slice", "[:digit:],[:alpha:]",
			},
			&Flags{},
			false,
			"",
			&Flags{
				unexported: "",
				String:     "content",
				Bool:       true,
				Floats: Floats{
					Float32: 3.1415,
					Float64: 3.14159265359,
				},
				Ints: Ints{
					Int:   -2147483648,
					Int8:  -128,
					Int16: -32768,
					Int32: -2147483648,
					Int64: -9223372036854775808,
				},
				Uints: Uints{
					Uint:   4294967295,
					Uint8:  255,
					Uint16: 65535,
					Uint32: 4294967295,
					Uint64: 18446744073709551615,
				},
				Duration: d90m,
				URL:      *url1,
				Regexp:   *re1,
				SliceGroup: SliceGroup{
					StringSlice: []string{"milad", "mona"},
					BoolSlice:   []bool{false, true},
					FloatSlices: FloatSlices{
						Float32Slice: []float32{3.1415, 2.7182},
						Float64Slice: []float64{3.14159265359, 2.71828182845},
					},
					IntSlices: IntSlices{
						IntSlice:   []int{-2147483648, 2147483647},
						Int8Slice:  []int8{-128, 127},
						Int16Slice: []int16{-32768, 32767},
						Int32Slice: []int32{-2147483648, 2147483647},
						Int64Slice: []int64{-9223372036854775808, 9223372036854775807},
					},
					UintSlices: UintSlices{
						UintSlice:   []uint{0, 4294967295},
						Uint8Slice:  []uint8{0, 255},
						Uint16Slice: []uint16{0, 65535},
						Uint32Slice: []uint32{0, 4294967295},
						Uint64Slice: []uint64{0, 18446744073709551615},
					},
					DurationSlice: []time.Duration{d90m, d120m},
					URLSlice:      []url.URL{*url1, *url2},
					RegexpSlice:   []regexp.Regexp{*re1, *re2},
				},
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
				"--regexp", "[:digit:]",
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
				"--regexp-slice", "[:digit:],[:alpha:]",
			},
			&Flags{},
			false,
			"",
			&Flags{
				unexported: "",
				String:     "content",
				Bool:       true,
				Floats: Floats{
					Float32: 3.1415,
					Float64: 3.14159265359,
				},
				Ints: Ints{
					Int:   -2147483648,
					Int8:  -128,
					Int16: -32768,
					Int32: -2147483648,
					Int64: -9223372036854775808,
				},
				Uints: Uints{
					Uint:   4294967295,
					Uint8:  255,
					Uint16: 65535,
					Uint32: 4294967295,
					Uint64: 18446744073709551615,
				},
				Duration: d90m,
				URL:      *url1,
				Regexp:   *re1,
				SliceGroup: SliceGroup{
					StringSlice: []string{"milad", "mona"},
					BoolSlice:   []bool{false, true},
					FloatSlices: FloatSlices{
						Float32Slice: []float32{3.1415, 2.7182},
						Float64Slice: []float64{3.14159265359, 2.71828182845},
					},
					IntSlices: IntSlices{
						IntSlice:   []int{-2147483648, 2147483647},
						Int8Slice:  []int8{-128, 127},
						Int16Slice: []int16{-32768, 32767},
						Int32Slice: []int32{-2147483648, 2147483647},
						Int64Slice: []int64{-9223372036854775808, 9223372036854775807},
					},
					UintSlices: UintSlices{
						UintSlice:   []uint{0, 4294967295},
						Uint8Slice:  []uint8{0, 255},
						Uint16Slice: []uint16{0, 65535},
						Uint32Slice: []uint32{0, 4294967295},
						Uint64Slice: []uint64{0, 18446744073709551615},
					},
					DurationSlice: []time.Duration{d90m, d120m},
					URLSlice:      []url.URL{*url1, *url2},
					RegexpSlice:   []regexp.Regexp{*re1, *re2},
				},
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
				"-regexp=[:digit:]",
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
				"-regexp-slice=[:digit:],[:alpha:]",
			},
			&Flags{},
			false,
			"",
			&Flags{
				unexported: "",
				String:     "content",
				Bool:       true,
				Floats: Floats{
					Float32: 3.1415,
					Float64: 3.14159265359,
				},
				Ints: Ints{
					Int:   -2147483648,
					Int8:  -128,
					Int16: -32768,
					Int32: -2147483648,
					Int64: -9223372036854775808,
				},
				Uints: Uints{
					Uint:   4294967295,
					Uint8:  255,
					Uint16: 65535,
					Uint32: 4294967295,
					Uint64: 18446744073709551615,
				},
				Duration: d90m,
				URL:      *url1,
				Regexp:   *re1,
				SliceGroup: SliceGroup{
					StringSlice: []string{"milad", "mona"},
					BoolSlice:   []bool{false, true},
					FloatSlices: FloatSlices{
						Float32Slice: []float32{3.1415, 2.7182},
						Float64Slice: []float64{3.14159265359, 2.71828182845},
					},
					IntSlices: IntSlices{
						IntSlice:   []int{-2147483648, 2147483647},
						Int8Slice:  []int8{-128, 127},
						Int16Slice: []int16{-32768, 32767},
						Int32Slice: []int32{-2147483648, 2147483647},
						Int64Slice: []int64{-9223372036854775808, 9223372036854775807},
					},
					UintSlices: UintSlices{
						UintSlice:   []uint{0, 4294967295},
						Uint8Slice:  []uint8{0, 255},
						Uint16Slice: []uint16{0, 65535},
						Uint32Slice: []uint32{0, 4294967295},
						Uint64Slice: []uint64{0, 18446744073709551615},
					},
					DurationSlice: []time.Duration{d90m, d120m},
					URLSlice:      []url.URL{*url1, *url2},
					RegexpSlice:   []regexp.Regexp{*re1, *re2},
				},
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
				"--regexp=[:digit:]",
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
				"--regexp-slice=[:digit:],[:alpha:]",
			},
			&Flags{},
			false,
			"",
			&Flags{
				unexported: "",
				String:     "content",
				Bool:       true,
				Floats: Floats{
					Float32: 3.1415,
					Float64: 3.14159265359,
				},
				Ints: Ints{
					Int:   -2147483648,
					Int8:  -128,
					Int16: -32768,
					Int32: -2147483648,
					Int64: -9223372036854775808,
				},
				Uints: Uints{
					Uint:   4294967295,
					Uint8:  255,
					Uint16: 65535,
					Uint32: 4294967295,
					Uint64: 18446744073709551615,
				},
				Duration: d90m,
				URL:      *url1,
				Regexp:   *re1,
				SliceGroup: SliceGroup{
					StringSlice: []string{"milad", "mona"},
					BoolSlice:   []bool{false, true},
					FloatSlices: FloatSlices{
						Float32Slice: []float32{3.1415, 2.7182},
						Float64Slice: []float64{3.14159265359, 2.71828182845},
					},
					IntSlices: IntSlices{
						IntSlice:   []int{-2147483648, 2147483647},
						Int8Slice:  []int8{-128, 127},
						Int16Slice: []int16{-32768, 32767},
						Int32Slice: []int32{-2147483648, 2147483647},
						Int64Slice: []int64{-9223372036854775808, 9223372036854775807},
					},
					UintSlices: UintSlices{
						UintSlice:   []uint{0, 4294967295},
						Uint8Slice:  []uint8{0, 255},
						Uint16Slice: []uint16{0, 65535},
						Uint32Slice: []uint32{0, 4294967295},
						Uint64Slice: []uint64{0, 18446744073709551615},
					},
					DurationSlice: []time.Duration{d90m, d120m},
					URLSlice:      []url.URL{*url1, *url2},
					RegexpSlice:   []regexp.Regexp{*re1, *re2},
				},
			},
		},
		{
			"StopOnError",
			[]string{
				"app",
				"-int", "invalid",
			},
			&Flags{},
			false,
			`strconv.ParseInt: parsing "invalid": invalid syntax`,
			&Flags{},
		},
		{
			"ContinueOnError",
			[]string{
				"app",
				"-bool", "invalid",
				"-float32", "invalid",
				"-float64", "invalid",
				"-int", "invalid",
				"-int8", "invalid",
				"-int16", "invalid",
				"-int32", "invalid",
				"-int64", "invalid",
				"-uint", "invalid",
				"-uint8", "invalid",
				"-uint16", "invalid",
				"-uint32", "invalid",
				"-uint64", "invalid",
				"-duration", "invalid",
				"-url", ":invalid",
				"-regexp", "[:invalid",
				"-bool-slice", "invalid",
				"-float32-slice", "invalid",
				"-float64-slice", "invalid",
				"-int-slice", "invalid",
				"-int8-slice", "invalid",
				"-int16-slice", "invalid",
				"-int32-slice", "invalid",
				"-int64-slice", "invalid",
				"-uint-slice", "invalid",
				"-uint8-slice", "invalid",
				"-uint16-slice", "invalid",
				"-uint32-slice", "invalid",
				"-uint64-slice", "invalid",
				"-duration-slice", "invalid",
				"-url-slice", ":invalid",
				"-regexp-slice", "[:invalid",
			},
			&Flags{},
			true,
			"",
			&Flags{},
		},
	}

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			os.Args = tc.args

			err := Populate(tc.s, tc.continueOnError)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, tc.s)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestRegisterFlags(t *testing.T) {
	d90m := 90 * time.Minute
	d120m := 120 * time.Minute

	url1, _ := url.Parse("service-1:8080")
	url2, _ := url.Parse("service-2:8080")

	re1 := regexp.MustCompilePOSIX("[:digit:]")
	re2 := regexp.MustCompilePOSIX("[:alpha:]")

	fs := flag.NewFlagSet("app", flag.ContinueOnError)
	fs.String("string", "", "")

	tests := []struct {
		name               string
		args               []string
		fs                 *flag.FlagSet
		s                  interface{}
		continueOnError    bool
		expectedError      error
		expectedParseError string
		expected           *Flags
	}{
		{
			"NonStruct",
			[]string{"app"},
			new(flag.FlagSet),
			new(string),
			false,
			errors.New("non-struct type: you should pass a pointer to a struct type"), "",
			&Flags{},
		},
		{
			"NonPointer",
			[]string{"app"},
			new(flag.FlagSet),
			Flags{},
			false,
			errors.New("non-pointer type: you should pass a pointer to a struct type"), "",
			&Flags{},
		},
		{
			"FlagRegistered_StopOnError",
			[]string{"app"},
			fs,
			&Flags{},
			false,
			errors.New("flag already registered: string"), "",
			&Flags{},
		},
		{
			"FlagRegistered_ContinueOnError",
			[]string{"app"},
			fs,
			&Flags{},
			true,
			nil, "",
			&Flags{},
		},
		{
			"FromDefaults",
			[]string{"app"},
			new(flag.FlagSet),
			&Flags{
				unexported: "internal",
				String:     "default",
				Bool:       false,
				Floats: Floats{
					Float32: 3.1415,
					Float64: 3.14159265359,
				},
				Ints: Ints{
					Int:   -2147483648,
					Int8:  -128,
					Int16: -32768,
					Int32: -2147483648,
					Int64: -9223372036854775808,
				},
				Uints: Uints{
					Uint:   4294967295,
					Uint8:  255,
					Uint16: 65535,
					Uint32: 4294967295,
					Uint64: 18446744073709551615,
				},
				Duration: d90m,
				URL:      *url1,
				Regexp:   *re1,
				SliceGroup: SliceGroup{
					StringSlice: []string{"milad", "mona"},
					BoolSlice:   []bool{false, true},
					FloatSlices: FloatSlices{
						Float32Slice: []float32{3.1415, 2.7182},
						Float64Slice: []float64{3.14159265359, 2.71828182845},
					},
					IntSlices: IntSlices{
						IntSlice:   []int{-2147483648, 2147483647},
						Int8Slice:  []int8{-128, 127},
						Int16Slice: []int16{-32768, 32767},
						Int32Slice: []int32{-2147483648, 2147483647},
						Int64Slice: []int64{-9223372036854775808, 9223372036854775807},
					},
					UintSlices: UintSlices{
						UintSlice:   []uint{0, 4294967295},
						Uint8Slice:  []uint8{0, 255},
						Uint16Slice: []uint16{0, 65535},
						Uint32Slice: []uint32{0, 4294967295},
						Uint64Slice: []uint64{0, 18446744073709551615},
					},
					DurationSlice: []time.Duration{d90m, d120m},
					URLSlice:      []url.URL{*url1, *url2},
					RegexpSlice:   []regexp.Regexp{*re1, *re2},
				},
			},
			false,
			nil, "",
			&Flags{
				unexported: "internal",
				String:     "default",
				Bool:       false,
				Floats: Floats{
					Float32: 3.1415,
					Float64: 3.14159265359,
				},
				Ints: Ints{
					Int:   -2147483648,
					Int8:  -128,
					Int16: -32768,
					Int32: -2147483648,
					Int64: -9223372036854775808,
				},
				Uints: Uints{
					Uint:   4294967295,
					Uint8:  255,
					Uint16: 65535,
					Uint32: 4294967295,
					Uint64: 18446744073709551615,
				},
				Duration: d90m,
				URL:      *url1,
				Regexp:   *re1,
				SliceGroup: SliceGroup{
					StringSlice: []string{"milad", "mona"},
					BoolSlice:   []bool{false, true},
					FloatSlices: FloatSlices{
						Float32Slice: []float32{3.1415, 2.7182},
						Float64Slice: []float64{3.14159265359, 2.71828182845},
					},
					IntSlices: IntSlices{
						IntSlice:   []int{-2147483648, 2147483647},
						Int8Slice:  []int8{-128, 127},
						Int16Slice: []int16{-32768, 32767},
						Int32Slice: []int32{-2147483648, 2147483647},
						Int64Slice: []int64{-9223372036854775808, 9223372036854775807},
					},
					UintSlices: UintSlices{
						UintSlice:   []uint{0, 4294967295},
						Uint8Slice:  []uint8{0, 255},
						Uint16Slice: []uint16{0, 65535},
						Uint32Slice: []uint32{0, 4294967295},
						Uint64Slice: []uint64{0, 18446744073709551615},
					},
					DurationSlice: []time.Duration{d90m, d120m},
					URLSlice:      []url.URL{*url1, *url2},
					RegexpSlice:   []regexp.Regexp{*re1, *re2},
				},
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
				"-regexp", "[:digit:]",
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
				"-regexp-slice", "[:digit:],[:alpha:]",
			},
			new(flag.FlagSet),
			&Flags{},
			false,
			nil, "",
			&Flags{
				unexported: "",
				String:     "content",
				Bool:       true,
				Floats: Floats{
					Float32: 3.1415,
					Float64: 3.14159265359,
				},
				Ints: Ints{
					Int:   -2147483648,
					Int8:  -128,
					Int16: -32768,
					Int32: -2147483648,
					Int64: -9223372036854775808,
				},
				Uints: Uints{
					Uint:   4294967295,
					Uint8:  255,
					Uint16: 65535,
					Uint32: 4294967295,
					Uint64: 18446744073709551615,
				},
				Duration: d90m,
				URL:      *url1,
				Regexp:   *re1,
				SliceGroup: SliceGroup{
					StringSlice: []string{"milad", "mona"},
					BoolSlice:   []bool{false, true},
					FloatSlices: FloatSlices{
						Float32Slice: []float32{3.1415, 2.7182},
						Float64Slice: []float64{3.14159265359, 2.71828182845},
					},
					IntSlices: IntSlices{
						IntSlice:   []int{-2147483648, 2147483647},
						Int8Slice:  []int8{-128, 127},
						Int16Slice: []int16{-32768, 32767},
						Int32Slice: []int32{-2147483648, 2147483647},
						Int64Slice: []int64{-9223372036854775808, 9223372036854775807},
					},
					UintSlices: UintSlices{
						UintSlice:   []uint{0, 4294967295},
						Uint8Slice:  []uint8{0, 255},
						Uint16Slice: []uint16{0, 65535},
						Uint32Slice: []uint32{0, 4294967295},
						Uint64Slice: []uint64{0, 18446744073709551615},
					},
					DurationSlice: []time.Duration{d90m, d120m},
					URLSlice:      []url.URL{*url1, *url2},
					RegexpSlice:   []regexp.Regexp{*re1, *re2},
				},
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
				"--regexp", "[:digit:]",
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
				"--regexp-slice", "[:digit:],[:alpha:]",
			},
			new(flag.FlagSet),
			&Flags{},
			false,
			nil, "",
			&Flags{
				unexported: "",
				String:     "content",
				Bool:       true,
				Floats: Floats{
					Float32: 3.1415,
					Float64: 3.14159265359,
				},
				Ints: Ints{
					Int:   -2147483648,
					Int8:  -128,
					Int16: -32768,
					Int32: -2147483648,
					Int64: -9223372036854775808,
				},
				Uints: Uints{
					Uint:   4294967295,
					Uint8:  255,
					Uint16: 65535,
					Uint32: 4294967295,
					Uint64: 18446744073709551615,
				},
				Duration: d90m,
				URL:      *url1,
				Regexp:   *re1,
				SliceGroup: SliceGroup{
					StringSlice: []string{"milad", "mona"},
					BoolSlice:   []bool{false, true},
					FloatSlices: FloatSlices{
						Float32Slice: []float32{3.1415, 2.7182},
						Float64Slice: []float64{3.14159265359, 2.71828182845},
					},
					IntSlices: IntSlices{
						IntSlice:   []int{-2147483648, 2147483647},
						Int8Slice:  []int8{-128, 127},
						Int16Slice: []int16{-32768, 32767},
						Int32Slice: []int32{-2147483648, 2147483647},
						Int64Slice: []int64{-9223372036854775808, 9223372036854775807},
					},
					UintSlices: UintSlices{
						UintSlice:   []uint{0, 4294967295},
						Uint8Slice:  []uint8{0, 255},
						Uint16Slice: []uint16{0, 65535},
						Uint32Slice: []uint32{0, 4294967295},
						Uint64Slice: []uint64{0, 18446744073709551615},
					},
					DurationSlice: []time.Duration{d90m, d120m},
					URLSlice:      []url.URL{*url1, *url2},
					RegexpSlice:   []regexp.Regexp{*re1, *re2},
				},
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
				"-regexp=[:digit:]",
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
				"-regexp-slice=[:digit:],[:alpha:]",
			},
			new(flag.FlagSet),
			&Flags{},
			false,
			nil, "",
			&Flags{
				unexported: "",
				String:     "content",
				Bool:       true,
				Floats: Floats{
					Float32: 3.1415,
					Float64: 3.14159265359,
				},
				Ints: Ints{
					Int:   -2147483648,
					Int8:  -128,
					Int16: -32768,
					Int32: -2147483648,
					Int64: -9223372036854775808,
				},
				Uints: Uints{
					Uint:   4294967295,
					Uint8:  255,
					Uint16: 65535,
					Uint32: 4294967295,
					Uint64: 18446744073709551615,
				},
				Duration: d90m,
				URL:      *url1,
				Regexp:   *re1,
				SliceGroup: SliceGroup{
					StringSlice: []string{"milad", "mona"},
					BoolSlice:   []bool{false, true},
					FloatSlices: FloatSlices{
						Float32Slice: []float32{3.1415, 2.7182},
						Float64Slice: []float64{3.14159265359, 2.71828182845},
					},
					IntSlices: IntSlices{
						IntSlice:   []int{-2147483648, 2147483647},
						Int8Slice:  []int8{-128, 127},
						Int16Slice: []int16{-32768, 32767},
						Int32Slice: []int32{-2147483648, 2147483647},
						Int64Slice: []int64{-9223372036854775808, 9223372036854775807},
					},
					UintSlices: UintSlices{
						UintSlice:   []uint{0, 4294967295},
						Uint8Slice:  []uint8{0, 255},
						Uint16Slice: []uint16{0, 65535},
						Uint32Slice: []uint32{0, 4294967295},
						Uint64Slice: []uint64{0, 18446744073709551615},
					},
					DurationSlice: []time.Duration{d90m, d120m},
					URLSlice:      []url.URL{*url1, *url2},
					RegexpSlice:   []regexp.Regexp{*re1, *re2},
				},
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
				"--regexp=[:digit:]",
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
				"--regexp-slice=[:digit:],[:alpha:]",
			},
			new(flag.FlagSet),
			&Flags{},
			false,
			nil, "",
			&Flags{
				unexported: "",
				String:     "content",
				Bool:       true,
				Floats: Floats{
					Float32: 3.1415,
					Float64: 3.14159265359,
				},
				Ints: Ints{
					Int:   -2147483648,
					Int8:  -128,
					Int16: -32768,
					Int32: -2147483648,
					Int64: -9223372036854775808,
				},
				Uints: Uints{
					Uint:   4294967295,
					Uint8:  255,
					Uint16: 65535,
					Uint32: 4294967295,
					Uint64: 18446744073709551615,
				},
				Duration: d90m,
				URL:      *url1,
				Regexp:   *re1,
				SliceGroup: SliceGroup{
					StringSlice: []string{"milad", "mona"},
					BoolSlice:   []bool{false, true},
					FloatSlices: FloatSlices{
						Float32Slice: []float32{3.1415, 2.7182},
						Float64Slice: []float64{3.14159265359, 2.71828182845},
					},
					IntSlices: IntSlices{
						IntSlice:   []int{-2147483648, 2147483647},
						Int8Slice:  []int8{-128, 127},
						Int16Slice: []int16{-32768, 32767},
						Int32Slice: []int32{-2147483648, 2147483647},
						Int64Slice: []int64{-9223372036854775808, 9223372036854775807},
					},
					UintSlices: UintSlices{
						UintSlice:   []uint{0, 4294967295},
						Uint8Slice:  []uint8{0, 255},
						Uint16Slice: []uint16{0, 65535},
						Uint32Slice: []uint32{0, 4294967295},
						Uint64Slice: []uint64{0, 18446744073709551615},
					},
					DurationSlice: []time.Duration{d90m, d120m},
					URLSlice:      []url.URL{*url1, *url2},
					RegexpSlice:   []regexp.Regexp{*re1, *re2},
				},
			},
		},
		{
			"StopOnError",
			[]string{
				"app",
				"-int", "invalid",
			},
			new(flag.FlagSet),
			&Flags{},
			false,
			nil, `invalid value "invalid" for flag -int: strconv.ParseInt: parsing "invalid": invalid syntax`,
			&Flags{},
		},
		{
			"ContinueOnError",
			[]string{
				"app",
				"-float32", "invalid",
				"-float64", "invalid",
				"-int", "invalid",
				"-int8", "invalid",
				"-int16", "invalid",
				"-int32", "invalid",
				"-int64", "invalid",
				"-uint", "invalid",
				"-uint8", "invalid",
				"-uint16", "invalid",
				"-uint32", "invalid",
				"-uint64", "invalid",
				"-duration", "invalid",
				"-url", ":invalid",
				"-regexp", "[:invalid",
				"-bool-slice", "invalid",
				"-float32-slice", "invalid",
				"-float64-slice", "invalid",
				"-int-slice", "invalid",
				"-int8-slice", "invalid",
				"-int16-slice", "invalid",
				"-int32-slice", "invalid",
				"-int64-slice", "invalid",
				"-uint-slice", "invalid",
				"-uint8-slice", "invalid",
				"-uint16-slice", "invalid",
				"-uint32-slice", "invalid",
				"-uint64-slice", "invalid",
				"-duration-slice", "invalid",
				"-url-slice", ":invalid",
				"-regexp-slice", "[:invalid",
			},
			new(flag.FlagSet),
			&Flags{},
			true,
			nil, "",
			&Flags{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := RegisterFlags(tc.fs, tc.s, tc.continueOnError)
			assert.Equal(t, tc.expectedError, err)

			if tc.expectedError == nil {
				err := tc.fs.Parse(tc.args[1:])

				if tc.expectedParseError == "" {
					assert.NoError(t, err)
					assert.Equal(t, tc.expected, tc.s)
				} else {
					assert.Error(t, err)
					assert.EqualError(t, err, tc.expectedParseError)
				}
			}
		})
	}
}
