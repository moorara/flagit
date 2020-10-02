package flagit

import (
	"net/url"
	"reflect"
	"strconv"
	"time"
)

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
