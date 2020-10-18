package set

import (
	"net/url"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/moorara/flagit/ptr"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	tests := []struct {
		name            string
		s               string
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  string
	}{
		{
			"NewValue",
			"old", "new",
			true, "",
			"new",
		},
		{
			"NoNewValue",
			"same", "same",
			false, "",
			"same",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.s).Elem()
			updated, err := String(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.s)
		})
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		name            string
		b               bool
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  bool
	}{
		{
			"NewValue",
			false, "true",
			true, "",
			true,
		},
		{
			"NoNewValue",
			true, "true",
			false, "",
			true,
		},
		{
			"InvalidValue",
			false, "invalid",
			false, `strconv.ParseBool: parsing "invalid": invalid syntax`,
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.b).Elem()
			updated, err := Bool(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.b)
		})
	}
}

func TestFloat32(t *testing.T) {
	tests := []struct {
		name            string
		f               float32
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  float32
	}{
		{
			"NewValue",
			3.1415, "2.7182",
			true, "",
			2.7182,
		},
		{
			"NoNewValue",
			2.7182, "2.7182",
			false, "",
			2.7182,
		},
		{
			"InvalidValue",
			3.1415, "invalid",
			false, `strconv.ParseFloat: parsing "invalid": invalid syntax`,
			3.1415,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.f).Elem()
			updated, err := Float32(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.f)
		})
	}
}

func TestFloat64(t *testing.T) {
	tests := []struct {
		name            string
		f               float64
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  float64
	}{
		{
			"NewValue",
			3.14159265359, "2.7182818284",
			true, "",
			2.7182818284,
		},
		{
			"NoNewValue",
			2.7182818284, "2.7182818284",
			false, "",
			2.7182818284,
		},
		{
			"InvalidValue",
			3.14159265359, "invalid",
			false, `strconv.ParseFloat: parsing "invalid": invalid syntax`,
			3.14159265359,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.f).Elem()
			updated, err := Float64(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.f)
		})
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		name            string
		i               int
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  int
	}{
		{
			"NewValue",
			-9223372036854775808, "9223372036854775807",
			true, "",
			9223372036854775807,
		},
		{
			"NoNewValue",
			9223372036854775807, "9223372036854775807",
			false, "",
			9223372036854775807,
		},
		{
			"InvalidValue",
			-9223372036854775808, "invalid",
			false, `strconv.ParseInt: parsing "invalid": invalid syntax`,
			-9223372036854775808,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.i).Elem()
			updated, err := Int(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.i)
		})
	}
}

func TestInt8(t *testing.T) {
	tests := []struct {
		name            string
		i               int8
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  int8
	}{
		{
			"NewValue",
			-128, "127",
			true, "",
			127,
		},
		{
			"NoNewValue",
			127, "127",
			false, "",
			127,
		},
		{
			"InvalidValue",
			-128, "invalid",
			false, `strconv.ParseInt: parsing "invalid": invalid syntax`,
			-128,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.i).Elem()
			updated, err := Int8(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.i)
		})
	}
}

func TestInt16(t *testing.T) {
	tests := []struct {
		name            string
		i               int16
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  int16
	}{
		{
			"NewValue",
			-32768, "32767",
			true, "",
			32767,
		},
		{
			"NoNewValue",
			32767, "32767",
			false, "",
			32767,
		},
		{
			"InvalidValue",
			-32768, "invalid",
			false, `strconv.ParseInt: parsing "invalid": invalid syntax`,
			-32768,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.i).Elem()
			updated, err := Int16(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.i)
		})
	}
}

func TestInt32(t *testing.T) {
	tests := []struct {
		name            string
		i               int32
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  int32
	}{
		{
			"NewValue",
			-2147483648, "2147483647",
			true, "",
			2147483647,
		},
		{
			"NoNewValue",
			2147483647, "2147483647",
			false, "",
			2147483647,
		},
		{
			"InvalidValue",
			-2147483648, "invalid",
			false, `strconv.ParseInt: parsing "invalid": invalid syntax`,
			-2147483648,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.i).Elem()
			updated, err := Int32(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.i)
		})
	}
}

func TestInt64(t *testing.T) {
	tests := []struct {
		name            string
		i               int64
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  int64
	}{
		{
			"NewValue",
			-9223372036854775808, "9223372036854775807",
			true, "",
			9223372036854775807,
		},
		{
			"NoNewValue",
			9223372036854775807, "9223372036854775807",
			false, "",
			9223372036854775807,
		},
		{
			"InvalidValue",
			-9223372036854775808, "invalid",
			false, `strconv.ParseInt: parsing "invalid": invalid syntax`,
			-9223372036854775808,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.i).Elem()
			updated, err := Int64(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.i)
		})
	}
}

func TestInt64Duration(t *testing.T) {
	tests := []struct {
		name            string
		i               time.Duration
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  time.Duration
	}{
		{
			"NewValue",
			time.Second, "1m",
			true, "",
			time.Minute,
		},
		{
			"NoNewValue",
			time.Minute, "1m",
			false, "",
			time.Minute,
		},
		{
			"InvalidValue",
			time.Second, "invalid",
			false, `time: invalid duration "invalid"`,
			time.Second,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.i).Elem()
			updated, err := Int64(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.i)
		})
	}
}

func TestUint(t *testing.T) {
	tests := []struct {
		name            string
		u               uint
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  uint
	}{
		{
			"NewValue",
			0, "18446744073709551615",
			true, "",
			18446744073709551615,
		},
		{
			"NoNewValue",
			18446744073709551615, "18446744073709551615",
			false, "",
			18446744073709551615,
		},
		{
			"InvalidValue",
			0, "invalid",
			false, `strconv.ParseUint: parsing "invalid": invalid syntax`,
			0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.u).Elem()
			updated, err := Uint(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.u)
		})
	}
}

