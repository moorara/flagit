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

	"github.com/moorara/flagit/ptr"
	"github.com/stretchr/testify/assert"
)

type (
	Values struct {
		String   string        `flag:"string"`
		Bool     bool          `flag:"bool"`
		Float32  float32       `flag:"float32"`
		Float64  float64       `flag:"float64"`
		Int      int           `flag:"int"`
		Int8     int8          `flag:"int8"`
		Int16    int16         `flag:"int16"`
		Int32    int32         `flag:"int32"`
		Int64    int64         `flag:"int64"`
		Uint     uint          `flag:"uint"`
		Uint8    uint8         `flag:"uint8"`
		Uint16   uint16        `flag:"uint16"`
		Uint32   uint32        `flag:"uint32"`
		Uint64   uint64        `flag:"uint64"`
		URL      url.URL       `flag:"url,the help text"`
		Regexp   regexp.Regexp `flag:"regexp,the help text"`
		Duration time.Duration `flag:"duration,the help text"`
	}

	Pointers struct {
		StringPointer   *string        `flag:"string-pointer"`
		BoolPointer     *bool          `flag:"bool-pointer"`
		Float32Pointer  *float32       `flag:"float32-pointer"`
		Float64Pointer  *float64       `flag:"float64-pointer"`
		IntPointer      *int           `flag:"int-pointer"`
		Int8Pointer     *int8          `flag:"int8-pointer"`
		Int16Pointer    *int16         `flag:"int16-pointer"`
		Int32Pointer    *int32         `flag:"int32-pointer"`
		Int64Pointer    *int64         `flag:"int64-pointer"`
		UintPointer     *uint          `flag:"uint-pointer"`
		Uint8Pointer    *uint8         `flag:"uint8-pointer"`
		Uint16Pointer   *uint16        `flag:"uint16-pointer"`
		Uint32Pointer   *uint32        `flag:"uint32-pointer"`
		Uint64Pointer   *uint64        `flag:"uint64-pointer"`
		URLPointer      *url.URL       `flag:"url-pointer,the help text"`
		RegexpPointer   *regexp.Regexp `flag:"regexp-pointer,the help text"`
		DurationPointer *time.Duration `flag:"duration-pointer,the help text"`
	}

	Slices struct {
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
		URLSlice      []url.URL       `flag:"url-slice,the help text"`
		RegexpSlice   []regexp.Regexp `flag:"regexp-slice,the help text"`
		DurationSlice []time.Duration `flag:"duration-slice,the help text"`
	}

	Flags struct {
		unexported  string
		WithoutFlag string
		Values
		Pointers
		Slices
	}
)

func getFields(vStruct reflect.Value, handle func(f fieldInfo)) {
	for i := 0; i < vStruct.NumField(); i++ {
		v := vStruct.Field(i)
		t := v.Type()
		f := vStruct.Type().Field(i)

		if isNestedStruct(t) {
			getFields(v, handle)
		}

		if !v.CanSet() || !isTypeSupported(t) || f.Tag.Get(flagTag) == "" {
			continue
		}

		handle(fieldInfo{
			value: v,
			name:  f.Name,
			sep:   ",",
		})
	}
}

