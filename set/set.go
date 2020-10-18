package set

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// String sets a string value.
func String(v reflect.Value, val string) (bool, error) {
	if v.String() == val {
		return false, nil
	}

	v.SetString(val)
	return true, nil
}

// Bool sets a bool value.
func Bool(v reflect.Value, val string) (bool, error) {
	b, err := strconv.ParseBool(val)
	if err != nil {
		return false, err
	}

	if v.Bool() == b {
		return false, nil
	}

	v.SetBool(b)
	return true, nil
}

// Float32 sets a float32 value.
func Float32(v reflect.Value, val string) (bool, error) {
	f, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return false, err
	}

	if v.Float() == f {
		return false, nil
	}

	v.SetFloat(f)
	return true, nil
}

// Float64 sets a float64 value.
func Float64(v reflect.Value, val string) (bool, error) {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return false, err
	}

	if v.Float() == f {
		return false, nil
	}

	v.SetFloat(f)
	return true, nil
}

// Int sets an int value.
func Int(v reflect.Value, val string) (bool, error) {
	// int size and range are platform-dependent
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return false, err
	}

	if v.Int() == i {
		return false, nil
	}

	v.SetInt(i)
	return true, nil
}

// Int8 sets an int8 value.
func Int8(v reflect.Value, val string) (bool, error) {
	i, err := strconv.ParseInt(val, 10, 8)
	if err != nil {
		return false, err
	}

	if v.Int() == i {
		return false, nil
	}

	v.SetInt(i)
	return true, nil
}

// Int16 sets an int16 value.
func Int16(v reflect.Value, val string) (bool, error) {
	i, err := strconv.ParseInt(val, 10, 16)
	if err != nil {
		return false, err
	}

	if v.Int() == i {
		return false, nil
	}

	v.SetInt(i)
	return true, nil
}

// Int32 sets an int32 value.
func Int32(v reflect.Value, val string) (bool, error) {
	i, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return false, err
	}

	if v.Int() == i {
		return false, nil
	}

	v.SetInt(i)
	return true, nil
}

// Int64 sets an int64 value.
func Int64(v reflect.Value, val string) (bool, error) {
	if t := v.Type(); t.PkgPath() == "time" && t.Name() == "Duration" {
		d, err := time.ParseDuration(val)
		if err != nil {
			return false, err
		}

		if v.Interface() == d {
			return false, nil
		}

		v.Set(reflect.ValueOf(d))
		return true, nil
	}

	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return false, err
	}

	if v.Int() == i {
		return false, nil
	}

	v.SetInt(i)
	return true, nil
}

// Uint sets an uint value.
func Uint(v reflect.Value, val string) (bool, error) {
	// uint size and range are platform-dependent
	u, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return false, err
	}

	if v.Uint() == u {
		return false, nil
	}

	v.SetUint(u)
	return true, nil
}

// Uint8 sets an uint8 value.
func Uint8(v reflect.Value, val string) (bool, error) {
	u, err := strconv.ParseUint(val, 10, 8)
	if err != nil {
		return false, err
	}

	if v.Uint() == u {
		return false, nil
	}

	v.SetUint(u)
	return true, nil
}

// Uint16 sets an uint16 value.
func Uint16(v reflect.Value, val string) (bool, error) {
	u, err := strconv.ParseUint(val, 10, 16)
	if err != nil {
		return false, err
	}

	if v.Uint() == u {
		return false, nil
	}

	v.SetUint(u)
	return true, nil
}

// Uint32 sets an uint32 value.
func Uint32(v reflect.Value, val string) (bool, error) {
	u, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return false, err
	}

	if v.Uint() == u {
		return false, nil
	}

	v.SetUint(u)
	return true, nil
}

// Uint64 sets an uint64 value.
func Uint64(v reflect.Value, val string) (bool, error) {
	u, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return false, err
	}

	if v.Uint() == u {
		return false, nil
	}

	v.SetUint(u)
	return true, nil
}