func TestUint8(t *testing.T) {
	tests := []struct {
		name            string
		u               uint8
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  uint8
	}{
		{
			"NewValue",
			0, "255",
			true, "",
			255,
		},
		{
			"NoNewValue",
			255, "255",
			false, "",
			255,
		},
		{
			"InvalidValue",
			0, "invalid",
			false, `strconv.ParseUint: parsing "invalid": invalid syntax`,
			0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.u).Elem()
			updated, err := Uint8(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.u)
		})
	}
}

func TestUint16(t *testing.T) {
	tests := []struct {
		name            string
		u               uint16
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  uint16
	}{
		{
			"NewValue",
			0, "65535",
			true, "",
			65535,
		},
		{
			"NoNewValue",
			65535, "65535",
			false, "",
			65535,
		},
		{
			"InvalidValue",
			0, "invalid",
			false, `strconv.ParseUint: parsing "invalid": invalid syntax`,
			0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.u).Elem()
			updated, err := Uint16(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.u)
		})
	}
}

func TestUint32(t *testing.T) {
	tests := []struct {
		name            string
		u               uint32
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  uint32
	}{
		{
			"NewValue",
			0, "4294967295",
			true, "",
			4294967295,
		},
		{
			"NoNewValue",
			4294967295, "4294967295",
			false, "",
			4294967295,
		},
		{
			"InvalidValue",
			0, "invalid",
			false, `strconv.ParseUint: parsing "invalid": invalid syntax`,
			0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.u).Elem()
			updated, err := Uint32(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.u)
		})
	}
}

func TestUint64(t *testing.T) {
	tests := []struct {
		name            string
		u               uint64
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  uint64
	}{
		{
			"NewValue",
			0, "18446744073709551615",
			true, "",
			18446744073709551615,
		},
		{
			"NoNewValue",
			18446744073709551615, "18446744073709551615",
			false, "",
			18446744073709551615,
		},
		{
			"InvalidValue",
			0, "invalid",
			false, `strconv.ParseUint: parsing "invalid": invalid syntax`,
			0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.u).Elem()
			updated, err := Uint64(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.u)
		})
	}
}

func TestStructURL(t *testing.T) {
	url1, _ := url.Parse("service-1")
	url2, _ := url.Parse("service-2")

	tests := []struct {
		name            string
		s               url.URL
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  url.URL
	}{
		{
			"URLNewValue",
			*url1, "service-2",
			true, "",
			*url2,
		},
		{
			"URLNoNewValue",
			*url2, "service-2",
			false, "",
			*url2,
		},
		{
			"URLInvalidValue",
			*url1, ":",
			false, `parse ":": missing protocol scheme`,
			*url1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.s).Elem()
			updated, err := Struct(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.s)
		})
	}
}

func TestStructRegexp(t *testing.T) {
	re1 := regexp.MustCompilePOSIX("[:digit:]")
	re2 := regexp.MustCompilePOSIX("[:alpha:]")

	tests := []struct {
		name            string
		s               regexp.Regexp
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  regexp.Regexp
	}{
		{
			"RegexpNewValue",
			*re1, "[:alpha:]",
			true, "",
			*re2,
		},
		{
			"RegexpNoNewValue",
			*re2, "[:alpha:]",
			false, "",
			*re2,
		},
		{
			"RegexpInvalidValue",
			*re1, "[:invalid:",
			false, "error parsing regexp: missing closing ]: `[:invalid:`",
			*re1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.s).Elem()
			updated, err := Struct(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.s)
		})
	}
}

func TestStringPtr(t *testing.T) {
	tests := []struct {
		name            string
		s               *string
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  *string
	}{
		{
			"Nil",
			nil, "new",
			true, "",
			ptr.String("new"),
		},
		{
			"NewValue",
			ptr.String("old"), "new",
			true, "",
			ptr.String("new"),
		},
		{
			"NoNewValue",
			ptr.String("same"), "same",
			false, "",
			ptr.String("same"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.s).Elem()
			updated, err := StringPtr(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.s)
		})
	}
}

func TestBoolPtr(t *testing.T) {
	tests := []struct {
		name            string
		b               *bool
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  *bool
	}{
		{
			"Nil",
			nil, "true",
			true, "",
			ptr.Bool(true),
		},
		{
			"NewValue",
			ptr.Bool(false), "true",
			true, "",
			ptr.Bool(true),
		},
		{
			"NoNewValue",
			ptr.Bool(true), "true",
			false, "",
			ptr.Bool(true),
		},
		{
			"InvalidValue",
			ptr.Bool(false), "invalid",
			false, `strconv.ParseBool: parsing "invalid": invalid syntax`,
			ptr.Bool(false),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.b).Elem()
			updated, err := BoolPtr(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.b)
		})
	}
}

func TestFloat32Ptr(t *testing.T) {
	tests := []struct {
		name            string
		f               *float32
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  *float32
	}{
		{
			"Nil",
			nil, "2.7182",
			true, "",
			ptr.Float32(2.7182),
		},
		{
			"NewValue",
			ptr.Float32(3.1415), "2.7182",
			true, "",
			ptr.Float32(2.7182),
		},
		{
			"NoNewValue",
			ptr.Float32(2.7182), "2.7182",
			false, "",
			ptr.Float32(2.7182),
		},
		{
			"InvalidValue",
			ptr.Float32(3.1415), "invalid",
			false, `strconv.ParseFloat: parsing "invalid": invalid syntax`,
			ptr.Float32(3.1415),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.f).Elem()
			updated, err := Float32Ptr(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.f)
		})
	}
}

func TestFloat64Ptr(t *testing.T) {
	tests := []struct {
		name            string
		f               *float64
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  *float64
	}{
		{
			"Nil",
			nil, "2.7182818284",
			true, "",
			ptr.Float64(2.7182818284),
		},
		{
			"NewValue",
			ptr.Float64(3.14159265359), "2.7182818284",
			true, "",
			ptr.Float64(2.7182818284),
		},
		{
			"NoNewValue",
			ptr.Float64(2.7182818284), "2.7182818284",
			false, "",
			ptr.Float64(2.7182818284),
		},
		{
			"InvalidValue",
			ptr.Float64(3.14159265359), "invalid",
			false, `strconv.ParseFloat: parsing "invalid": invalid syntax`,
			ptr.Float64(3.14159265359),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.f).Elem()
			updated, err := Float64Ptr(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.f)
		})
	}
}

