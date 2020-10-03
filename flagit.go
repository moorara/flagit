// Package flagit TODO:
// TODO: support nested fields
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

var (
	flagNameRE = regexp.MustCompile(`^[A-Za-z]([0-9A-Za-z-.]*[0-9A-Za-z])?$`)
	flagArgRE  = regexp.MustCompile("^-{1,2}[A-Za-z]([0-9A-Za-z-.]*[0-9A-Za-z])?")
)

type fieldInfo struct {
	value   reflect.Value
	name    string
	flag    string
	listSep string
}

// flagValue implements the flag.Value interface.
type flagValue struct {
	fieldInfo
	continueOnError bool
}

// String is called for getting and printing the default value.
// Default value is already included in the usage string.
func (v flagValue) String() string {
	return ""
}

func (v flagValue) Set(val string) error {
	if _, err := setFieldValue(v.fieldInfo, val); err != nil {
		if v.continueOnError {
			return nil
		}
		return err
	}

	return nil
}

func validateStruct(s interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(s) // reflect.Value --> v.Type(), v.Kind(), v.NumField()
	t := reflect.TypeOf(s)  // reflect.Type --> t.Kind(), t.Name(), t.NumField()

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

func getFlagValue(flag string) string {
	flagRegex := regexp.MustCompile("-{1,2}" + flag)

	for i, arg := range os.Args {
		if flagRegex.MatchString(arg) {
			if s := strings.Index(arg, "="); s > 0 {
				return arg[s+1:]
			}

			if i+1 < len(os.Args) {
				if val := os.Args[i+1]; !flagArgRE.MatchString(val) {
					return val
				}
			}

			// For boolean flags
			return "true"
		}
	}

	return ""
}

func setFieldValue(f fieldInfo, val string) (bool, error) {
	switch f.value.Kind() {
	case reflect.String:
		return setString(f.value, val)
	case reflect.Bool:
		return setBool(f.value, val)
	case reflect.Float32:
		return setFloat32(f.value, val)
	case reflect.Float64:
		return setFloat64(f.value, val)
	case reflect.Int:
		return setInt(f.value, val)
	case reflect.Int8:
		return setInt8(f.value, val)
	case reflect.Int16:
		return setInt16(f.value, val)
	case reflect.Int32:
		return setInt32(f.value, val)
	case reflect.Int64:
		return setInt64(f.value, val)
	case reflect.Uint:
		return setUint(f.value, val)
	case reflect.Uint8:
		return setUint8(f.value, val)
	case reflect.Uint16:
		return setUint16(f.value, val)
	case reflect.Uint32:
		return setUint32(f.value, val)
	case reflect.Uint64:
		return setUint64(f.value, val)
	case reflect.Struct:
		return setStruct(f.value, val)

	case reflect.Slice:
		tSlice := reflect.TypeOf(f.value.Interface()).Elem()
		vals := strings.Split(val, f.listSep)

		switch tSlice.Kind() {
		case reflect.String:
			return setStringSlice(f.value, vals)
		case reflect.Bool:
			return setBoolSlice(f.value, vals)
		case reflect.Float32:
			return setFloat32Slice(f.value, vals)
		case reflect.Float64:
			return setFloat64Slice(f.value, vals)
		case reflect.Int:
			return setIntSlice(f.value, vals)
		case reflect.Int8:
			return setInt8Slice(f.value, vals)
		case reflect.Int16:
			return setInt16Slice(f.value, vals)
		case reflect.Int32:
			return setInt32Slice(f.value, vals)
		case reflect.Int64:
			return setInt64Slice(f.value, vals)
		case reflect.Uint:
			return setUintSlice(f.value, vals)
		case reflect.Uint8:
			return setUint8Slice(f.value, vals)
		case reflect.Uint16:
			return setUint16Slice(f.value, vals)
		case reflect.Uint32:
			return setUint32Slice(f.value, vals)
		case reflect.Uint64:
			return setUint64Slice(f.value, vals)
		case reflect.Struct:
			return setURLSlice(f.value, vals)
		}
	}

	return false, fmt.Errorf("unsupported kind: %s", f.value.Kind())
}

func iterateOnFields(vStruct reflect.Value, continueOnError bool, handle func(f fieldInfo) error) error {
	// Iterate over struct fields
	for i := 0; i < vStruct.NumField(); i++ {
		v := vStruct.Field(i)        // reflect.Value       --> vField.Kind(), vField.Type().Name(), vField.Type().Kind(), vField.Interface()
		t := v.Type()                // reflect.Type        --> t.Kind(), t.PkgPath(), t.Name(), t.NumField()
		f := vStruct.Type().Field(i) // reflect.StructField --> f.Name, f.Type.Name(), f.Type.Kind(), f.Tag.Get(tag)

		// Skip unexported and unsupported fields
		if !v.CanSet() || !isTypeSupported(t) {
			continue
		}

		// `flag:"..."`
		flagName := f.Tag.Get(flagTag)
		if flagName == "" {
			continue
		}

		// Sanitize the flag name
		if !flagNameRE.MatchString(flagName) {
			if continueOnError {
				continue
			}
			return fmt.Errorf("invalid flag name: %s", flagName)
		}

		// `sep:"..."`
		listSep := f.Tag.Get(sepTag)
		if listSep == "" {
			listSep = ","
		}

		err := handle(fieldInfo{
			value:   v,
			name:    f.Name,
			flag:    flagName,
			listSep: listSep,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// Populate accepts the pointer to a struct type.
// For those struct fields that have the flag tag, it will read values from command-line flags and parse them to the appropriate types.
// This method does not use the built-in flag package for parsing and reading the flags.
func Populate(s interface{}, continueOnError bool) error {
	v, err := validateStruct(s)
	if err != nil {
		return err
	}

	return iterateOnFields(v, continueOnError, func(f fieldInfo) error {
		if val := getFlagValue(f.flag); val != "" {
			if _, err := setFieldValue(f, val); err != nil {
				if continueOnError {
					return nil
				}
				return err
			}
		}

		return nil
	})
}

// RegisterFlags accepts a flag set and the pointer to a struct type.
// For those struct fields that have the flag tag, it will register a flag on the given flag set.
// The current values of the struct fields will be used as default values for the registered flags.
// Once the Parse method on the flag set is called, the values will be read, parsed to the appropriate types, and assigned to the corresponding struct fields.
func RegisterFlags(fs *flag.FlagSet, s interface{}, continueOnError bool) error {
	v, err := validateStruct(s)
	if err != nil {
		return err
	}

	return iterateOnFields(v, continueOnError, func(f fieldInfo) error {
		if fs.Lookup(f.flag) != nil {
			if continueOnError {
				return nil
			}
			return fmt.Errorf("flag already registered: %s", f.flag)
		}

		// Create usage string
		var usage string
		switch f.value.Kind() {
		case reflect.Bool:
			usage = fmt.Sprintf(
				"%-15s %s\n%-15s %v",
				"data type:", f.value.Type(),
				"default value:", f.value.Interface(),
			)
		case reflect.Slice:
			usage = fmt.Sprintf(
				"%-15s []%s\n%-15s %v\n%-15s %s",
				"data type:", reflect.TypeOf(f.value.Interface()).Elem(),
				"default value:", "[]",
				"separator:", f.listSep,
			)
		case reflect.Struct:
			t := f.value.Type()
			if t.PkgPath() == "net/url" && t.Name() == "URL" {
				usage = fmt.Sprintf(
					"%-15s %s\n%-15s %v",
					"data type:", f.value.Type(),
					"default value:", "",
				)
			}
		default:
			usage = fmt.Sprintf(
				"%-15s %s\n%-15s %v",
				"data type:", f.value.Type(),
				"default value:", f.value.Interface(),
			)
		}

		// Register the flag
		switch f.value.Kind() {
		case reflect.Bool:
			// f.value.CanAddr() expected to be true
			// f.value.Addr().Interface().(*bool) expected to be ok
			ptr := f.value.Addr().Interface().(*bool)
			fs.BoolVar(ptr, f.flag, f.value.Bool(), usage)
		default:
			fv := &flagValue{f, continueOnError}
			fs.Var(fv, f.flag, usage)
		}

		return nil
	})
}
