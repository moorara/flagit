package flagit

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

func setString(v reflect.Value, val string) (bool, error) {
	if v.String() == val {
		return false, nil
	}

	v.SetString(val)
	return true, nil
}

func setBool(v reflect.Value, val string) (bool, error) {
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

func setFloat32(v reflect.Value, val string) (bool, error) {
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

func setFloat64(v reflect.Value, val string) (bool, error) {
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

func setInt(v reflect.Value, val string) (bool, error) {
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

func setInt8(v reflect.Value, val string) (bool, error) {
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

func setInt16(v reflect.Value, val string) (bool, error) {
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

func setInt32(v reflect.Value, val string) (bool, error) {
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

func setInt64(v reflect.Value, val string) (bool, error) {
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

func setUint(v reflect.Value, val string) (bool, error) {
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

func setUint8(v reflect.Value, val string) (bool, error) {
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

func setUint16(v reflect.Value, val string) (bool, error) {
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

func setUint32(v reflect.Value, val string) (bool, error) {
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

func setUint64(v reflect.Value, val string) (bool, error) {
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

func setStruct(v reflect.Value, val string) (bool, error) {
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

		v.Set(reflect.ValueOf(u).Elem())
		return true, nil
	}

	return false, fmt.Errorf("unsupported type: %s.%s", t.PkgPath(), t.Name())
}

func setStringSlice(v reflect.Value, vals []string) (bool, error) {
	if reflect.DeepEqual(v.Interface(), vals) {
		return false, nil
	}

	v.Set(reflect.ValueOf(vals))
	return true, nil
}

func setBoolSlice(v reflect.Value, vals []string) (bool, error) {
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

func setFloat32Slice(v reflect.Value, vals []string) (bool, error) {
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

func setFloat64Slice(v reflect.Value, vals []string) (bool, error) {
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

func setIntSlice(v reflect.Value, vals []string) (bool, error) {
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

func setInt8Slice(v reflect.Value, vals []string) (bool, error) {
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

func setInt16Slice(v reflect.Value, vals []string) (bool, error) {
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

func setInt32Slice(v reflect.Value, vals []string) (bool, error) {
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

func setInt64Slice(v reflect.Value, vals []string) (bool, error) {
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

func setUintSlice(v reflect.Value, vals []string) (bool, error) {
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

func setUint8Slice(v reflect.Value, vals []string) (bool, error) {
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

func setUint16Slice(v reflect.Value, vals []string) (bool, error) {
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

func setUint32Slice(v reflect.Value, vals []string) (bool, error) {
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

func setUint64Slice(v reflect.Value, vals []string) (bool, error) {
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

func setURLSlice(v reflect.Value, vals []string) (bool, error) {
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
	}

	return false, fmt.Errorf("unsupported type: %s.%s", t.PkgPath(), t.Name())
}