func TestIntPtr(t *testing.T) {
	tests := []struct {
		name            string
		i               *int
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  *int
	}{
		{
			"Nil",
			nil, "9223372036854775807",
			true, "",
			ptr.Int(9223372036854775807),
		},
		{
			"NewValue",
			ptr.Int(-9223372036854775808), "9223372036854775807",
			true, "",
			ptr.Int(9223372036854775807),
		},
		{
			"NoNewValue",
			ptr.Int(9223372036854775807), "9223372036854775807",
			false, "",
			ptr.Int(9223372036854775807),
		},
		{
			"InvalidValue",
			ptr.Int(-9223372036854775808), "invalid",
			false, `strconv.ParseInt: parsing "invalid": invalid syntax`,
			ptr.Int(-9223372036854775808),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.i).Elem()
			updated, err := IntPtr(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.i)
		})
	}
}

func TestInt8Ptr(t *testing.T) {
	tests := []struct {
		name            string
		i               *int8
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  *int8
	}{
		{
			"Nil",
			nil, "127",
			true, "",
			ptr.Int8(127),
		},
		{
			"NewValue",
			ptr.Int8(-128), "127",
			true, "",
			ptr.Int8(127),
		},
		{
			"NoNewValue",
			ptr.Int8(127), "127",
			false, "",
			ptr.Int8(127),
		},
		{
			"InvalidValue",
			ptr.Int8(-128), "invalid",
			false, `strconv.ParseInt: parsing "invalid": invalid syntax`,
			ptr.Int8(-128),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.i).Elem()
			updated, err := Int8Ptr(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.i)
		})
	}
}

func TestInt16Ptr(t *testing.T) {
	tests := []struct {
		name            string
		i               *int16
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  *int16
	}{
		{
			"Nil",
			nil, "32767",
			true, "",
			ptr.Int16(32767),
		},
		{
			"NewValue",
			ptr.Int16(-32768), "32767",
			true, "",
			ptr.Int16(32767),
		},
		{
			"NoNewValue",
			ptr.Int16(32767), "32767",
			false, "",
			ptr.Int16(32767),
		},
		{
			"InvalidValue",
			ptr.Int16(-32768), "invalid",
			false, `strconv.ParseInt: parsing "invalid": invalid syntax`,
			ptr.Int16(-32768),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.i).Elem()
			updated, err := Int16Ptr(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.i)
		})
	}
}

func TestInt32Ptr(t *testing.T) {
	tests := []struct {
		name            string
		i               *int32
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  *int32
	}{
		{
			"Nil",
			nil, "2147483647",
			true, "",
			ptr.Int32(2147483647),
		},
		{
			"NewValue",
			ptr.Int32(-2147483648), "2147483647",
			true, "",
			ptr.Int32(2147483647),
		},
		{
			"NoNewValue",
			ptr.Int32(2147483647), "2147483647",
			false, "",
			ptr.Int32(2147483647),
		},
		{
			"InvalidValue",
			ptr.Int32(-2147483648), "invalid",
			false, `strconv.ParseInt: parsing "invalid": invalid syntax`,
			ptr.Int32(-2147483648),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.i).Elem()
			updated, err := Int32Ptr(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.i)
		})
	}
}

func TestInt64Ptr(t *testing.T) {
	tests := []struct {
		name            string
		i               *int64
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  *int64
	}{
		{
			"Nil",
			nil, "9223372036854775807",
			true, "",
			ptr.Int64(9223372036854775807),
		},
		{
			"NewValue",
			ptr.Int64(-9223372036854775808), "9223372036854775807",
			true, "",
			ptr.Int64(9223372036854775807),
		},
		{
			"NoNewValue",
			ptr.Int64(9223372036854775807), "9223372036854775807",
			false, "",
			ptr.Int64(9223372036854775807),
		},
		{
			"InvalidValue",
			ptr.Int64(-9223372036854775808), "invalid",
			false, `strconv.ParseInt: parsing "invalid": invalid syntax`,
			ptr.Int64(-9223372036854775808),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.i).Elem()
			updated, err := Int64Ptr(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.i)
		})
	}
}

func TestInt64DurationPtr(t *testing.T) {
	tests := []struct {
		name            string
		i               *time.Duration
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  *time.Duration
	}{
		{
			"Nil",
			nil, "1m",
			true, "",
			ptr.Duration(time.Minute),
		},
		{
			"NewValue",
			ptr.Duration(time.Second), "1m",
			true, "",
			ptr.Duration(time.Minute),
		},
		{
			"NoNewValue",
			ptr.Duration(time.Minute), "1m",
			false, "",
			ptr.Duration(time.Minute),
		},
		{
			"InvalidValue",
			ptr.Duration(time.Second), "invalid",
			false, `time: invalid duration "invalid"`,
			ptr.Duration(time.Second),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.i).Elem()
			updated, err := Int64Ptr(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.i)
		})
	}
}

func TestUintPtr(t *testing.T) {
	tests := []struct {
		name            string
		u               *uint
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  *uint
	}{
		{
			"Nil",
			nil, "18446744073709551615",
			true, "",
			ptr.Uint(18446744073709551615),
		},
		{
			"NewValue",
			ptr.Uint(0), "18446744073709551615",
			true, "",
			ptr.Uint(18446744073709551615),
		},
		{
			"NoNewValue",
			ptr.Uint(18446744073709551615), "18446744073709551615",
			false, "",
			ptr.Uint(18446744073709551615),
		},
		{
			"InvalidValue",
			ptr.Uint(0), "invalid",
			false, `strconv.ParseUint: parsing "invalid": invalid syntax`,
			ptr.Uint(0),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.u).Elem()
			updated, err := UintPtr(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.u)
		})
	}
}

