// Package flagit TODO:
// TODO: Decide how to handle errors
package flagit

import (
	"errors"
	"flag"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	flagTag = "flag"
	sepTag  = "sep"
)

type fieldInfo struct {
	v       reflect.Value
	name    string
	listSep string
}

func validateStruct(s interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(s) // reflect.Value --> v.Type(), v.Kind(), v.NumField()
	t := reflect.TypeOf(s)  // reflect.Type --> t.Name(), t.Kind(), t.NumField()

	// A pointer to a struct should be passed
	if t.Kind() != reflect.Ptr {
		return reflect.Value{}, errors.New("non-pointer type: you should pass a pointer to a struct type")
	}

	// Navigate to the pointer value
	v = v.Elem()
	t = t.Elem()

	if t.Kind() != reflect.Struct {
		return reflect.Value{}, errors.New("non-struct type: you should pass a pointer to a struct type")
	}

	return v, nil
}

func isTypeSupported(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.String:
		return true
	case reflect.Bool:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Slice:
		return isTypeSupported(t.Elem())
	case reflect.Struct:
		if t.PkgPath() == "net/url" && t.Name() == "URL" {
			return true
		}
	}

	return false
}

func getFlagValue(flagName string) string {
	flagRegex := regexp.MustCompile("-{1,2}" + flagName)
	genericRegex := regexp.MustCompile("^-{1,2}[A-Za-z].*")

	for i, arg := range os.Args {
		if flagRegex.MatchString(arg) {
			if s := strings.Index(arg, "="); s > 0 {
				return arg[s+1:]
			}

			if i+1 < len(os.Args) {
				if val := os.Args[i+1]; !genericRegex.MatchString(val) {
					return val
				}
			}

			// For boolean flags
			return "true"
		}
	}

	return ""
}

func setString(v reflect.Value, val string) bool {
	if v.String() != val {
		v.SetString(val)
		return true
	}

	return false
}

func setBool(v reflect.Value, val string) bool {
	if b, err := strconv.ParseBool(val); err == nil {
		if v.Bool() != b {
			v.SetBool(b)
			return true
		}
	}

	return false
}

func setFloat32(v reflect.Value, val string) bool {
	if f, err := strconv.ParseFloat(val, 32); err == nil {
		if v.Float() != f {
			v.SetFloat(f)
			return true
		}
	}

	return false
}

func setFloat64(v reflect.Value, val string) bool {
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		if v.Float() != f {
			v.SetFloat(f)
			return true
		}
	}

	return false
}

func setInt(v reflect.Value, val string) bool {
	// int size and range are platform-dependent
	if i, err := strconv.ParseInt(val, 10, 64); err == nil {
		if v.Int() != i {
			v.SetInt(i)
			return true
		}
	}

	return false
}

func setInt8(v reflect.Value, val string) bool {
	if i, err := strconv.ParseInt(val, 10, 8); err == nil {
		if v.Int() != i {
			v.SetInt(i)
			return true
		}
	}

	return false
}

func setInt16(v reflect.Value, val string) bool {
	if i, err := strconv.ParseInt(val, 10, 16); err == nil {
		if v.Int() != i {
			v.SetInt(i)
			return true
		}
	}

	return false
}

func setInt32(v reflect.Value, val string) bool {
	if i, err := strconv.ParseInt(val, 10, 32); err == nil {
		if v.Int() != i {
			v.SetInt(i)
			return true
		}
	}

	return false
}

func setInt64(v reflect.Value, val string) bool {
	if t := v.Type(); t.PkgPath() == "time" && t.Name() == "Duration" {
		// time.Duration
		if d, err := time.ParseDuration(val); err == nil {
			if v.Interface() != d {
				v.Set(reflect.ValueOf(d))
				return true
			}
		}
	} else if i, err := strconv.ParseInt(val, 10, 64); err == nil {
		if v.Int() != i {
			v.SetInt(i)
			return true
		}
	}

	return false
}

func setUint(v reflect.Value, val string) bool {
	// uint size and range are platform-dependent
	if u, err := strconv.ParseUint(val, 10, 64); err == nil {
		if v.Uint() != u {
			v.SetUint(u)
			return true
		}
	}

	return false
}

func setUint8(v reflect.Value, val string) bool {
	if u, err := strconv.ParseUint(val, 10, 8); err == nil {
		if v.Uint() != u {
			v.SetUint(u)
			return true
		}
	}

	return false
}

func setUint16(v reflect.Value, val string) bool {
	if u, err := strconv.ParseUint(val, 10, 16); err == nil {
		if v.Uint() != u {
			v.SetUint(u)
			return true
		}
	}

	return false
}

func setUint32(v reflect.Value, val string) bool {
	if u, err := strconv.ParseUint(val, 10, 32); err == nil {
		if v.Uint() != u {
			v.SetUint(u)
			return true
		}
	}

	return false
}

func setUint64(v reflect.Value, val string) bool {
	if u, err := strconv.ParseUint(val, 10, 64); err == nil {
		if v.Uint() != u {
			v.SetUint(u)
			return true
		}
	}

	return false
}