// Struct sets a struct value.
func Struct(v reflect.Value, val string) (bool, error) {
	t := v.Type()

	if t.PkgPath() == "net/url" && t.Name() == "URL" {
		u, err := url.Parse(val)
		if err != nil {
			return false, err
		}

		// u is a pointer
		if reflect.DeepEqual(v.Interface(), *u) {
			return false, nil
		}

		// u is a pointer
		v.Set(reflect.ValueOf(u).Elem())
		return true, nil
	} else if t.PkgPath() == "regexp" && t.Name() == "Regexp" {
		r, err := regexp.CompilePOSIX(val)
		if err != nil {
			return false, err
		}

		// r is a pointer
		if reflect.DeepEqual(v.Interface(), *r) {
			return false, nil
		}

		// r is a pointer
		v.Set(reflect.ValueOf(r).Elem())
		return true, nil
	}

	return false, fmt.Errorf("unsupported type: %s.%s", t.PkgPath(), t.Name())
}

// StringPtr sets a string pointer.
func StringPtr(v reflect.Value, val string) (bool, error) {
	if !v.IsZero() && reflect.DeepEqual(v.Elem().Interface(), val) {
		return false, nil
	}

	v.Set(reflect.ValueOf(&val))
	return true, nil
}

// BoolPtr sets a bool pointer.
func BoolPtr(v reflect.Value, val string) (bool, error) {
	b, err := strconv.ParseBool(val)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Bool() == b {
		return false, nil
	}

	v.Set(reflect.ValueOf(&b))
	return true, nil
}

// Float32Ptr sets a float32 pointer.
func Float32Ptr(v reflect.Value, val string) (bool, error) {
	f64, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Float() == f64 {
		return false, nil
	}

	f32 := float32(f64)
	v.Set(reflect.ValueOf(&f32))
	return true, nil
}

// Float64Ptr sets a float64 pointer.
func Float64Ptr(v reflect.Value, val string) (bool, error) {
	f64, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Float() == f64 {
		return false, nil
	}

	v.Set(reflect.ValueOf(&f64))
	return true, nil
}

// IntPtr sets an int pointer.
func IntPtr(v reflect.Value, val string) (bool, error) {
	// int size and range are platform-dependent
	i64, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Int() == i64 {
		return false, nil
	}

	i := int(i64)
	v.Set(reflect.ValueOf(&i))
	return true, nil
}

// Int8Ptr sets an int8 pointer.
func Int8Ptr(v reflect.Value, val string) (bool, error) {
	i64, err := strconv.ParseInt(val, 10, 8)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Int() == i64 {
		return false, nil
	}

	i8 := int8(i64)
	v.Set(reflect.ValueOf(&i8))
	return true, nil
}

// Int16Ptr sets an int16 pointer.
func Int16Ptr(v reflect.Value, val string) (bool, error) {
	i64, err := strconv.ParseInt(val, 10, 16)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Int() == i64 {
		return false, nil
	}

	i16 := int16(i64)
	v.Set(reflect.ValueOf(&i16))
	return true, nil
}

// Int32Ptr sets an int32 pointer.
func Int32Ptr(v reflect.Value, val string) (bool, error) {
	i64, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Int() == i64 {
		return false, nil
	}

	i32 := int32(i64)
	v.Set(reflect.ValueOf(&i32))
	return true, nil
}

// Int64Ptr sets an int64 pointer.
func Int64Ptr(v reflect.Value, val string) (bool, error) {
	t := reflect.TypeOf(v.Interface()).Elem()

	if t.PkgPath() == "time" && t.Name() == "Duration" {
		d, err := time.ParseDuration(val)
		if err != nil {
			return false, err
		}

		if !v.IsZero() && v.Elem().Interface() == d {
			return false, nil
		}

		v.Set(reflect.ValueOf(&d))
		return true, nil
	}

	i64, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Int() == i64 {
		return false, nil
	}

	v.Set(reflect.ValueOf(&i64))
	return true, nil
}

// UintPtr sets an uint pointer.
func UintPtr(v reflect.Value, val string) (bool, error) {
	// uint size and range are platform-dependent
	u64, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Uint() == u64 {
		return false, nil
	}

	u := uint(u64)
	v.Set(reflect.ValueOf(&u))
	return true, nil
}

// Uint8Ptr sets an uint8 pointer.
func Uint8Ptr(v reflect.Value, val string) (bool, error) {
	u64, err := strconv.ParseUint(val, 10, 8)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Uint() == u64 {
		return false, nil
	}

	u8 := uint8(u64)
	v.Set(reflect.ValueOf(&u8))
	return true, nil
}