func TestUint8Ptr(t *testing.T) {
	tests := []struct {
		name            string
		u               *uint8
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  *uint8
	}{
		{
			"Nil",
			nil, "255",
			true, "",
			ptr.Uint8(255),
		},
		{
			"NewValue",
			ptr.Uint8(0), "255",
			true, "",
			ptr.Uint8(255),
		},
		{
			"NoNewValue",
			ptr.Uint8(255), "255",
			false, "",
			ptr.Uint8(255),
		},
		{
			"InvalidValue",
			ptr.Uint8(0), "invalid",
			false, `strconv.ParseUint: parsing "invalid": invalid syntax`,
			ptr.Uint8(0),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.u).Elem()
			updated, err := Uint8Ptr(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.u)
		})
	}
}

func TestUint16Ptr(t *testing.T) {
	tests := []struct {
		name            string
		u               *uint16
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  *uint16
	}{
		{
			"Nil",
			nil, "65535",
			true, "",
			ptr.Uint16(65535),
		},
		{
			"NewValue",
			ptr.Uint16(0), "65535",
			true, "",
			ptr.Uint16(65535),
		},
		{
			"NoNewValue",
			ptr.Uint16(65535), "65535",
			false, "",
			ptr.Uint16(65535),
		},
		{
			"InvalidValue",
			ptr.Uint16(0), "invalid",
			false, `strconv.ParseUint: parsing "invalid": invalid syntax`,
			ptr.Uint16(0),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.u).Elem()
			updated, err := Uint16Ptr(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.u)
		})
	}
}

func TestUint32Ptr(t *testing.T) {
	tests := []struct {
		name            string
		u               *uint32
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  *uint32
	}{
		{
			"Nil",
			nil, "4294967295",
			true, "",
			ptr.Uint32(4294967295),
		},
		{
			"NewValue",
			ptr.Uint32(0), "4294967295",
			true, "",
			ptr.Uint32(4294967295),
		},
		{
			"NoNewValue",
			ptr.Uint32(4294967295), "4294967295",
			false, "",
			ptr.Uint32(4294967295),
		},
		{
			"InvalidValue",
			ptr.Uint32(0), "invalid",
			false, `strconv.ParseUint: parsing "invalid": invalid syntax`,
			ptr.Uint32(0),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.u).Elem()
			updated, err := Uint32Ptr(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.u)
		})
	}
}

func TestUint64Ptr(t *testing.T) {
	tests := []struct {
		name            string
		u               *uint64
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  *uint64
	}{
		{
			"Nil",
			nil, "18446744073709551615",
			true, "",
			ptr.Uint64(18446744073709551615),
		},
		{
			"NewValue",
			ptr.Uint64(0), "18446744073709551615",
			true, "",
			ptr.Uint64(18446744073709551615),
		},
		{
			"NoNewValue",
			ptr.Uint64(18446744073709551615), "18446744073709551615",
			false, "",
			ptr.Uint64(18446744073709551615),
		},
		{
			"InvalidValue",
			ptr.Uint64(0), "invalid",
			false, `strconv.ParseUint: parsing "invalid": invalid syntax`,
			ptr.Uint64(0),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.u).Elem()
			updated, err := Uint64Ptr(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.u)
		})
	}
}

func TestStructPtrURL(t *testing.T) {
	url1, _ := url.Parse("service-1")
	url2, _ := url.Parse("service-2")

	tests := []struct {
		name            string
		s               *url.URL
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  *url.URL
	}{
		{
			"URLNil",
			nil, "service-2",
			true, "",
			url2,
		},
		{
			"URLNewValue",
			url1, "service-2",
			true, "",
			url2,
		},
		{
			"URLNoNewValue",
			url2, "service-2",
			false, "",
			url2,
		},
		{
			"URLInvalidValue",
			url1, ":",
			false, `parse ":": missing protocol scheme`,
			url1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.s).Elem()
			updated, err := StructPtr(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.s)
		})
	}
}

func TestStructPtrRegexp(t *testing.T) {
	re1 := regexp.MustCompilePOSIX("[:digit:]")
	re2 := regexp.MustCompilePOSIX("[:alpha:]")

	tests := []struct {
		name            string
		s               *regexp.Regexp
		val             string
		expectedUpdated bool
		expectedError   string
		expectedResult  *regexp.Regexp
	}{
		{
			"RegexpNil",
			nil, "[:alpha:]",
			true, "",
			re2,
		},
		{
			"RegexpNewValue",
			re1, "[:alpha:]",
			true, "",
			re2,
		},
		{
			"RegexpNoNewValue",
			re2, "[:alpha:]",
			false, "",
			re2,
		},
		{
			"RegexpInvalidValue",
			re1, "[:invalid:",
			false, "error parsing regexp: missing closing ]: `[:invalid:`",
			re1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.s).Elem()
			updated, err := StructPtr(v, tc.val)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.s)
		})
	}
}

func TestStringSlice(t *testing.T) {
	tests := []struct {
		name            string
		s               []string
		vals            []string
		expectedUpdated bool
		expectedError   string
		expectedResult  []string
	}{
		{
			"Nil",
			nil, []string{"new"},
			true, "",
			[]string{"new"},
		},
		{
			"NewValue",
			[]string{"old"}, []string{"new"},
			true, "",
			[]string{"new"},
		},
		{
			"NoNewValue",
			[]string{"same"}, []string{"same"},
			false, "",
			[]string{"same"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.s).Elem()
			updated, err := StringSlice(v, tc.vals)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.s)
		})
	}
}