func setStruct(v reflect.Value, val string) bool {
	if t := v.Type(); t.PkgPath() == "net/url" && t.Name() == "URL" {
		// url.URL
		if u, err := url.Parse(val); err == nil {
			// u is a pointer
			if !reflect.DeepEqual(v.Interface(), *u) {
				v.Set(reflect.ValueOf(u).Elem())
				return true
			}
		}
	}

	return false
}

func setStringSlice(v reflect.Value, vals []string) bool {
	if !reflect.DeepEqual(v.Interface(), vals) {
		v.Set(reflect.ValueOf(vals))
		return true
	}

	return false
}

func setBoolSlice(v reflect.Value, vals []string) bool {
	bools := []bool{}
	for _, val := range vals {
		if b, err := strconv.ParseBool(val); err == nil {
			bools = append(bools, b)
		}
	}

	if !reflect.DeepEqual(v.Interface(), bools) {
		v.Set(reflect.ValueOf(bools))
		return true
	}

	return false
}

func setFloat32Slice(v reflect.Value, vals []string) bool {
	floats := []float32{}
	for _, val := range vals {
		if f, err := strconv.ParseFloat(val, 32); err == nil {
			floats = append(floats, float32(f))
		}
	}

	if !reflect.DeepEqual(v.Interface(), floats) {
		v.Set(reflect.ValueOf(floats))
		return true
	}

	return false
}

func setFloat64Slice(v reflect.Value, vals []string) bool {
	floats := []float64{}
	for _, val := range vals {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			floats = append(floats, f)
		}
	}

	if !reflect.DeepEqual(v.Interface(), floats) {
		v.Set(reflect.ValueOf(floats))
		return true
	}

	return false
}

func setIntSlice(v reflect.Value, vals []string) bool {
	// int size and range are platform-dependent
	ints := []int{}
	for _, val := range vals {
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			ints = append(ints, int(i))
		}
	}

	if !reflect.DeepEqual(v.Interface(), ints) {
		v.Set(reflect.ValueOf(ints))
		return true
	}

	return false
}

func setInt8Slice(v reflect.Value, vals []string) bool {
	ints := []int8{}
	for _, val := range vals {
		if i, err := strconv.ParseInt(val, 10, 8); err == nil {
			ints = append(ints, int8(i))
		}
	}

	if !reflect.DeepEqual(v.Interface(), ints) {
		v.Set(reflect.ValueOf(ints))
		return true
	}

	return false
}

func setInt16Slice(v reflect.Value, vals []string) bool {
	ints := []int16{}
	for _, val := range vals {
		if i, err := strconv.ParseInt(val, 10, 16); err == nil {
			ints = append(ints, int16(i))
		}
	}

	if !reflect.DeepEqual(v.Interface(), ints) {
		v.Set(reflect.ValueOf(ints))
		return true
	}

	return false
}

func setInt32Slice(v reflect.Value, vals []string) bool {
	ints := []int32{}
	for _, val := range vals {
		if i, err := strconv.ParseInt(val, 10, 32); err == nil {
			ints = append(ints, int32(i))
		}
	}

	if !reflect.DeepEqual(v.Interface(), ints) {
		v.Set(reflect.ValueOf(ints))
		return true
	}

	return false
}

func setInt64Slice(v reflect.Value, vals []string) bool {
	if t := reflect.TypeOf(v.Interface()).Elem(); t.PkgPath() == "time" && t.Name() == "Duration" {
		durations := []time.Duration{}
		for _, val := range vals {
			if d, err := time.ParseDuration(val); err == nil {
				durations = append(durations, d)
			}
		}

		// []time.Duration
		if !reflect.DeepEqual(v.Interface(), durations) {
			v.Set(reflect.ValueOf(durations))
			return true
		}
	} else {
		ints := []int64{}
		for _, val := range vals {
			if i, err := strconv.ParseInt(val, 10, 64); err == nil {
				ints = append(ints, i)
			}
		}

		if !reflect.DeepEqual(v.Interface(), ints) {
			v.Set(reflect.ValueOf(ints))
			return true
		}
	}

	return false
}

func setUintSlice(v reflect.Value, vals []string) bool {
	// uint size and range are platform-dependent
	uints := []uint{}
	for _, val := range vals {
		if u, err := strconv.ParseUint(val, 10, 64); err == nil {
			uints = append(uints, uint(u))
		}
	}

	if !reflect.DeepEqual(v.Interface(), uints) {
		v.Set(reflect.ValueOf(uints))
		return true
	}

	return false
}

func setUint8Slice(v reflect.Value, vals []string) bool {
	uints := []uint8{}
	for _, val := range vals {
		if u, err := strconv.ParseUint(val, 10, 8); err == nil {
			uints = append(uints, uint8(u))
		}
	}

	if !reflect.DeepEqual(v.Interface(), uints) {
		v.Set(reflect.ValueOf(uints))
		return true
	}

	return false
}