// Uint16Ptr sets an uint16 pointer.
func Uint16Ptr(v reflect.Value, val string) (bool, error) {
	u64, err := strconv.ParseUint(val, 10, 16)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Uint() == u64 {
		return false, nil
	}

	u16 := uint16(u64)
	v.Set(reflect.ValueOf(&u16))
	return true, nil
}

// Uint32Ptr sets an uint32 pointer.
func Uint32Ptr(v reflect.Value, val string) (bool, error) {
	u64, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Uint() == u64 {
		return false, nil
	}

	u32 := uint32(u64)
	v.Set(reflect.ValueOf(&u32))
	return true, nil
}

// Uint64Ptr sets an uint64 pointer.
func Uint64Ptr(v reflect.Value, val string) (bool, error) {
	u64, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Uint() == u64 {
		return false, nil
	}

	v.Set(reflect.ValueOf(&u64))
	return true, nil
}

// StructPtr sets a struct pointer.
func StructPtr(v reflect.Value, val string) (bool, error) {
	t := reflect.TypeOf(v.Interface()).Elem()

	if t.PkgPath() == "net/url" && t.Name() == "URL" {
		u, err := url.Parse(val)
		if err != nil {
			return false, err
		}

		if !v.IsZero() && reflect.DeepEqual(v.Elem().Interface(), *u) {
			return false, nil
		}

		// u is a pointer
		v.Set(reflect.ValueOf(u))
		return true, nil
	} else if t.PkgPath() == "regexp" && t.Name() == "Regexp" {
		r, err := regexp.CompilePOSIX(val)
		if err != nil {
			return false, err
		}

		if !v.IsZero() && reflect.DeepEqual(v.Elem().Interface(), *r) {
			return false, nil
		}

		// r is a pointer
		v.Set(reflect.ValueOf(r))
		return true, nil
	}

	return false, fmt.Errorf("unsupported type: %s.%s", t.PkgPath(), t.Name())
}

// StringSlice sets a string slice.
func StringSlice(v reflect.Value, vals []string) (bool, error) {
	if reflect.DeepEqual(v.Interface(), vals) {
		return false, nil
	}

	v.Set(reflect.ValueOf(vals))
	return true, nil
}

// BoolSlice sets a bool slice.
func BoolSlice(v reflect.Value, vals []string) (bool, error) {
	bools := []bool{}
	for _, val := range vals {
		b, err := strconv.ParseBool(val)
		if err != nil {
			return false, err
		}

		bools = append(bools, b)
	}

	if reflect.DeepEqual(v.Interface(), bools) {
		return false, nil
	}

	v.Set(reflect.ValueOf(bools))
	return true, nil
}

// Float32Slice sets a float32 slice.
func Float32Slice(v reflect.Value, vals []string) (bool, error) {
	floats := []float32{}
	for _, val := range vals {
		f, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return false, err
		}

		floats = append(floats, float32(f))
	}

	if reflect.DeepEqual(v.Interface(), floats) {
		return false, nil
	}

	v.Set(reflect.ValueOf(floats))
	return true, nil
}

// Float64Slice sets a float64 slice.
func Float64Slice(v reflect.Value, vals []string) (bool, error) {
	floats := []float64{}
	for _, val := range vals {
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return false, err
		}

		floats = append(floats, f)
	}

	if reflect.DeepEqual(v.Interface(), floats) {
		return false, nil
	}

	v.Set(reflect.ValueOf(floats))
	return true, nil
}

// IntSlice sets an int slice.
func IntSlice(v reflect.Value, vals []string) (bool, error) {
	// int size and range are platform-dependent
	ints := []int{}
	for _, val := range vals {
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return false, err
		}

		ints = append(ints, int(i))
	}

	if reflect.DeepEqual(v.Interface(), ints) {
		return false, nil
	}

	v.Set(reflect.ValueOf(ints))
	return true, nil
}

// Int8Slice sets an int8 slice.
func Int8Slice(v reflect.Value, vals []string) (bool, error) {
	ints := []int8{}
	for _, val := range vals {
		i, err := strconv.ParseInt(val, 10, 8)
		if err != nil {
			return false, err
		}

		ints = append(ints, int8(i))
	}

	if reflect.DeepEqual(v.Interface(), ints) {
		return false, nil
	}

	v.Set(reflect.ValueOf(ints))
	return true, nil
}