func TestBoolSlice(t *testing.T) {
	tests := []struct {
		name            string
		b               []bool
		vals            []string
		expectedUpdated bool
		expectedError   string
		expectedResult  []bool
	}{
		{
			"Nil",
			nil, []string{"true"},
			true, "",
			[]bool{true},
		},
		{
			"NewValue",
			[]bool{false}, []string{"true"},
			true, "",
			[]bool{true},
		},
		{
			"NoNewValue",
			[]bool{true}, []string{"true"},
			false, "",
			[]bool{true},
		},
		{
			"InvalidValue",
			[]bool{false}, []string{"invalid"},
			false, `strconv.ParseBool: parsing "invalid": invalid syntax`,
			[]bool{false},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.b).Elem()
			updated, err := BoolSlice(v, tc.vals)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.b)
		})
	}
}

func TestFloat32Slice(t *testing.T) {
	tests := []struct {
		name            string
		f               []float32
		vals            []string
		expectedUpdated bool
		expectedError   string
		expectedResult  []float32
	}{
		{
			"Nil",
			nil, []string{"2.7182"},
			true, "",
			[]float32{2.7182},
		},
		{
			"NewValue",
			[]float32{3.1415}, []string{"2.7182"},
			true, "",
			[]float32{2.7182},
		},
		{
			"NoNewValue",
			[]float32{2.7182}, []string{"2.7182"},
			false, "",
			[]float32{2.7182},
		},
		{
			"InvalidValue",
			[]float32{3.1415}, []string{"invalid"},
			false, `strconv.ParseFloat: parsing "invalid": invalid syntax`,
			[]float32{3.1415},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.f).Elem()
			updated, err := Float32Slice(v, tc.vals)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.f)
		})
	}
}

func TestFloat64Slice(t *testing.T) {
	tests := []struct {
		name            string
		f               []float64
		vals            []string
		expectedUpdated bool
		expectedError   string
		expectedResult  []float64
	}{
		{
			"Nil",
			nil, []string{"2.7182818284"},
			true, "",
			[]float64{2.7182818284},
		},
		{
			"NewValue",
			[]float64{3.14159265359}, []string{"2.7182818284"},
			true, "",
			[]float64{2.7182818284},
		},
		{
			"NoNewValue",
			[]float64{2.7182818284}, []string{"2.7182818284"},
			false, "",
			[]float64{2.7182818284},
		},
		{
			"InvalidValue",
			[]float64{3.14159265359}, []string{"invalid"},
			false, `strconv.ParseFloat: parsing "invalid": invalid syntax`,
			[]float64{3.14159265359},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.f).Elem()
			updated, err := Float64Slice(v, tc.vals)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.f)
		})
	}
}

func TestIntSlice(t *testing.T) {
	tests := []struct {
		name            string
		i               []int
		vals            []string
		expectedUpdated bool
		expectedError   string
		expectedResult  []int
	}{
		{
			"Nil",
			nil, []string{"9223372036854775807"},
			true, "",
			[]int{9223372036854775807},
		},
		{
			"NewValue",
			[]int{-9223372036854775808}, []string{"9223372036854775807"},
			true, "",
			[]int{9223372036854775807},
		},
		{
			"NoNewValue",
			[]int{9223372036854775807}, []string{"9223372036854775807"},
			false, "",
			[]int{9223372036854775807},
		},
		{
			"InvalidValue",
			[]int{-9223372036854775808}, []string{"invalid"},
			false, `strconv.ParseInt: parsing "invalid": invalid syntax`,
			[]int{-9223372036854775808},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.i).Elem()
			updated, err := IntSlice(v, tc.vals)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.i)
		})
	}
}

func TestInt8Slice(t *testing.T) {
	tests := []struct {
		name            string
		i               []int8
		vals            []string
		expectedUpdated bool
		expectedError   string
		expectedResult  []int8
	}{
		{
			"Nil",
			nil, []string{"127"},
			true, "",
			[]int8{127},
		},
		{
			"NewValue",
			[]int8{-128}, []string{"127"},
			true, "",
			[]int8{127},
		},
		{
			"NoNewValue",
			[]int8{127}, []string{"127"},
			false, "",
			[]int8{127},
		},
		{
			"InvalidValue",
			[]int8{-128}, []string{"invalid"},
			false, `strconv.ParseInt: parsing "invalid": invalid syntax`,
			[]int8{-128},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.i).Elem()
			updated, err := Int8Slice(v, tc.vals)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.i)
		})
	}
}

func TestInt16Slice(t *testing.T) {
	tests := []struct {
		name            string
		i               []int16
		vals            []string
		expectedUpdated bool
		expectedError   string
		expectedResult  []int16
	}{
		{
			"Nil",
			nil, []string{"32767"},
			true, "",
			[]int16{32767},
		},
		{
			"NewValue",
			[]int16{-32768}, []string{"32767"},
			true, "",
			[]int16{32767},
		},
		{
			"NoNewValue",
			[]int16{32767}, []string{"32767"},
			false, "",
			[]int16{32767},
		},
		{
			"InvalidValue",
			[]int16{-32768}, []string{"invalid"},
			false, `strconv.ParseInt: parsing "invalid": invalid syntax`,
			[]int16{-32768},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.i).Elem()
			updated, err := Int16Slice(v, tc.vals)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.i)
		})
	}
}

func TestInt32Slice(t *testing.T) {
	tests := []struct {
		name            string
		i               []int32
		vals            []string
		expectedUpdated bool
		expectedError   string
		expectedResult  []int32
	}{
		{
			"Nil",
			nil, []string{"2147483647"},
			true, "",
			[]int32{2147483647},
		},
		{
			"NewValue",
			[]int32{-2147483648}, []string{"2147483647"},
			true, "",
			[]int32{2147483647},
		},
		{
			"NoNewValue",
			[]int32{2147483647}, []string{"2147483647"},
			false, "",
			[]int32{2147483647},
		},
		{
			"InvalidValue",
			[]int32{-2147483648}, []string{"invalid"},
			false, `strconv.ParseInt: parsing "invalid": invalid syntax`,
			[]int32{-2147483648},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.i).Elem()
			updated, err := Int32Slice(v, tc.vals)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.i)
		})
	}
}