func TestFlagValue(t *testing.T) {
	d := time.Second

	tests := []struct {
		name             string
		v                flagValue
		setVal           string
		expectedSetError string
	}{
		{
			name: "OK",
			v: flagValue{
				continueOnError: false,
				value:           reflect.ValueOf(&d).Elem(),
				sep:             ",",
			},
			setVal:           "1m",
			expectedSetError: "",
		},
		{
			name: "Error",
			v: flagValue{
				continueOnError: false,
				value:           reflect.ValueOf(&d).Elem(),
				sep:             ",",
			},
			setVal:           "invalid",
			expectedSetError: `time: invalid duration "invalid"`,
		},
		{
			name: "ContinueOnError",
			v: flagValue{
				continueOnError: true,
				value:           reflect.ValueOf(&d).Elem(),
				sep:             ",",
			},
			setVal:           "invalid",
			expectedSetError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Empty(t, tc.v.String())

			err := tc.v.Set(tc.setVal)
			if tc.expectedSetError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedSetError)
			}
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
			struct{}{},
			errors.New("non-pointer type: you should pass a pointer to a struct type"),
		},
		{
			"OK",
			new(struct{}),
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
	u, _ := url.Parse("service-1")
	r := regexp.MustCompilePOSIX("[:digit:]")

	tests := []struct {
		name     string
		field    interface{}
		expected bool
	}{
		{"String", "content", true},
		{"Bool", true, true},
		{"Float32", float32(3.1415), true},
		{"Float64", float64(3.14159265359), true},
		{"Int", int(-9223372036854775808), true},
		{"Int8", int8(-128), true},
		{"Int16", int16(-32768), true},
		{"Int32", int32(-2147483648), true},
		{"Int64", int64(-9223372036854775808), true},
		{"Uint", uint(18446744073709551615), true},
		{"Uint8", uint8(255), true},
		{"Uint16", uint16(65535), true},
		{"Uint32", uint32(4294967295), true},
		{"Uint64", uint64(18446744073709551615), true},
		{"URL", *u, true},
		{"Regexp", *r, true},
		{"Duration", time.Second, true},
		{"StringPointer", ptr.String("content"), true},
		{"BoolPointer", ptr.Bool(true), true},
		{"Float32Pointer", ptr.Float32(3.1415), true},
		{"Float64Pointer", ptr.Float64(3.14159265359), true},
		{"IntPointer", ptr.Int(-9223372036854775808), true},
		{"Int8Pointer", ptr.Int8(-128), true},
		{"Int16Pointer", ptr.Int16(-32768), true},
		{"Int32Pointer", ptr.Int32(-2147483648), true},
		{"Int64Pointer", ptr.Int64(-9223372036854775808), true},
		{"UintPointer", ptr.Uint(18446744073709551615), true},
		{"Uint8Pointer", ptr.Uint8(255), true},
		{"Uint16Pointer", ptr.Uint16(65535), true},
		{"Uint32Pointer", ptr.Uint32(4294967295), true},
		{"Uint64Pointer", ptr.Uint64(18446744073709551615), true},
		{"URLPointer", u, true},
		{"RegexpPointer", r, true},
		{"DurationPointer", ptr.Duration(time.Second), true},
		{"StringSlice", []string{"content"}, true},
		{"BoolSlice", []bool{true}, true},
		{"Float32Slice", []float32{3.1415}, true},
		{"Float64Slice", []float64{3.14159265359}, true},
		{"IntSlice", []int{-9223372036854775808}, true},
		{"Int8Slice", []int8{-128}, true},
		{"Int16Slice", []int16{-32768}, true},
		{"Int32Slice", []int32{-2147483648}, true},
		{"Int64Slice", []int64{-9223372036854775808}, true},
		{"UintSlice", []uint{18446744073709551615}, true},
		{"Uint8Slice", []uint8{255}, true},
		{"Uint16Slice", []uint16{65535}, true},
		{"Uint32Slice", []uint32{4294967295}, true},
		{"Uint64Slice", []uint64{18446744073709551615}, true},
		{"URLSlice", []url.URL{*u}, true},
		{"RegexpSlice", []regexp.Regexp{*r}, true},
		{"DurationSlice", []time.Duration{time.Second}, true},
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
		{[]string{"app=invalid"}, "invalid", ""},

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
				"URL", "Regexp", "Duration",
				"StringPointer",
				"BoolPointer",
				"Float32Pointer", "Float64Pointer",
				"IntPointer", "Int8Pointer", "Int16Pointer", "Int32Pointer", "Int64Pointer",
				"UintPointer", "Uint8Pointer", "Uint16Pointer", "Uint32Pointer", "Uint64Pointer",
				"URLPointer", "RegexpPointer", "DurationPointer",
				"StringSlice",
				"BoolSlice",
				"Float32Slice", "Float64Slice",
				"IntSlice", "Int8Slice", "Int16Slice", "Int32Slice", "Int64Slice",
				"UintSlice", "Uint8Slice", "Uint16Slice", "Uint32Slice", "Uint64Slice",
				"URLSlice", "RegexpSlice", "DurationSlice",
			},
			expectedFlagNames: []string{
				"string",
				"bool",
				"float32", "float64",
				"int", "int8", "int16", "int32", "int64",
				"uint", "uint8", "uint16", "uint32", "uint64",
				"url", "regexp", "duration",
				"string-pointer",
				"bool-pointer",
				"float32-pointer", "float64-pointer",
				"int-pointer", "int8-pointer", "int16-pointer", "int32-pointer", "int64-pointer",
				"uint-pointer", "uint8-pointer", "uint16-pointer", "uint32-pointer", "uint64-pointer",
				"url-pointer", "regexp-pointer", "duration-pointer",
				"string-slice",
				"bool-slice",
				"float32-slice", "float64-slice",
				"int-slice", "int8-slice", "int16-slice", "int32-slice", "int64-slice",
				"uint-slice", "uint8-slice", "uint16-slice", "uint32-slice", "uint64-slice",
				"url-slice", "regexp-slice", "duration-slice",
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
	url1, _ := url.Parse("service-1")
	url2, _ := url.Parse("service-2")

	re1 := regexp.MustCompilePOSIX("[:digit:]")
	re2 := regexp.MustCompilePOSIX("[:alpha:]")

	flags := &Flags{
		Values: Values{
			String:   "foo",
			Bool:     false,
			Float32:  3.1415,
			Float64:  3.14159265359,
			Int:      -9223372036854775808,
			Int8:     -128,
			Int16:    -32768,
			Int32:    -2147483648,
			Int64:    -9223372036854775808,
			Uint:     0,
			Uint8:    0,
			Uint16:   0,
			Uint32:   0,
			Uint64:   0,
			URL:      *url1,
			Regexp:   *re1,
			Duration: time.Second,
		},
		Pointers: Pointers{
			StringPointer:   ptr.String("foo"),
			BoolPointer:     ptr.Bool(false),
			Float32Pointer:  ptr.Float32(3.1415),
			Float64Pointer:  ptr.Float64(3.14159265359),
			IntPointer:      ptr.Int(-9223372036854775808),
			Int8Pointer:     ptr.Int8(-128),
			Int16Pointer:    ptr.Int16(-32768),
			Int32Pointer:    ptr.Int32(-2147483648),
			Int64Pointer:    ptr.Int64(-9223372036854775808),
			UintPointer:     ptr.Uint(0),
			Uint8Pointer:    ptr.Uint8(0),
			Uint16Pointer:   ptr.Uint16(0),
			Uint32Pointer:   ptr.Uint32(0),
			Uint64Pointer:   ptr.Uint64(0),
			URLPointer:      url1,
			RegexpPointer:   re1,
			DurationPointer: ptr.Duration(time.Second),
		},
		Slices: Slices{
			StringSlice:   []string{"foo", "bar"},
			BoolSlice:     []bool{false, true},
			Float32Slice:  []float32{3.1415, 2.7182},
			Float64Slice:  []float64{3.14159265359, 2.71828182845},
			IntSlice:      []int{-9223372036854775808, 9223372036854775807},
			Int8Slice:     []int8{-128, 127},
			Int16Slice:    []int16{-32768, 32767},
			Int32Slice:    []int32{-2147483648, 2147483647},
			Int64Slice:    []int64{-9223372036854775808, 9223372036854775807},
			UintSlice:     []uint{0, 18446744073709551615},
			Uint8Slice:    []uint8{0, 255},
			Uint16Slice:   []uint16{0, 65535},
			Uint32Slice:   []uint32{0, 4294967295},
			Uint64Slice:   []uint64{0, 18446744073709551615},
			URLSlice:      []url.URL{*url1, *url2},
			RegexpSlice:   []regexp.Regexp{*re1, *re2},
			DurationSlice: []time.Duration{time.Second, time.Minute},
		},
	}

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
			flags,
			false,
			"",
			flags,
		},
		{
			"FromFlags",
			[]string{
				"app",
				"-string=foo",
				"-bool=false",
				"-float32=3.1415",
				"-float64=3.14159265359",
				"-int=-9223372036854775808",
				"-int8=-128",
				"-int16=-32768",
				"-int32=-2147483648",
				"-int64=-9223372036854775808",
				"-uint=0",
				"-uint8=0",
				"-uint16=0",
				"-uint32=0",
				"-uint64=0",
				"-url=service-1",
				"-regexp=[:digit:]",
				"-duration=1s",
				"-string-pointer=foo",
				"-bool-pointer=false",
				"-float32-pointer=3.1415",
				"-float64-pointer=3.14159265359",
				"-int-pointer=-9223372036854775808",
				"-int8-pointer=-128",
				"-int16-pointer=-32768",
				"-int32-pointer=-2147483648",
				"-int64-pointer=-9223372036854775808",
				"-uint-pointer=0",
				"-uint8-pointer=0",
				"-uint16-pointer=0",
				"-uint32-pointer=0",
				"-uint64-pointer=0",
				"-url-pointer=service-1",
				"-regexp-pointer=[:digit:]",
				"-duration-pointer=1s",
				"-string-slice=foo,bar",
				"-bool-slice=false,true",
				"-float32-slice=3.1415,2.7182",
				"-float64-slice=3.14159265359,2.71828182845",
				"-int-slice=-9223372036854775808,9223372036854775807",
				"-int8-slice=-128,127",
				"-int16-slice=-32768,32767",
				"-int32-slice=-2147483648,2147483647",
				"-int64-slice=-9223372036854775808,9223372036854775807",
				"-uint-slice=0,18446744073709551615",
				"-uint8-slice=0,255",
				"-uint16-slice=0,65535",
				"-uint32-slice=0,4294967295",
				"-uint64-slice=0,18446744073709551615",
				"-url-slice=service-1,service-2",
				"-regexp-slice=[:digit:],[:alpha:]",
				"-duration-slice=1s,1m",
			},
			&Flags{},
			false,
			"",
			flags,
		},
		{
			"StopOnError",
			[]string{
				"app",
				"-int=invalid",
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
				"-bool=invalid",
				"-float32=invalid",
				"-float64=invalid",
				"-int=invalid",
				"-int8=invalid",
				"-int16=invalid",
				"-int32=invalid",
				"-int64=invalid",
				"-uint=invalid",
				"-uint8=invalid",
				"-uint16=invalid",
				"-uint32=invalid",
				"-uint64=invalid",
				"-url=:",
				"-regexp=[:invalid:",
				"-duration=invalid",
				"-bool-pointer=invalid",
				"-float32-pointer=invalid",
				"-float64-pointer=invalid",
				"-int-pointer=invalid",
				"-int8-pointer=invalid",
				"-int16-pointer=invalid",
				"-int32-pointer=invalid",
				"-int64-pointer=invalid",
				"-uint-pointer=invalid",
				"-uint8-pointer=invalid",
				"-uint16-pointer=invalid",
				"-uint32-pointer=invalid",
				"-uint64-pointer=invalid",
				"-url-pointer=:",
				"-regexp-pointer=[:invalid:",
				"-duration-pointer=invalid",
				"-bool-slice=invalid",
				"-float32-slice=invalid",
				"-float64-slice=invalid",
				"-int-slice=invalid",
				"-int8-slice=invalid",
				"-int16-slice=invalid",
				"-int32-slice=invalid",
				"-int64-slice=invalid",
				"-uint-slice=invalid",
				"-uint8-slice=invalid",
				"-uint16-slice=invalid",
				"-uint32-slice=invalid",
				"-uint64-slice=invalid",
				"-url-slice=:",
				"-regexp-slice=[:invalid:",
				"-duration-slice=invalid",
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
	fs := flag.NewFlagSet("app", flag.ContinueOnError)
	fs.String("string", "", "")

	url1, _ := url.Parse("service-1")
	url2, _ := url.Parse("service-2")

	re1 := regexp.MustCompilePOSIX("[:digit:]")
	re2 := regexp.MustCompilePOSIX("[:alpha:]")

	flags := &Flags{
		Values: Values{
			String:   "foo",
			Bool:     false,
			Float32:  3.1415,
			Float64:  3.14159265359,
			Int:      -9223372036854775808,
			Int8:     -128,
			Int16:    -32768,
			Int32:    -2147483648,
			Int64:    -9223372036854775808,
			Uint:     0,
			Uint8:    0,
			Uint16:   0,
			Uint32:   0,
			Uint64:   0,
			URL:      *url1,
			Regexp:   *re1,
			Duration: time.Second,
		},
		Pointers: Pointers{
			StringPointer:   ptr.String("foo"),
			BoolPointer:     ptr.Bool(false),
			Float32Pointer:  ptr.Float32(3.1415),
			Float64Pointer:  ptr.Float64(3.14159265359),
			IntPointer:      ptr.Int(-9223372036854775808),
			Int8Pointer:     ptr.Int8(-128),
			Int16Pointer:    ptr.Int16(-32768),
			Int32Pointer:    ptr.Int32(-2147483648),
			Int64Pointer:    ptr.Int64(-9223372036854775808),
			UintPointer:     ptr.Uint(0),
			Uint8Pointer:    ptr.Uint8(0),
			Uint16Pointer:   ptr.Uint16(0),
			Uint32Pointer:   ptr.Uint32(0),
			Uint64Pointer:   ptr.Uint64(0),
			URLPointer:      url1,
			RegexpPointer:   re1,
			DurationPointer: ptr.Duration(time.Second),
		},
		Slices: Slices{
			StringSlice:   []string{"foo", "bar"},
			BoolSlice:     []bool{false, true},
			Float32Slice:  []float32{3.1415, 2.7182},
			Float64Slice:  []float64{3.14159265359, 2.71828182845},
			IntSlice:      []int{-9223372036854775808, 9223372036854775807},
			Int8Slice:     []int8{-128, 127},
			Int16Slice:    []int16{-32768, 32767},
			Int32Slice:    []int32{-2147483648, 2147483647},
			Int64Slice:    []int64{-9223372036854775808, 9223372036854775807},
			UintSlice:     []uint{0, 18446744073709551615},
			Uint8Slice:    []uint8{0, 255},
			Uint16Slice:   []uint16{0, 65535},
			Uint32Slice:   []uint32{0, 4294967295},
			Uint64Slice:   []uint64{0, 18446744073709551615},
			URLSlice:      []url.URL{*url1, *url2},
			RegexpSlice:   []regexp.Regexp{*re1, *re2},
			DurationSlice: []time.Duration{time.Second, time.Minute},
		},
	}

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
			flags,
			false,
			nil, "",
			flags,
		},
		{
			"FromFlags",
			[]string{
				"app",
				"-string=foo",
				"-bool=false",
				"-float32=3.1415",
				"-float64=3.14159265359",
				"-int=-9223372036854775808",
				"-int8=-128",
				"-int16=-32768",
				"-int32=-2147483648",
				"-int64=-9223372036854775808",
				"-uint=0",
				"-uint8=0",
				"-uint16=0",
				"-uint32=0",
				"-uint64=0",
				"-url=service-1",
				"-regexp=[:digit:]",
				"-duration=1s",
				"-string-pointer=foo",
				"-bool-pointer=false",
				"-float32-pointer=3.1415",
				"-float64-pointer=3.14159265359",
				"-int-pointer=-9223372036854775808",
				"-int8-pointer=-128",
				"-int16-pointer=-32768",
				"-int32-pointer=-2147483648",
				"-int64-pointer=-9223372036854775808",
				"-uint-pointer=0",
				"-uint8-pointer=0",
				"-uint16-pointer=0",
				"-uint32-pointer=0",
				"-uint64-pointer=0",
				"-url-pointer=service-1",
				"-regexp-pointer=[:digit:]",
				"-duration-pointer=1s",
				"-string-slice=foo,bar",
				"-bool-slice=false,true",
				"-float32-slice=3.1415,2.7182",
				"-float64-slice=3.14159265359,2.71828182845",
				"-int-slice=-9223372036854775808,9223372036854775807",
				"-int8-slice=-128,127",
				"-int16-slice=-32768,32767",
				"-int32-slice=-2147483648,2147483647",
				"-int64-slice=-9223372036854775808,9223372036854775807",
				"-uint-slice=0,18446744073709551615",
				"-uint8-slice=0,255",
				"-uint16-slice=0,65535",
				"-uint32-slice=0,4294967295",
				"-uint64-slice=0,18446744073709551615",
				"-url-slice=service-1,service-2",
				"-regexp-slice=[:digit:],[:alpha:]",
				"-duration-slice=1s,1m",
			},
			new(flag.FlagSet),
			&Flags{},
			false,
			nil, "",
			flags,
		},
		{
			"StopOnError",
			[]string{
				"app",
				"-int=invalid",
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
				"-bool=invalid",
				"-float32=invalid",
				"-float64=invalid",
				"-int=invalid",
				"-int8=invalid",
				"-int16=invalid",
				"-int32=invalid",
				"-int64=invalid",
				"-uint=invalid",
				"-uint8=invalid",
				"-uint16=invalid",
				"-uint32=invalid",
				"-uint64=invalid",
				"-url=:",
				"-regexp=[:invalid:",
				"-duration=invalid",
				"-bool-pointer=invalid",
				"-float32-pointer=invalid",
				"-float64-pointer=invalid",
				"-int-pointer=invalid",
				"-int8-pointer=invalid",
				"-int16-pointer=invalid",
				"-int32-pointer=invalid",
				"-int64-pointer=invalid",
				"-uint-pointer=invalid",
				"-uint8-pointer=invalid",
				"-uint16-pointer=invalid",
				"-uint32-pointer=invalid",
				"-uint64-pointer=invalid",
				"-url-pointer=:",
				"-regexp-pointer=[:invalid:",
				"-duration-pointer=invalid",
				"-bool-slice=invalid",
				"-float32-slice=invalid",
				"-float64-slice=invalid",
				"-int-slice=invalid",
				"-int8-slice=invalid",
				"-int16-slice=invalid",
				"-int32-slice=invalid",
				"-int64-slice=invalid",
				"-uint-slice=invalid",
				"-uint8-slice=invalid",
				"-uint16-slice=invalid",
				"-uint32-slice=invalid",
				"-uint64-slice=invalid",
				"-url-slice=:",
				"-regexp-slice=[:invalid:",
				"-duration-slice=invalid",
			},
			new(flag.FlagSet),
			&Flags{},
			true,
			nil, `invalid boolean value "invalid" for -bool: parse error`,
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