// Int16Slice sets an int16 slice.
func Int16Slice(v reflect.Value, vals []string) (bool, error) {
	ints := []int16{}
	for _, val := range vals {
		i, err := strconv.ParseInt(val, 10, 16)
		if err != nil {
			return false, err
		}

		ints = append(ints, int16(i))
	}

	if reflect.DeepEqual(v.Interface(), ints) {
		return false, nil
	}

	v.Set(reflect.ValueOf(ints))
	return true, nil
}

// Int32Slice sets an int32 slice.
func Int32Slice(v reflect.Value, vals []string) (bool, error) {
	ints := []int32{}
	for _, val := range vals {
		i, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return false, err
		}

		ints = append(ints, int32(i))
	}

	if reflect.DeepEqual(v.Interface(), ints) {
		return false, nil
	}

	v.Set(reflect.ValueOf(ints))
	return true, nil
}

// Int64Slice sets an int64 slice.
func Int64Slice(v reflect.Value, vals []string) (bool, error) {
	t := reflect.TypeOf(v.Interface()).Elem()

	if t.PkgPath() == "time" && t.Name() == "Duration" {
		durations := []time.Duration{}
		for _, val := range vals {
			d, err := time.ParseDuration(val)
			if err != nil {
				return false, err
			}

			durations = append(durations, d)
		}

		if reflect.DeepEqual(v.Interface(), durations) {
			return false, nil
		}

		v.Set(reflect.ValueOf(durations))
		return true, nil
	}

	ints := []int64{}
	for _, val := range vals {
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return false, err
		}

		ints = append(ints, i)
	}

	if reflect.DeepEqual(v.Interface(), ints) {
		return false, nil
	}

	v.Set(reflect.ValueOf(ints))
	return true, nil
}

// UintSlice sets an uint slice.
func UintSlice(v reflect.Value, vals []string) (bool, error) {
	// uint size and range are platform-dependent
	uints := []uint{}
	for _, val := range vals {
		u, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return false, err
		}

		uints = append(uints, uint(u))
	}

	if reflect.DeepEqual(v.Interface(), uints) {
		return false, nil
	}

	v.Set(reflect.ValueOf(uints))
	return true, nil
}

// Uint8Slice sets an uint8 slice.
func Uint8Slice(v reflect.Value, vals []string) (bool, error) {
	uints := []uint8{}
	for _, val := range vals {
		u, err := strconv.ParseUint(val, 10, 8)
		if err != nil {
			return false, err
		}

		uints = append(uints, uint8(u))
	}

	if reflect.DeepEqual(v.Interface(), uints) {
		return false, nil
	}

	v.Set(reflect.ValueOf(uints))
	return true, nil
}

// Uint16Slice sets an uint16 slice.
func Uint16Slice(v reflect.Value, vals []string) (bool, error) {
	uints := []uint16{}
	for _, val := range vals {
		u, err := strconv.ParseUint(val, 10, 16)
		if err != nil {
			return false, err
		}

		uints = append(uints, uint16(u))
	}

	if reflect.DeepEqual(v.Interface(), uints) {
		return false, nil
	}

	v.Set(reflect.ValueOf(uints))
	return true, nil
}

// Uint32Slice sets an uint32 slice.
func Uint32Slice(v reflect.Value, vals []string) (bool, error) {
	uints := []uint32{}
	for _, val := range vals {
		u, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return false, err
		}

		uints = append(uints, uint32(u))
	}

	if reflect.DeepEqual(v.Interface(), uints) {
		return false, nil
	}

	v.Set(reflect.ValueOf(uints))
	return true, nil
}

// Uint64Slice sets an uint64 slice.
func Uint64Slice(v reflect.Value, vals []string) (bool, error) {
	uints := []uint64{}
	for _, val := range vals {
		u, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return false, err
		}

		uints = append(uints, u)
	}

	if reflect.DeepEqual(v.Interface(), uints) {
		return false, nil
	}

	v.Set(reflect.ValueOf(uints))
	return true, nil
}

