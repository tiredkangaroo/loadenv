package loadenv

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

func parseLines(lines []string) (map[string]string, error) {
	envVariables := make(map[string]string)
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		ls := strings.Split(line, "=")
		if len(ls) != 2 {
			return nil, fmt.Errorf("bad syntax on line number %d.", i)
		}
		envVariables[ls[0]] = ls[1]
	}
	return envVariables, nil
}

func linesFromFile(filepath string) ([]string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return []string{}, err
	}
	lines := strings.Split(string(data), "\n")
	return lines, nil
}

func loadFile(filepath string) error {
	lines, err := linesFromFile(filepath)
	if err != nil {
		return err
	}
	variables, err := parseLines(lines)
	if err != nil {
		return err
	}
	for k, v := range variables {
		err := os.Setenv(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// Load loads environment variables from the filepaths specified.
// If no filepaths are specified, it defaults to reading .env.
// It may return an error when a file cannot be read, the file has
// bad syntax, or the syscall to set the environment variable fails.
func Load(filepaths ...string) error {
	if len(filepaths) == 0 {
		filepaths = []string{".env"}
	}
	for _, filepath := range filepaths {
		err := loadFile(filepath)
		if err != nil {
			return err
		}
	}
	return nil
}

// Unmarshal unmarshals environment variables into a struct.
// If no filepaths are specified, it defaults to reading .env.
//
// Supported types are: string, int, int8, int16, int32, int64,
// uint, uint8, uint16, uint32, and uint64.
//
// A struct field can specify whether or not the environment variable
// is required to be provided in any of the .env files with the `required`
// struct tag. If it is not provided, it defaults to true.
//
// It may return an error if s is not a pointer to a struct, a file cannot be read,
// the file has bad syntax, a required field is not provided, or a provided field cannot
// be made into the type in the struct field.
func Unmarshal(s any, filepaths ...string) error {
	if len(filepaths) == 0 {
		filepaths = append(filepaths, ".env")
	}

	st := reflect.TypeOf(s)
	if st.Kind() != reflect.Ptr {
		return fmt.Errorf("s must be POINTER to a struct")
	}
	st = st.Elem()
	if st.Kind() != reflect.Struct {
		return fmt.Errorf("s must be a pointer to a STRUCT")
	}
	sv := reflect.ValueOf(s).Elem()

	variables := make(map[string]string)
	for _, filepath := range filepaths {
		lines, err := linesFromFile(filepath)
		if err != nil {
			return err
		}
		vs, err := parseLines(lines)
		if err != nil {
			return err
		}
		for k, v := range vs {
			variables[k] = v // merge maps
		}
	}
	for i := range st.NumField() {
		fieldt := st.Field(i)
		field := sv.Field(i)

		requireds := fieldt.Tag.Get("required")
		var required bool
		if requireds == "" {
			required = true
		} else {
			r, err := strconv.ParseBool(requireds)
			if err != nil {
				return fmt.Errorf("required struct tag must be valid bool value")
			}
			required = r
		}

		value, ok := variables[fieldt.Name]
		if !ok {
			if required {
				return fmt.Errorf("required environment variable %s not provided.", fieldt.Name)
			}
			continue
		}

		switch fieldt.Type.Kind() {
		case reflect.String:
			field.SetString(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("environment variable %s has value %s. it cannot be made into type %s.",
					fieldt.Name,
					value,
					fieldt.Type.Kind().String(),
				)
			}
			vpi64 := (*int64)(unsafe.Pointer(&v))
			field.SetInt(*vpi64)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, err := strconv.Atoi(value)
			if err != nil || v < 0 {
				return fmt.Errorf("environment variable %s has value %s. it cannot be made into type %s.",
					fieldt.Name,
					value,
					fieldt.Type.Kind().String(),
				)
			}
			vpui64 := (*uint64)(unsafe.Pointer(&v))
			field.SetUint(*vpui64)
		case reflect.Bool:
			v, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("environment variable %s has value %s. it cannot be made into type %s.",
					fieldt.Name,
					value,
					fieldt.Type.Kind().String(),
				)
			}
			field.SetBool(v)
		default:
			return fmt.Errorf("unsupported field kind: %s", field.Type().Kind().String())
		}
	}
	return nil
}