func TestInt64Slice(t *testing.T) {
	tests := []struct {
		name            string
		i               []int64
		vals            []string
		expectedUpdated bool
		expectedError   string
		expectedResult  []int64
	}{
		{
			"Nil",
			nil, []string{"9223372036854775807"},
			true, "",
			[]int64{9223372036854775807},
		},
		{
			"NewValue",
			[]int64{-9223372036854775808}, []string{"9223372036854775807"},
			true, "",
			[]int64{9223372036854775807},
		},
		{
			"NoNewValue",
			[]int64{9223372036854775807}, []string{"9223372036854775807"},
			false, "",
			[]int64{9223372036854775807},
		},
		{
			"InvalidValue",
			[]int64{-9223372036854775808}, []string{"invalid"},
			false, `strconv.ParseInt: parsing "invalid": invalid syntax`,
			[]int64{-9223372036854775808},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.i).Elem()
			updated, err := Int64Slice(v, tc.vals)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.i)
		})
	}
}

func TestInt64DurationSlice(t *testing.T) {
	tests := []struct {
		name            string
		i               []time.Duration
		vals            []string
		expectedUpdated bool
		expectedError   string
		expectedResult  []time.Duration
	}{
		{
			"Nil",
			nil, []string{"1m"},
			true, "",
			[]time.Duration{time.Minute},
		},
		{
			"NewValue",
			[]time.Duration{time.Second}, []string{"1m"},
			true, "",
			[]time.Duration{time.Minute},
		},
		{
			"NoNewValue",
			[]time.Duration{time.Minute}, []string{"1m"},
			false, "",
			[]time.Duration{time.Minute},
		},
		{
			"InvalidValue",
			[]time.Duration{time.Second}, []string{"invalid"},
			false, `time: invalid duration "invalid"`,
			[]time.Duration{time.Second},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.i).Elem()
			updated, err := Int64Slice(v, tc.vals)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.i)
		})
	}
}

func TestUintSlice(t *testing.T) {
	tests := []struct {
		name            string
		u               []uint
		vals            []string
		expectedUpdated bool
		expectedError   string
		expectedResult  []uint
	}{
		{
			"Nil",
			nil, []string{"18446744073709551615"},
			true, "",
			[]uint{18446744073709551615},
		},
		{
			"NewValue",
			[]uint{0}, []string{"18446744073709551615"},
			true, "",
			[]uint{18446744073709551615},
		},
		{
			"NoNewValue",
			[]uint{18446744073709551615}, []string{"18446744073709551615"},
			false, "",
			[]uint{18446744073709551615},
		},
		{
			"InvalidValue",
			[]uint{0}, []string{"invalid"},
			false, `strconv.ParseUint: parsing "invalid": invalid syntax`,
			[]uint{0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.u).Elem()
			updated, err := UintSlice(v, tc.vals)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.u)
		})
	}
}

func TestUint8Slice(t *testing.T) {
	tests := []struct {
		name            string
		u               []uint8
		vals            []string
		expectedUpdated bool
		expectedError   string
		expectedResult  []uint8
	}{
		{
			"Nil",
			nil, []string{"255"},
			true, "",
			[]uint8{255},
		},
		{
			"NewValue",
			[]uint8{0}, []string{"255"},
			true, "",
			[]uint8{255},
		},
		{
			"NoNewValue",
			[]uint8{255}, []string{"255"},
			false, "",
			[]uint8{255},
		},
		{
			"InvalidValue",
			[]uint8{0}, []string{"invalid"},
			false, `strconv.ParseUint: parsing "invalid": invalid syntax`,
			[]uint8{0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.u).Elem()
			updated, err := Uint8Slice(v, tc.vals)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.u)
		})
	}
}

func TestUint16Slice(t *testing.T) {
	tests := []struct {
		name            string
		u               []uint16
		vals            []string
		expectedUpdated bool
		expectedError   string
		expectedResult  []uint16
	}{
		{
			"Nil",
			nil, []string{"65535"},
			true, "",
			[]uint16{65535},
		},
		{
			"NewValue",
			[]uint16{0}, []string{"65535"},
			true, "",
			[]uint16{65535},
		},
		{
			"NoNewValue",
			[]uint16{65535}, []string{"65535"},
			false, "",
			[]uint16{65535},
		},
		{
			"InvalidValue",
			[]uint16{0}, []string{"invalid"},
			false, `strconv.ParseUint: parsing "invalid": invalid syntax`,
			[]uint16{0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.u).Elem()
			updated, err := Uint16Slice(v, tc.vals)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.u)
		})
	}
}

func TestUint32Slice(t *testing.T) {
	tests := []struct {
		name            string
		u               []uint32
		vals            []string
		expectedUpdated bool
		expectedError   string
		expectedResult  []uint32
	}{
		{
			"Nil",
			nil, []string{"4294967295"},
			true, "",
			[]uint32{4294967295},
		},
		{
			"NewValue",
			[]uint32{0}, []string{"4294967295"},
			true, "",
			[]uint32{4294967295},
		},
		{
			"NoNewValue",
			[]uint32{4294967295}, []string{"4294967295"},
			false, "",
			[]uint32{4294967295},
		},
		{
			"InvalidValue",
			[]uint32{0}, []string{"invalid"},
			false, `strconv.ParseUint: parsing "invalid": invalid syntax`,
			[]uint32{0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.u).Elem()
			updated, err := Uint32Slice(v, tc.vals)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.u)
		})
	}
}

func TestUint64Slice(t *testing.T) {
	tests := []struct {
		name            string
		u               []uint64
		vals            []string
		expectedUpdated bool
		expectedError   string
		expectedResult  []uint64
	}{
		{
			"Nil",
			nil, []string{"18446744073709551615"},
			true, "",
			[]uint64{18446744073709551615},
		},
		{
			"NewValue",
			[]uint64{0}, []string{"18446744073709551615"},
			true, "",
			[]uint64{18446744073709551615},
		},
		{
			"NoNewValue",
			[]uint64{18446744073709551615}, []string{"18446744073709551615"},
			false, "",
			[]uint64{18446744073709551615},
		},
		{
			"InvalidValue",
			[]uint64{0}, []string{"invalid"},
			false, `strconv.ParseUint: parsing "invalid": invalid syntax`,
			[]uint64{0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.u).Elem()
			updated, err := Uint64Slice(v, tc.vals)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.u)
		})
	}
}

