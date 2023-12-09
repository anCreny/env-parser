package env_parser

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type EnvParser struct {
	tag     string
	useName bool
	safe    bool
	delim   string
}

func New(tag, delim string, useName, safe bool) *EnvParser {
	return &EnvParser{tag, useName, safe, delim}
}

func (e *EnvParser) Parse(structure interface{}) error {
	structType := reflect.TypeOf(structure)

	switch {
	case structType.Kind() != reflect.Ptr:
		return fmt.Errorf("structure to parse must be a pointer, but got %v type", structType.Kind())
	case structType.Elem().Kind() != reflect.Struct:
		return fmt.Errorf("object %v had to be a pointer to a structure, but it is not", structType.String())
	}

	t := structType.Elem()
	v := reflect.ValueOf(structure)

	for i := 0; i < t.NumField(); i++ {
		fType := t.Field(i)
		fValue := v.Elem().Field(i)

		allocIfNil(&fValue)

		if err := e.fillField(fType, fValue, ""); err != nil {
			return err
		}
	}

	return nil
}

func (e *EnvParser) fillField(fType reflect.StructField, fValue reflect.Value, tagValue string) error {
	if fValue.Kind() != reflect.Ptr {
		if tValue := fType.Tag.Get(e.tag); tValue != "" {
			tagValue += strings.ToUpper(tValue) + e.delim
		} else if e.useName {
			tagValue += strings.ToUpper(fType.Name) + e.delim
		}
	}

	kind := fValue.Kind()

	switch kind {
	case reflect.Struct:
		v := fValue.Addr().Interface()

		tt := reflect.TypeOf(v).Elem()
		vv := reflect.ValueOf(v)

		for i := 0; i < tt.NumField(); i++ {
			ffType := tt.Field(i)
			ffValue := vv.Elem().Field(i)

			allocIfNil(&ffValue)

			if err := e.fillField(ffType, ffValue, tagValue); err != nil {
				return err
			}
		}

	case reflect.Ptr:
		if err := e.fillField(fType, fValue.Elem(), tagValue); err != nil {
			return err
		}

	case reflect.Array, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Func, reflect.Map, reflect.Slice, reflect.Interface:
		return nil

	case reflect.Invalid:
		return fmt.Errorf("type is invalid")

	default:
		if e.safe && !fValue.IsZero() {
			return nil
		}

		tagValue = strings.TrimSuffix(tagValue, e.delim)

		envValue := os.Getenv(tagValue)

		if envValue == "" {
			return nil
		}

		switch kind {
		case reflect.Bool:
			value, err := strconv.ParseBool(envValue)
			if err != nil {
				return err
			}

			fValue.SetBool(value)

		case reflect.String:
			fValue.SetString(envValue)

		case reflect.Int:
			value, err := strconv.Atoi(envValue)
			if err != nil {
				return err
			}

			fValue.SetInt(int64(value))

		case reflect.Int8:
			value, err := strconv.ParseInt(envValue, 10, 8)
			if err != nil {
				return err
			}

			fValue.SetInt(value)

		case reflect.Int16:
			value, err := strconv.ParseInt(envValue, 10, 16)
			if err != nil {
				return err
			}

			fValue.SetInt(value)

		case reflect.Int32:
			value, err := strconv.ParseInt(envValue, 10, 32)
			if err != nil {
				return err
			}

			fValue.SetInt(value)

		case reflect.Int64:
			value, err := strconv.ParseInt(envValue, 10, 64)
			if err != nil {
				return err
			}

			fValue.SetInt(value)

		case reflect.Uint:
			value, err := strconv.ParseUint(envValue, 10, 64)
			if err != nil {
				return err
			}

			fValue.SetUint(value)

		case reflect.Uint8:
			value, err := strconv.ParseUint(envValue, 10, 8)
			if err != nil {
				return err
			}

			fValue.SetUint(value)

		case reflect.Uint16:
			value, err := strconv.ParseUint(envValue, 10, 16)
			if err != nil {
				return err
			}

			fValue.SetUint(value)

		case reflect.Uint32:
			value, err := strconv.ParseUint(envValue, 10, 32)
			if err != nil {
				return err
			}

			fValue.SetUint(value)

		case reflect.Uint64:
			value, err := strconv.ParseUint(envValue, 10, 64)
			if err != nil {
				return err
			}

			fValue.SetUint(value)

		case reflect.Float32:
			value, err := strconv.ParseFloat(envValue, 32)
			if err != nil {
				return err
			}

			fValue.SetFloat(value)

		case reflect.Float64:
			value, err := strconv.ParseFloat(envValue, 64)
			if err != nil {
				return err
			}

			fValue.SetFloat(value)
		}

	}

	return nil
}

func allocIfNil(value *reflect.Value) {
	if value.Kind() == reflect.Pointer && value.Elem().Kind() == reflect.Invalid {
		allocField := reflect.New(value.Type().Elem())
		value.Set(allocField)
	}
}
