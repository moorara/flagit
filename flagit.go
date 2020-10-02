// Package flagit TODO:
// TODO: support nested fields
// TODO: Decide how to handle errors
package flagit

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
)

const (
	flagTag = "flag"
	sepTag  = "sep"
)

type fieldInfo struct {
	v    reflect.Value
	name string
	flag string
	sep  string
}

// flagValue implements the flag.Value interface.
type flagValue struct{}

func (v *flagValue) String() string {
	return ""
}

func (v *flagValue) Set(string) error {
	return nil
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
		vals := strings.Split(val, f.sep)

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

func iterateOnFields(vStruct reflect.Value, handle func(f fieldInfo)) {
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
		if flagName == "" {
			continue
		}

		// `sep:"..."`
		listSep := f.Tag.Get(sepTag)
		if listSep == "" {
			listSep = ","
		}

		handle(fieldInfo{
			v:    v,
			name: f.Name,
			flag: flagName,
			sep:  listSep,
		})
	}
}

// Populate accepts the pointer to a struct type.
// For those struct fields that have the flag tag, it will read values from command-line flags and parse them to the appropriate types.
// This method does not use the built-in flag package for parsing and reading the flags.
func Populate(s interface{}) error {
	v, err := validateStruct(s)
	if err != nil {
		return err
	}

	iterateOnFields(v, func(f fieldInfo) {
		if val := getFlagValue(f.flag); val != "" {
			setFieldValue(f, val)
		}
	})

	return nil
}

// RegisterFlags accepts a flag set and the pointer to a struct type.
// For those struct fields that have the flag tag, it will register a flag on the given flag set.
// The current values of the struct fields will be used as default values for the registered flags.
// Once the Parse method on the flag set is called, the values will be read, parsed to the appropriate types, and assigned to the corresponding struct fields.
func RegisterFlags(fs flag.FlagSet, s interface{}) error {
	v, err := validateStruct(s)
	if err != nil {
		return err
	}

	iterateOnFields(v, func(f fieldInfo) {
		var dataType string
		if v.Kind() == reflect.Slice {
			dataType = fmt.Sprintf("[]%s", reflect.TypeOf(v.Interface()).Elem())
		} else {
			dataType = v.Type().String()
		}

		usage := fmt.Sprintf(
			"%-15s %s\n%-15s %v",
			"data type", dataType,
			"default value", v.Interface(),
		)

		fs.Var(new(flagValue), f.flag, usage)

		/* if flag.Lookup(flagName) == nil {
		} */

		/* switch v.Kind() {
		case reflect.String:
			fs.StringVar(nil, f.flag, nil, usage)
		case reflect.Bool:
			fs.BoolVar(nil, f.flag, nil, usage)
		case reflect.Float32, reflect.Float64:
			fs.Float64Var(nil, f.flag, nil, usage)
		case reflect.Int:
			fs.IntVar(nil, f.flag, nil, usage)
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fs.Int64Var(nil, f.flag, nil, usage)
		case reflect.Uint:
			fs.UintVar(nil, f.flag, nil, usage)
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fs.Uint64Var(nil, f.flag, nil, usage)
		case reflect.Slice:
			// TODO:
		case reflect.Struct:
			if t.PkgPath() == "net/url" && t.Name() == "URL" {
				// TODO:
			}
		} */
	})

	return nil
}