func TestStructSliceURL(t *testing.T) {
	url1, _ := url.Parse("service-1")
	url2, _ := url.Parse("service-2")

	tests := []struct {
		name            string
		s               []url.URL
		vals            []string
		expectedUpdated bool
		expectedError   string
		expectedResult  []url.URL
	}{
		{
			"URLNil",
			nil, []string{"service-2"},
			true, "",
			[]url.URL{*url2},
		},
		{
			"URLNewValue",
			[]url.URL{*url1}, []string{"service-2"},
			true, "",
			[]url.URL{*url2},
		},
		{
			"URLNoNewValue",
			[]url.URL{*url2}, []string{"service-2"},
			false, "",
			[]url.URL{*url2},
		},
		{
			"URLInvalidValue",
			[]url.URL{*url1}, []string{":"},
			false, `parse ":": missing protocol scheme`,
			[]url.URL{*url1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.s).Elem()
			updated, err := StructSlice(v, tc.vals)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.s)
		})
	}
}

func TestStructSliceRegexp(t *testing.T) {
	re1 := regexp.MustCompilePOSIX("[:digit:]")
	re2 := regexp.MustCompilePOSIX("[:alpha:]")

	tests := []struct {
		name            string
		s               []regexp.Regexp
		vals            []string
		expectedUpdated bool
		expectedError   string
		expectedResult  []regexp.Regexp
	}{
		{
			"RegexpNil",
			nil, []string{"[:alpha:]"},
			true, "",
			[]regexp.Regexp{*re2},
		},
		{
			"RegexpNewValue",
			[]regexp.Regexp{*re1}, []string{"[:alpha:]"},
			true, "",
			[]regexp.Regexp{*re2},
		},
		{
			"RegexpNoNewValue",
			[]regexp.Regexp{*re2}, []string{"[:alpha:]"},
			false, "",
			[]regexp.Regexp{*re2},
		},
		{
			"RegexpInvalidValue",
			[]regexp.Regexp{*re1}, []string{"[:invalid:"},
			false, "error parsing regexp: missing closing ]: `[:invalid:`",
			[]regexp.Regexp{*re1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.s).Elem()
			updated, err := StructSlice(v, tc.vals)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, tc.expectedResult, tc.s)
		})
	}
}