// StructSlice sets a struct slice.
func StructSlice(v reflect.Value, vals []string) (bool, error) {
	t := reflect.TypeOf(v.Interface()).Elem()

	if t.PkgPath() == "net/url" && t.Name() == "URL" {
		urls := []url.URL{}
		for _, val := range vals {
			u, err := url.Parse(val)
			if err != nil {
				return false, err
			}

			urls = append(urls, *u)
		}

		// []url.URL
		if reflect.DeepEqual(v.Interface(), urls) {
			return false, nil
		}

		v.Set(reflect.ValueOf(urls))
		return true, nil
	} else if t.PkgPath() == "regexp" && t.Name() == "Regexp" {
		regexps := []regexp.Regexp{}
		for _, val := range vals {
			r, err := regexp.CompilePOSIX(val)
			if err != nil {
				return false, err
			}

			regexps = append(regexps, *r)
		}

		// []regexp.Regexp
		if reflect.DeepEqual(v.Interface(), regexps) {
			return false, nil
		}

		v.Set(reflect.ValueOf(regexps))
		return true, nil
	}

	return false, fmt.Errorf("unsupported type: %s.%s", t.PkgPath(), t.Name())
}

// Value sets a supported value.
func Value(v reflect.Value, sep, val string) (bool, error) {
	switch v.Kind() {
	case reflect.String:
		return String(v, val)
	case reflect.Bool:
		return Bool(v, val)
	case reflect.Float32:
		return Float32(v, val)
	case reflect.Float64:
		return Float64(v, val)
	case reflect.Int:
		return Int(v, val)
	case reflect.Int8:
		return Int8(v, val)
	case reflect.Int16:
		return Int16(v, val)
	case reflect.Int32:
		return Int32(v, val)
	case reflect.Int64:
		return Int64(v, val)
	case reflect.Uint:
		return Uint(v, val)
	case reflect.Uint8:
		return Uint8(v, val)
	case reflect.Uint16:
		return Uint16(v, val)
	case reflect.Uint32:
		return Uint32(v, val)
	case reflect.Uint64:
		return Uint64(v, val)
	case reflect.Struct:
		return Struct(v, val)

	case reflect.Ptr:
		tPtr := reflect.TypeOf(v.Interface()).Elem()

		switch tPtr.Kind() {
		case reflect.String:
			return StringPtr(v, val)
		case reflect.Bool:
			return BoolPtr(v, val)
		case reflect.Float32:
			return Float32Ptr(v, val)
		case reflect.Float64:
			return Float64Ptr(v, val)
		case reflect.Int:
			return IntPtr(v, val)
		case reflect.Int8:
			return Int8Ptr(v, val)
		case reflect.Int16:
			return Int16Ptr(v, val)
		case reflect.Int32:
			return Int32Ptr(v, val)
		case reflect.Int64:
			return Int64Ptr(v, val)
		case reflect.Uint:
			return UintPtr(v, val)
		case reflect.Uint8:
			return Uint8Ptr(v, val)
		case reflect.Uint16:
			return Uint16Ptr(v, val)
		case reflect.Uint32:
			return Uint32Ptr(v, val)
		case reflect.Uint64:
			return Uint64Ptr(v, val)
		case reflect.Struct:
			return StructPtr(v, val)
		}

	case reflect.Slice:
		tSlice := reflect.TypeOf(v.Interface()).Elem()
		vals := strings.Split(val, sep)

		switch tSlice.Kind() {
		case reflect.String:
			return StringSlice(v, vals)
		case reflect.Bool:
			return BoolSlice(v, vals)
		case reflect.Float32:
			return Float32Slice(v, vals)
		case reflect.Float64:
			return Float64Slice(v, vals)
		case reflect.Int:
			return IntSlice(v, vals)
		case reflect.Int8:
			return Int8Slice(v, vals)
		case reflect.Int16:
			return Int16Slice(v, vals)
		case reflect.Int32:
			return Int32Slice(v, vals)
		case reflect.Int64:
			return Int64Slice(v, vals)
		case reflect.Uint:
			return UintSlice(v, vals)
		case reflect.Uint8:
			return Uint8Slice(v, vals)
		case reflect.Uint16:
			return Uint16Slice(v, vals)
		case reflect.Uint32:
			return Uint32Slice(v, vals)
		case reflect.Uint64:
			return Uint64Slice(v, vals)
		case reflect.Struct:
			return StructSlice(v, vals)
		}
	}

	return false, fmt.Errorf("unsupported kind: %s", v.Kind())
}
