// Package flagit TODO:
// TODO: Decide how to handle errors
package flagit

import (
	"errors"
	"flag"
	"os"
	"reflect"
	"regexp"
	"strings"
)

const (
	flagTag = "flag"
	sepTag  = "sep"
)

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