func setUint16Slice(v reflect.Value, vals []string) bool {
	uints := []uint16{}
	for _, val := range vals {
		if u, err := strconv.ParseUint(val, 10, 16); err == nil {
			uints = append(uints, uint16(u))
		}
	}

	if !reflect.DeepEqual(v.Interface(), uints) {
		v.Set(reflect.ValueOf(uints))
		return true
	}

	return false
}

func setUint32Slice(v reflect.Value, vals []string) bool {
	uints := []uint32{}
	for _, val := range vals {
		if u, err := strconv.ParseUint(val, 10, 32); err == nil {
			uints = append(uints, uint32(u))
		}
	}

	if !reflect.DeepEqual(v.Interface(), uints) {
		v.Set(reflect.ValueOf(uints))
		return true
	}

	return false
}

func setUint64Slice(v reflect.Value, vals []string) bool {
	uints := []uint64{}
	for _, val := range vals {
		if u, err := strconv.ParseUint(val, 10, 64); err == nil {
			uints = append(uints, u)
		}
	}

	if !reflect.DeepEqual(v.Interface(), uints) {
		v.Set(reflect.ValueOf(uints))
		return true
	}

	return false
}

func setURLSlice(v reflect.Value, vals []string) bool {
	t := reflect.TypeOf(v.Interface()).Elem()

	if t.PkgPath() == "net/url" && t.Name() == "URL" {
		urls := []url.URL{}
		for _, val := range vals {
			if u, err := url.Parse(val); err == nil {
				urls = append(urls, *u)
			}
		}

		// []url.URL
		if !reflect.DeepEqual(v.Interface(), urls) {
			v.Set(reflect.ValueOf(urls))
			return true
		}
	}

	return false
}

func setFieldValue(f fieldInfo, val string) bool {
	switch f.v.Kind() {
	case reflect.String:
		return setString(f.v, val)
	case reflect.Bool:
		return setBool(f.v, val)
	case reflect.Float32:
		return setFloat32(f.v, val)
	case reflect.Float64:
		return setFloat64(f.v, val)
	case reflect.Int:
		return setInt(f.v, val)
	case reflect.Int8:
		return setInt8(f.v, val)
	case reflect.Int16:
		return setInt16(f.v, val)
	case reflect.Int32:
		return setInt32(f.v, val)
	case reflect.Int64:
		return setInt64(f.v, val)
	case reflect.Uint:
		return setUint(f.v, val)
	case reflect.Uint8:
		return setUint8(f.v, val)
	case reflect.Uint16:
		return setUint16(f.v, val)
	case reflect.Uint32:
		return setUint32(f.v, val)
	case reflect.Uint64:
		return setUint64(f.v, val)
	case reflect.Struct:
		return setStruct(f.v, val)

	case reflect.Slice:
		tSlice := reflect.TypeOf(f.v.Interface()).Elem()
		vals := strings.Split(val, f.listSep)

		switch tSlice.Kind() {
		case reflect.String:
			return setStringSlice(f.v, vals)
		case reflect.Bool:
			return setBoolSlice(f.v, vals)
		case reflect.Float32:
			return setFloat32Slice(f.v, vals)
		case reflect.Float64:
			return setFloat64Slice(f.v, vals)
		case reflect.Int:
			return setIntSlice(f.v, vals)
		case reflect.Int8:
			return setInt8Slice(f.v, vals)
		case reflect.Int16:
			return setInt16Slice(f.v, vals)
		case reflect.Int32:
			return setInt32Slice(f.v, vals)
		case reflect.Int64:
			return setInt64Slice(f.v, vals)
		case reflect.Uint:
			return setUintSlice(f.v, vals)
		case reflect.Uint8:
			return setUint8Slice(f.v, vals)
		case reflect.Uint16:
			return setUint16Slice(f.v, vals)
		case reflect.Uint32:
			return setUint32Slice(f.v, vals)
		case reflect.Uint64:
			return setUint64Slice(f.v, vals)
		case reflect.Struct:
			return setURLSlice(f.v, vals)
		}
	}

	return false
}

func iterateOnFields(vStruct reflect.Value, handle func(v reflect.Value, fieldName, flagName, listSep string)) {
	// Iterate over struct fields
	for i := 0; i < vStruct.NumField(); i++ {
		v := vStruct.Field(i)        // reflect.Value --> vField.Kind(), vField.Type().Name(), vField.Type().Kind(), vField.Interface()
		f := vStruct.Type().Field(i) // reflect.StructField --> tField.Name, tField.Type.Name(), tField.Type.Kind(), tField.Tag.Get(tag)

		// Skip unexported and unsupported fields
		if !v.CanSet() || !isTypeSupported(v.Type()) {
			continue
		}

		// `flag:"..."`
		flagName := f.Tag.Get(flagTag)

		// `sep:"..."`
		listSep := f.Tag.Get(sepTag)
		if listSep == "" {
			listSep = ","
		}

		handle(v, f.Name, flagName, listSep)
	}
}

// Populate TODO:
func Populate(s interface{}) error {
	_, err := validateStruct(s)
	if err != nil {
		return err
	}

	return nil
}

// RegisterFlags TODO:
func RegisterFlags(fs flag.FlagSet, s interface{}) {

}