func TestValue(t *testing.T) {
	type Fields struct {
		String        string
		Bool          bool
		Float32       float32
		Float64       float64
		Int           int
		Int8          int8
		Int16         int16
		Int32         int32
		Int64         int64
		Uint          uint
		Uint8         uint8
		Uint16        uint16
		Uint32        uint32
		Uint64        uint64
		Duration      time.Duration
		URL           url.URL
		Regexp        regexp.Regexp
		StringPtr     *string
		BoolPtr       *bool
		Float32Ptr    *float32
		Float64Ptr    *float64
		IntPtr        *int
		Int8Ptr       *int8
		Int16Ptr      *int16
		Int32Ptr      *int32
		Int64Ptr      *int64
		UintPtr       *uint
		Uint8Ptr      *uint8
		Uint16Ptr     *uint16
		Uint32Ptr     *uint32
		Uint64Ptr     *uint64
		DurationPtr   *time.Duration
		URLPtr        *url.URL
		RegexpPtr     *regexp.Regexp
		StringSlice   []string
		BoolSlice     []bool
		Float32Slice  []float32
		Float64Slice  []float64
		IntSlice      []int
		Int8Slice     []int8
		Int16Slice    []int16
		Int32Slice    []int32
		Int64Slice    []int64
		UintSlice     []uint
		Uint8Slice    []uint8
		Uint16Slice   []uint16
		Uint32Slice   []uint32
		Uint64Slice   []uint64
		DurationSlice []time.Duration
		URLSlice      []url.URL
		RegexpSlice   []regexp.Regexp
	}

	url1, _ := url.Parse("service-1")
	url2, _ := url.Parse("service-2")

	re1 := regexp.MustCompilePOSIX("[:digit:]")
	re2 := regexp.MustCompilePOSIX("[:alpha:]")

	f1 := Fields{
		String:        "old",
		Bool:          false,
		Float32:       3.1415,
		Float64:       3.14159265359,
		Int:           -9223372036854775808,
		Int8:          -128,
		Int16:         -32768,
		Int32:         -2147483648,
		Int64:         -9223372036854775808,
		Uint:          0,
		Uint8:         0,
		Uint16:        0,
		Uint32:        0,
		Uint64:        0,
		Duration:      time.Second,
		URL:           *url1,
		Regexp:        *re1,
		StringPtr:     ptr.String("old"),
		BoolPtr:       ptr.Bool(false),
		Float32Ptr:    ptr.Float32(3.1415),
		Float64Ptr:    ptr.Float64(3.14159265359),
		IntPtr:        ptr.Int(-9223372036854775808),
		Int8Ptr:       ptr.Int8(-128),
		Int16Ptr:      ptr.Int16(-32768),
		Int32Ptr:      ptr.Int32(-2147483648),
		Int64Ptr:      ptr.Int64(-9223372036854775808),
		UintPtr:       ptr.Uint(0),
		Uint8Ptr:      ptr.Uint8(0),
		Uint16Ptr:     ptr.Uint16(0),
		Uint32Ptr:     ptr.Uint32(0),
		Uint64Ptr:     ptr.Uint64(0),
		DurationPtr:   ptr.Duration(time.Second),
		URLPtr:        url1,
		RegexpPtr:     re1,
		StringSlice:   []string{"old"},
		BoolSlice:     []bool{false},
		Float32Slice:  []float32{3.1415},
		Float64Slice:  []float64{3.14159265359},
		IntSlice:      []int{-2147483648},
		Int8Slice:     []int8{-128},
		Int16Slice:    []int16{-32768},
		Int32Slice:    []int32{-2147483648},
		Int64Slice:    []int64{-9223372036854775808},
		UintSlice:     []uint{0},
		Uint8Slice:    []uint8{0},
		Uint16Slice:   []uint16{0},
		Uint32Slice:   []uint32{0},
		Uint64Slice:   []uint64{0},
		DurationSlice: []time.Duration{time.Second},
		URLSlice:      []url.URL{*url1, *url2},
		RegexpSlice:   []regexp.Regexp{*re1, *re2},
	}

	f2 := Fields{
		String:        "new",
		Bool:          true,
		Float32:       2.7182,
		Float64:       2.7182818284,
		Int:           9223372036854775807,
		Int8:          127,
		Int16:         32767,
		Int32:         2147483647,
		Int64:         9223372036854775807,
		Uint:          18446744073709551615,
		Uint8:         255,
		Uint16:        65535,
		Uint32:        4294967295,
		Uint64:        18446744073709551615,
		Duration:      time.Minute,
		URL:           *url2,
		Regexp:        *re2,
		StringPtr:     ptr.String("new"),
		BoolPtr:       ptr.Bool(true),
		Float32Ptr:    ptr.Float32(2.7182),
		Float64Ptr:    ptr.Float64(2.7182818284),
		IntPtr:        ptr.Int(9223372036854775807),
		Int8Ptr:       ptr.Int8(127),
		Int16Ptr:      ptr.Int16(32767),
		Int32Ptr:      ptr.Int32(2147483647),
		Int64Ptr:      ptr.Int64(9223372036854775807),
		UintPtr:       ptr.Uint(18446744073709551615),
		Uint8Ptr:      ptr.Uint8(255),
		Uint16Ptr:     ptr.Uint16(65535),
		Uint32Ptr:     ptr.Uint32(4294967295),
		Uint64Ptr:     ptr.Uint64(18446744073709551615),
		DurationPtr:   ptr.Duration(time.Minute),
		URLPtr:        url2,
		RegexpPtr:     re2,
		StringSlice:   []string{"new"},
		BoolSlice:     []bool{true},
		Float32Slice:  []float32{2.7182},
		Float64Slice:  []float64{2.7182818284},
		IntSlice:      []int{9223372036854775807},
		Int8Slice:     []int8{127},
		Int16Slice:    []int16{32767},
		Int32Slice:    []int32{2147483647},
		Int64Slice:    []int64{9223372036854775807},
		UintSlice:     []uint{18446744073709551615},
		Uint8Slice:    []uint8{255},
		Uint16Slice:   []uint16{65535},
		Uint32Slice:   []uint32{4294967295},
		Uint64Slice:   []uint64{18446744073709551615},
		DurationSlice: []time.Duration{time.Minute},
		URLSlice:      []url.URL{*url2},
		RegexpSlice:   []regexp.Regexp{*re2},
	}

	values := map[string]string{
		"String":        "new",
		"Bool":          "true",
		"Float32":       "2.7182",
		"Float64":       "2.7182818284",
		"Int":           "9223372036854775807",
		"Int8":          "127",
		"Int16":         "32767",
		"Int32":         "2147483647",
		"Int64":         "9223372036854775807",
		"Uint":          "18446744073709551615",
		"Uint8":         "255",
		"Uint16":        "65535",
		"Uint32":        "4294967295",
		"Uint64":        "18446744073709551615",
		"Duration":      "1m",
		"URL":           "service-2",
		"Regexp":        "[:alpha:]",
		"StringPtr":     "new",
		"BoolPtr":       "true",
		"Float32Ptr":    "2.7182",
		"Float64Ptr":    "2.7182818284",
		"IntPtr":        "9223372036854775807",
		"Int8Ptr":       "127",
		"Int16Ptr":      "32767",
		"Int32Ptr":      "2147483647",
		"Int64Ptr":      "9223372036854775807",
		"UintPtr":       "18446744073709551615",
		"Uint8Ptr":      "255",
		"Uint16Ptr":     "65535",
		"Uint32Ptr":     "4294967295",
		"Uint64Ptr":     "18446744073709551615",
		"DurationPtr":   "1m",
		"URLPtr":        "service-2",
		"RegexpPtr":     "[:alpha:]",
		"StringSlice":   "new",
		"BoolSlice":     "true",
		"Float32Slice":  "2.7182",
		"Float64Slice":  "2.7182818284",
		"IntSlice":      "9223372036854775807",
		"Int8Slice":     "127",
		"Int16Slice":    "32767",
		"Int32Slice":    "2147483647",
		"Int64Slice":    "9223372036854775807",
		"UintSlice":     "18446744073709551615",
		"Uint8Slice":    "255",
		"Uint16Slice":   "65535",
		"Uint32Slice":   "4294967295",
		"Uint64Slice":   "18446744073709551615",
		"DurationSlice": "1m",
		"URLSlice":      "service-2",
		"RegexpSlice":   "[:alpha:]",
	}

	tests := []struct {
		name            string
		s               Fields
		values          map[string]string
		expectedUpdated bool
		expectedError   string
		expectedResult  Fields
	}{
		{
			"NewValues",
			f1,
			values,
			true, "",
			f2,
		},
		{
			"NoNewValues",
			f2,
			values,
			false, "",
			f2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vStruct := reflect.ValueOf(&tc.s).Elem()
			for i := 0; i < vStruct.NumField(); i++ {
				v := vStruct.Field(i)
				f := vStruct.Type().Field(i)

				updated, err := Value(v, ",", tc.values[f.Name])

				if tc.expectedError == "" {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, tc.expectedError)
				}

				assert.Equal(t, tc.expectedUpdated, updated)
			}

			assert.Equal(t, tc.expectedResult, tc.s)
		})
	}
}
